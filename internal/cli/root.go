package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/inference"
	"github.com/satoricorp/gx/internal/postlist"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/satoricorp/gx/internal/version"
)

func NewRoot(ctx context.Context) *cobra.Command {
	engine := authoring.NewEngine()

	root := &cobra.Command{
		Use:           "gx",
		Short:         "gx CLI for Git-native capture, commits, and review",
		Long:          gxTagline,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.Contains(cmd.CommandPath(), "__") {
				return nil
			}
			inference.ApplyToEnvironment()
			telemetry.EmitInstallOnce(commandTelemetryContext(ctx, cmd))
			return ensureAutoInitializedRepo(ctx, engine, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			printInitNoteIfNeeded(cmd)
			return printRootHelp(cmd)
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	root.AddGroup(
		&cobra.Group{ID: groupSetup, Title: "Setup:"},
		&cobra.Group{ID: groupWork, Title: "Work:"},
		&cobra.Group{ID: groupHelp, Title: "Help:"},
	)

	root.AddCommand(
		newInternalHooksCommand(ctx),
		newVersionCommand(),
		newDoctorCommand(ctx),
		newAuthCommand(ctx),
		newInitCommand(ctx, engine),
		newCaptureCommand(ctx),
		newPublishUploadCommand(ctx),
		newReportCommand(ctx),
		newReviewCommand(ctx),
		newIndexCommand(ctx),
	)
	assignCommandGroups(root)
	rootRef := root
	root.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c != rootRef {
			if err := renderHelpTemplate(c, c.OutOrStdout()); err != nil {
				c.PrintErrln(err)
			}
			return
		}
		if err := printRootHelp(c); err != nil {
			c.PrintErrln(err)
		}
	})
	installHelpStyling(root)

	return root
}

func assignCommandGroups(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		switch cmd.Name() {
		case "init", "auth":
			cmd.GroupID = groupSetup
		case "review":
			cmd.GroupID = groupWork
		case "doctor", "version":
			cmd.GroupID = groupHelp
		}
	}
}

func printInitNoteIfNeeded(cmd *cobra.Command) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	cfg, err := gxconfig.LoadAt(cwd)
	if err != nil || cfg.HasIdentity() {
		return
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, danger("Run `gx init` first."))
	fmt.Fprintln(out)
}

func newVersionCommand() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "gx version",
		Run: func(cmd *cobra.Command, args []string) {
			if jsonOut {
				if err := json.NewEncoder(cmd.OutOrStdout()).Encode(version.BuildInfo()); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err)
				}
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), logoText("version "+version.Current()))
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable version metadata")
	return cmd
}

func newInitCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var name string
	var email string
	var yes bool
	var global bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Set up gx in the current repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			cleanupLegacyAmbientCaptureFromInit(ctx, cmd, yes)
			if global {
				if !yes {
					fmt.Fprintln(cmd.OutOrStdout(), commandLine("gx init --global", true))
					fmt.Fprintln(cmd.OutOrStdout())
				}
				err := runGlobalInit(ctx, cmd, yes)
				telemetry.EmitProductEvent(ctx, telemetry.EventCLIInitRun, map[string]any{
					"status":      initRunStatus(err),
					"interactive": !yes,
					"global":      true,
				})
				return err
			}
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout(), commandLine("gx init", true))
				fmt.Fprintln(cmd.OutOrStdout())
			}
			result, err := engine.Init(ctx, authoring.InitOptions{
				Name:        name,
				Email:       email,
				Interactive: !yes,
				In:          initInput(cmd, yes),
				Out:         initOutput(cmd, yes),
			})
			telemetry.EmitProductEvent(ctx, telemetry.EventCLIInitRun, map[string]any{
				"status":      initRunStatus(err),
				"interactive": !yes,
			})
			if err != nil {
				return err
			}
			if !yes && result.IdentityName != "" && result.IdentityEmail != "" {
				if result.IdentityApplied {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Configured", fmt.Sprintf("gx identity as %s <%s>", result.IdentityName, result.IdentityEmail)))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Using", fmt.Sprintf("gx identity %s <%s>", result.IdentityName, result.IdentityEmail)))
				}
				if client := postlist.NewFromEnv(); client != nil {
					err := client.UpsertIdentity(ctx, postlist.Identity{
						Name:      result.IdentityName,
						Email:     result.IdentityEmail,
						GxVersion: version.Current(),
					})
					if err != nil {
						fmt.Fprintln(cmd.ErrOrStderr(), danger(fmt.Sprintf("gx signup upload failed: %v", err)))
					}
				}
			}
			installHook := installCaptureHook
			if yes {
				installHook = installCaptureHookQuiet
			}
			if err := installHook(cmd, result.Repo.RootPath); err != nil {
				if !yes {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Pre-push hook", err.Error()))
				}
			}
			if err := installClaudeCaptureHooks(cmd, result.Repo.RootPath, yes); err != nil {
				if !yes {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Claude hooks", err.Error()))
				}
			}
			if err := initSetupFromCommand(cmd, result.Repo.RootPath, yes); err != nil {
				return err
			}
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Next", "git add <files> && git commit -m \"...\""))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "user name to store in gx config")
	cmd.Flags().StringVar(&email, "email", "", "user email to store in gx config")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "accept defaults and suppress successful init output")
	cmd.Flags().BoolVar(&global, "global", false, "install machine-wide git hooks (~/.gx/hooks) so gx works in every repo; skips per-repo setup")
	return cmd
}

func initRunStatus(err error) string {
	if err != nil {
		return "error"
	}
	return "success"
}

func initInput(cmd *cobra.Command, quiet bool) io.Reader {
	if quiet {
		return nil
	}
	return cmd.InOrStdin()
}

func initOutput(cmd *cobra.Command, quiet bool) io.Writer {
	if quiet {
		return nil
	}
	return cmd.OutOrStdout()
}

func currentHeadCommit(ctx context.Context, repoRoot string) string {
	if strings.TrimSpace(repoRoot) == "" {
		return ""
	}
	out, err := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func Execute(ctx context.Context) error {
	args := os.Args[1:]
	switch filepath.Base(os.Args[0]) {
	case "gxr":
		args = append([]string{"review"}, args...)
	}
	root := NewRoot(ctx)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		if !errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, err)
		}
		return err
	}
	return nil
}

// ExitCode maps an error to a process exit status. Commands that need to fail
// with a specific status wrap their error with vcs.CodedErrorf; everything else
// falls through to the caller's generic failure code.
func ExitCode(err error) int {
	var coded *vcs.CodedError
	if errors.As(err, &coded) && coded.Code != 0 {
		return coded.Code
	}
	return 0
}
