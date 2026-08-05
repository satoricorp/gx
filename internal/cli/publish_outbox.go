package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/satoricorp/lgtm/internal/cloud"
	"github.com/satoricorp/lgtm/internal/publication"
	"github.com/spf13/cobra"
)

func newPublishUploadCommand(ctx context.Context) *cobra.Command {
	var quiet bool
	var limit int
	cmd := &cobra.Command{
		Use:    "__lgtm-upload-outbox",
		Short:  "Upload queued lgtm Cloud context",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return drainPublishUploadOutbox(ctx, cmd.OutOrStdout(), quiet, limit)
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "suppress upload summary")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum queued artifacts to upload")
	return cmd
}

func drainPublishUploadOutbox(ctx context.Context, out io.Writer, quiet bool, limit int) error {
	client := cloud.NewClient()
	if client == nil {
		return fmt.Errorf("lgtm cloud is not configured; set LGTM_CLOUD_URL or rebuild with cloud endpoints")
	}
	result, err := publication.DrainQueuedUploads(ctx, client, limit)
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Fprintln(out, labelValue("lgtm Cloud uploads", fmt.Sprintf("%d uploaded, %d failed, %d pending", result.Uploaded, result.Failed, result.Pending)))
	}
	return nil
}
