package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/satoricorp/lgtm/internal/hooks"
)

func newInternalHooksCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "__hooks",
		Hidden: true,
	}
	cmd.AddCommand(newPrepareCommitMsgHookCommand(ctx))
	cmd.AddCommand(newPostCommitHookCommand(ctx))
	cmd.AddCommand(newPostRewriteHookCommand(ctx))
	return cmd
}

func newPrepareCommitMsgHookCommand(ctx context.Context) *cobra.Command {
	var repoRoot string
	var messagePath string
	cmd := &cobra.Command{
		Use:   "prepare-commit-msg",
		Short: "Git prepare-commit-msg hook adapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repoRoot == "" {
				repoRoot, _ = os.Getwd()
			}
			return hooks.PrepareCommitMsg(hooks.PrepareCommitMsgOptions{
				RepoRoot:    repoRoot,
				MessagePath: messagePath,
			})
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root")
	cmd.Flags().StringVar(&messagePath, "message-path", "", "commit message file path")
	return cmd
}

func newPostCommitHookCommand(ctx context.Context) *cobra.Command {
	var repoRoot string
	cmd := &cobra.Command{
		Use:   "post-commit",
		Short: "Git post-commit hook adapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repoRoot == "" {
				repoRoot, _ = os.Getwd()
			}
			return hooks.RunPostCommit(ctx, hooks.PostCommitOptions{RepoRoot: repoRoot})
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root")
	return cmd
}

func newPostRewriteHookCommand(ctx context.Context) *cobra.Command {
	var repoRoot string
	cmd := &cobra.Command{
		Use:   "post-rewrite",
		Short: "Git post-rewrite hook adapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repoRoot == "" {
				repoRoot, _ = os.Getwd()
			}
			return hooks.RunPostRewrite(ctx, hooks.PostRewriteOptions{RepoRoot: repoRoot})
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root")
	return cmd
}
