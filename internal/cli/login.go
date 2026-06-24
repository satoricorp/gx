package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/telemetry"
)

func newLoginCommand(ctx context.Context) *cobra.Command {
	var token string
	var apiURL string
	var orgID string
	cmd := &cobra.Command{
		Use:    "login",
		Short:  "Log in for capture upload to the GX server",
		Hidden: true,
		Long:   "Store upload credentials for POST /v1/extracts and /v1/sessions. For gx cloud GitHub login, use `gx auth login`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return fmt.Errorf("upload token required — pass --token (device flow not yet wired for capture server)")
			}
			creds := auth.Credentials{
				Token:  token,
				APIURL: apiURL,
				OrgID:  orgID,
			}
			if creds.APIURL == "" {
				creds.APIURL = auth.DefaultAPIURL
			}
			if err := auth.Save(creds); err != nil {
				return err
			}
			telemetry.EmitProductEvent(ctx, telemetry.EventCLIAuthLogin, map[string]any{
				"status":             "success",
				"auth_kind":          "upload_token",
				"api_url_configured": apiURL != "",
				"org_id_set":         orgID != "",
			})
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Upload API", creds.APIURL))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Capture upload", success("ok")))
			if creds.OrgID != "" {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Org", creds.OrgID))
			}
			fmt.Fprintln(cmd.OutOrStdout(), muted("Cloud GitHub login: `gx auth login`"))
			return nil
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "GX CLI upload bearer token")
	cmd.Flags().StringVar(&apiURL, "api-url", "", "GX server base URL (default http://localhost:3201)")
	cmd.Flags().StringVar(&orgID, "org-id", "", "optional org id for telemetry")
	return cmd
}
