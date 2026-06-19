package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	"github.com/satoricorp/gx/internal/storage"
)

func newIngestCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest external agent sessions into the gx store",
	}
	cmd.AddCommand(newIngestCursorCommand(ctx))
	return cmd
}

func newIngestCursorCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "cursor",
		Short: "Sync Cursor chat/composer sessions into the gx database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.Open(ctx)
			if err != nil {
				return err
			}
			defer db.Close()
			store, err := storage.NewStore(ctx, db)
			if err != nil {
				return err
			}
			defer store.Close()

			result, err := cursoringest.Sync(ctx, store)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Source", result.VSCDBPath))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Composers", fmt.Sprintf("%d", result.Composers)))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("New sessions", fmt.Sprintf("%d", result.NewSessions)))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("New messages", fmt.Sprintf("%d", result.NewMessages)))
			return nil
		},
	}
}
