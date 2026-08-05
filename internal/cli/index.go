package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/semantic"
	"github.com/satoricorp/gx/internal/vcs"
)

// newIndexCommand exposes repository indexing as a command so an index can be
// built (or rebuilt) without waiting for a push. It calls
// semantic.IndexRepository, the same incremental indexer every other indexing
// path uses.
func newIndexCommand(ctx context.Context) *cobra.Command {
	var full bool
	var jsonOut bool
	var quiet bool
	var concurrency int
	var namespace string

	cmd := &cobra.Command{
		Use:    "index",
		Short:  "Index this repository's source for review retrieval",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := indexRepoRoot(ctx)
			if err != nil {
				return err
			}
			orgID := ""
			if creds, ok := auth.LoadUpload(); ok {
				orgID = creds.OrgID
			}
			logf := func(format string, args ...any) {}
			if !quiet && !jsonOut {
				logf = func(format string, args ...any) {
					fmt.Fprintln(cmd.ErrOrStderr(), muted(fmt.Sprintf(format, args...)))
				}
			}
			result, err := semantic.IndexRepository(ctx, semantic.RepoIndexOptions{
				RepoRoot:    root,
				OrgID:       orgID,
				Namespace:   namespace,
				Reason:      "manual",
				Full:        full,
				Concurrency: concurrency,
				Logf:        logf,
			})
			if err != nil {
				return err
			}
			if jsonOut {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(result)
			}
			if quiet {
				return nil
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, labelValue("Namespace", result.Namespace))
			fmt.Fprintln(out, labelValue("Repository", firstNonEmptyString(result.RepoFullName, result.RepoRoot)))
			fmt.Fprintln(out, labelValue("Files", fmt.Sprintf("%d scanned, %d indexed, %d unchanged, %d removed",
				result.FilesScanned, result.FilesIndexed, result.FilesSkipped, result.FilesRemoved)))
			fmt.Fprintln(out, labelValue("Chunks", fmt.Sprintf("%d total, %d uploaded, %d deleted",
				result.ChunksTotal, result.ChunksUpserted, result.ChunksDeleted)))
			fmt.Fprintln(out, labelValue("Embedded", fmt.Sprintf("%d bytes in %d batches", result.BytesEmbedded, result.EmbedBatches)))
			fmt.Fprintln(out, labelValue("Duration", result.Duration.Round(1e6).String()))
			if result.UpToDate() {
				fmt.Fprintln(out, success("Index already up to date."))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&full, "full", false, "re-embed every file, ignoring the incremental manifest")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit the index result as JSON")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "suppress output")
	cmd.Flags().IntVar(&concurrency, "concurrency", 0, "parallel embed/upload batches (0 = default)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "override the TurboPuffer namespace")
	return cmd
}

// indexRepoRoot resolves the repository without requiring `gx init`, matching
// how review resolves it: indexing must work on a checkout gx has never
// touched, which is the entire point of on-demand indexing.
func indexRepoRoot(ctx context.Context) (string, error) {
	repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
	if err == nil && repo.RootPath != "" {
		return repo.RootPath, nil
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return "", err
	}
	return cwd, nil
}
