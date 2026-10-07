package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/vcs"
)

// gx enhance — the one fix worth making, written for a model.
//
// It runs the same review `gx review` runs and reports its top fix as a
// paste-ready prompt. Two commands over one engine rather than a flag on
// review, because the audience differs and the audience decides the shape:
// review is read by a person deciding whether to ship, and prints lanes, a
// story, and a verdict; enhance is read by a coding model that has been asked
// to change one thing, and prints exactly one thing with no verdict, no
// second-place findings, and no ANSI.
//
// It never fails the build. `gx review` owns the gate and its exit codes;
// asking for a suggestion is not a gate, so enhance exits 0 whether or not
// the change would ship — including when there is nothing to fix, which it
// says plainly rather than exiting non-zero on a clean tree.
func newEnhanceCommand(ctx context.Context) *cobra.Command {
	var focus string
	var base string
	var wholeRepo bool
	var deep bool
	var fast bool
	cmd := &cobra.Command{
		Use:     "enhance [intent]",
		Aliases: []string{"gxe"},
		Short:   "The single highest-value fix in the current change, written to hand to your coding model",
		Long: `Print the one change most worth making right now, as a prompt you can paste
straight into a coding model — the problem, the code it lives in, the change to
make, and how to know it worked.

It runs the same review as ` + "`gx review`" + ` and reports the top of that review's
fix plan: blocking findings first, then the judge's own ranking, skipping
anything with no concrete action attached. So enhance and review never disagree
about what matters most — enhance is review's answer, narrowed to one item and
rewritten for a model rather than a person.

The optional INTENT is one sentence saying what the change was supposed to do
("make checkout survive gateway blips"), and travels into the prompt so the
model knows what the code was trying to achieve.

Output is plain text on stdout, so it pipes:

    gx enhance | pbcopy
    gx enhance | claude -p "apply this"

Exit status is always 0 — the gate is ` + "`gx review`" + `'s job, not this one. When
the review finds nothing to fix, enhance says so.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := codereview.ValidateModelPresetEnv(); err != nil {
				return err
			}
			prompt := ""
			if len(args) > 0 {
				prompt = strings.TrimSpace(args[0])
			}
			// Store-free, like review: suggesting a fix must not require (or
			// leave behind) gx state in the repository.
			repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
			if err != nil {
				return err
			}
			runReview := func(progress io.Writer) (codereview.Report, error) {
				return codereview.Review(ctx, repo.RootPath, codereview.Options{
					Deep:      deep,
					Focus:     focus,
					Base:      base,
					WholeRepo: wholeRepo,
					Prompt:    prompt,
					Fast:      fast,
					// The prompt is plain text for a model to read; color would
					// ride into its context as escape codes.
					Color:          false,
					ProgressWriter: progress,
				})
			}
			// The spinner goes to stderr so stdout carries the prompt alone and
			// `gx enhance | pbcopy` copies something usable.
			report, err := runReviewWithLoader(cmd.InOrStdin(), cmd.ErrOrStderr(), runReview)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if !report.Reviewed {
				fmt.Fprintln(cmd.ErrOrStderr(), labelValue("Enhance", muted("nothing to review — no change detected")))
				return nil
			}
			fix, ok := codereview.TopFix(report)
			if !ok {
				fmt.Fprintln(cmd.ErrOrStderr(), labelValue("Enhance", muted("the review found nothing to fix")))
				return nil
			}
			fmt.Fprint(out, codereview.RenderFixPrompt(report, fix))
			// The degradation notice belongs on stderr with the other
			// out-of-band lines: a partial review can still produce a real top
			// fix, and the reader should know the ranking behind it was thin.
			if reason := gateDegradedReason(report); reason != "" {
				fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Enhance", reason))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&focus, "focus", "", "limit the review to files under this path prefix")
	cmd.Flags().StringVar(&base, "base", "", "review the commit range <ref>...HEAD instead of the working tree, e.g. --base origin/main")
	cmd.Flags().BoolVar(&wholeRepo, "repo", false, "review the whole repository rather than just the current change")
	cmd.Flags().BoolVar(&deep, "deep", false, "run full-spectrum review with more local and indexed context")
	cmd.Flags().BoolVar(&fast, "fast", false, "optimize for wall clock: one reviewer instead of two, no verification pass")
	return cmd
}
