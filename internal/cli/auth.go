package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	AuthValid     *bool  `json:"authValid,omitempty"`
	AuthError     string `json:"authError,omitempty"`
	APIValid      *bool  `json:"apiValid,omitempty"`
	APIError      string `json:"apiError,omitempty"`
	ConvexSiteURL string `json:"convexSiteURL,omitempty"`
	CloudURL      string `json:"cloudURL,omitempty"`
}

type authTokenJSON struct {
	Token    string `json:"token"`
	AuthKind string `json:"authKind"`
	CloudURL string `json:"cloudURL,omitempty"`
}

func buildAuthStatusJSON(ctx context.Context) (authStatusJSON, error) {
	convexSiteURL := cloud.ConvexSiteURL()
	cloudURL := cloud.CloudURL()
	creds, err := cloud.LoadCloudCredentials()
	if err != nil {
		return authStatusJSON{}, err
	}
	if kind := authKindForCredentials(creds); kind != "none" {
		status := authStatusJSON{
			LoggedIn:      true,
			Login:         creds.Login,
			AvatarURL:     creds.AvatarURL,
			MachineName:   creds.MachineName,
			AuthKind:      kind,
			ConvexSiteURL: convexSiteURL,
			CloudURL:      cloudURL,
		}
		status = verifyAuthStatus(ctx, status, creds)
		return status, nil
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
		newAuthStatusCommand(ctx),
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

func newAuthStatusCommand(ctx context.Context) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show gx cloud login status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOut {
				status, err := buildAuthStatusJSON(ctx)
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
			kind := authKindForCredentials(creds)
			if kind != "none" {
				status := verifyAuthStatus(ctx, authStatusJSON{AuthKind: kind}, creds)
				fmt.Fprintln(cmd.OutOrStdout(), section("Cloud auth"))
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Auth", kind))
				if strings.TrimSpace(creds.Login) != "" {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Login", creds.Login))
				}
				if status.APIValid != nil && !*status.APIValid {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GX API", danger("warn")+": "+status.APIError))
				} else if status.APIValid != nil && *status.APIValid {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GX API", success("ok")))
				} else if status.APIError != "" {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GX API", danger("warn")+": "+status.APIError))
				}
				if status.AuthValid != nil && !*status.AuthValid {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GitHub token", danger("warn")+": "+status.AuthError))
				} else if status.AuthValid != nil && *status.AuthValid {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GitHub token", success("ok")))
				} else if status.AuthError != "" {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("GitHub token", danger("warn")+": "+status.AuthError))
				}
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

func verifyAuthStatus(ctx context.Context, status authStatusJSON, creds *cloud.CloudCredentials) authStatusJSON {
	if creds == nil {
		return status
	}
	verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if strings.TrimSpace(creds.CLISessionToken) != "" {
		token, _, err := cloud.CloudAPITokenWithKind()
		if err != nil {
			status.APIError = err.Error()
		} else {
			validation, err := cloud.ValidateCloudAPISession(verifyCtx, nil, token)
			if err != nil {
				status.APIError = "could not verify GX API session: " + err.Error()
			} else {
				status.APIValid = &validation.Valid
				if validation.Valid {
					if strings.TrimSpace(validation.Login) != "" {
						status.Login = validation.Login
					}
				} else {
					status.APIError = strings.TrimSpace(validation.Error)
					if status.APIError == "" {
						status.APIError = "GX API rejected stored session"
					}
					status.APIError += "; run `gx auth logout` then `gx auth login`"
				}
			}
		}
	}

	if strings.TrimSpace(creds.GitHubAccessToken) != "" {
		validation, err := cloud.ValidateGitHubAccessToken(verifyCtx, nil, creds.GitHubAccessToken)
		if err != nil {
			status.AuthError = "could not verify GitHub token: " + err.Error()
			return status
		}
		status.AuthValid = &validation.Valid
		if validation.Valid {
			if strings.TrimSpace(validation.Login) != "" {
				status.Login = validation.Login
			}
			return status
		}
		status.AuthError = strings.TrimSpace(validation.Error)
		if status.AuthError == "" {
			status.AuthError = "GitHub rejected stored token"
		}
		status.AuthError += "; run `gx auth logout` then `gx auth login`"
	}
	return status
}

func newAuthTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "token",
		Short:  "Print the resolved gx cloud bearer token",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, kind, err := cloud.CloudAPITokenWithKind()
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(authTokenJSON{
				Token:    token,
				AuthKind: kind,
				CloudURL: cloud.CloudURL(),
			})
		},
	}
	return cmd
}

func authKindForCredentials(creds *cloud.CloudCredentials) string {
	if creds == nil {
		return "none"
	}
	if strings.TrimSpace(creds.CLISessionToken) != "" {
		return "gx-cli"
	}
	if strings.TrimSpace(creds.GitHubAccessToken) != "" {
		return "github"
	}
	return "none"
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
