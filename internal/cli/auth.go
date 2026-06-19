package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/cloud"
)

type authStatusJSON struct {
	LoggedIn      bool   `json:"loggedIn"`
	Login         string `json:"login,omitempty"`
	AvatarURL     string `json:"avatarURL,omitempty"`
	MachineName   string `json:"machineName,omitempty"`
	AuthKind      string `json:"authKind"`
	ConvexSiteURL string `json:"convexSiteURL,omitempty"`
	CloudURL      string `json:"cloudURL,omitempty"`
}

type authTokenJSON struct {
	Token    string `json:"token"`
	AuthKind string `json:"authKind"`
	CloudURL string `json:"cloudURL,omitempty"`
}

func buildAuthStatusJSON() (authStatusJSON, error) {
	convexSiteURL := cloud.ConvexSiteURL()
	cloudURL := cloud.CloudURL()
	creds, err := cloud.LoadCloudCredentials()
	if err != nil {
		return authStatusJSON{}, err
	}
	if creds != nil {
		return authStatusJSON{
			LoggedIn:      true,
			Login:         creds.Login,
			AvatarURL:     creds.AvatarURL,
			MachineName:   creds.MachineName,
			AuthKind:      "cloud",
			ConvexSiteURL: convexSiteURL,
			CloudURL:      cloudURL,
		}, nil
	}
	return authStatusJSON{LoggedIn: false, AuthKind: "none", ConvexSiteURL: convexSiteURL, CloudURL: cloudURL}, nil
}

func newAuthCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Log in to gx cloud",
	}
	cmd.AddCommand(
		newAuthLoginCommand(ctx),
		newAuthLogoutCommand(ctx),
		newAuthStatusCommand(),
		newAuthTokenCommand(),
	)
	return cmd
}

func newAuthLoginCommand(ctx context.Context) *cobra.Command {
	var machineName string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to gx cloud with GitHub",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := cloud.Login(ctx, cloud.LoginOptions{
				MachineName: machineName,
				Out:         cmd.OutOrStdout(),
			})
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Logged in", creds.Login))
			return nil
		},
	}
	cmd.Flags().StringVar(&machineName, "name", "", "Machine name shown in the gx console")
	return cmd
}

func newAuthLogoutCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out of gx cloud on this machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cloud.Logout(ctx, cloud.AuthEndpoints{}, nil); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), success("Logged out"))
			return nil
		},
	}
}

func newAuthStatusCommand() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show gx cloud login status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOut {
				status, err := buildAuthStatusJSON()
				if err != nil {
					return err
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(status)
			}

			creds, err := cloud.LoadCloudCredentials()
			if err != nil {
				return err
			}
			if creds != nil {
				fmt.Fprintln(cmd.OutOrStdout(), section("Cloud auth"))
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Login", creds.Login))
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Machine", creds.MachineName))
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Machine ID", creds.MachineID))
				printEffectiveCloudURLs(cmd)
				age := time.Since(creds.ObtainedAt.UTC())
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Session age", formatDuration(age)))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), muted("Not logged in."))
			printEffectiveCloudURLs(cmd)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

func newAuthTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "token",
		Short:  "Print the resolved gx cloud bearer token",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cloud.BearerToken()
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(authTokenJSON{
				Token:    token,
				AuthKind: "cloud",
				CloudURL: cloud.CloudURL(),
			})
		},
	}
	return cmd
}

func printEffectiveCloudURLs(cmd *cobra.Command) {
	if url := cloud.ConvexSiteURL(); url != "" {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Convex site", url))
	}
	if url := cloud.CloudURL(); url != "" {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Cloud URL", url))
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
