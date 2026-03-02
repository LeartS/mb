package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/config"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with a Metabase instance",
	}
	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newLogoutCmd())
	return cmd
}

func newLoginCmd() *cobra.Command {
	var host, apiKey, username string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Metabase instance",
		Long: `Authenticate with a Metabase instance using an API key (recommended)
or username/password credentials. When using --username, the password
is prompted interactively (hidden input).

Examples:
  mb auth login --host https://metabase.example.com --api-key mb_XXXX
  mb auth login --host https://metabase.example.com --username admin@example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if host == "" {
				return fmt.Errorf("--host is required")
			}

			var sessionToken string

			if apiKey == "" && username == "" {
				return fmt.Errorf("provide --api-key or --username")
			}

			if apiKey == "" {
				// Authenticate via session with interactive password prompt.
				fmt.Fprintf(os.Stderr, "Password for %s: ", username)
				passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(os.Stderr) // newline after hidden input
				if err != nil {
					return fmt.Errorf("reading password: %w", err)
				}
				password := string(passwordBytes)
				if password == "" {
					return fmt.Errorf("password cannot be empty")
				}

				c := client.NewWithCredentials(host, "", "")
				resp, err := c.Post("/api/session", map[string]string{
					"username": username,
					"password": password,
				})
				if err != nil {
					return fmt.Errorf("login failed: %w", err)
				}
				var result struct {
					ID string `json:"id"`
				}
				if err := json.Unmarshal(resp, &result); err != nil {
					return fmt.Errorf("parsing session response: %w", err)
				}
				sessionToken = result.ID
			}

			// Verify the credentials work.
			c := client.NewWithCredentials(host, apiKey, sessionToken)
			resp, err := c.Get("/api/user/current")
			if err != nil {
				return fmt.Errorf("credential verification failed: %w", err)
			}

			var user struct {
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}
			if err := json.Unmarshal(resp, &user); err != nil {
				return fmt.Errorf("parsing user info: %w", err)
			}

			// Save to config.
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cfg.DefaultHost = host
			cfg.Hosts[host] = config.Host{
				APIKey:       apiKey,
				SessionToken: sessionToken,
			}
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Logged in to %s as %s %s (%s)\n", host, user.FirstName, user.LastName, user.Email)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "Metabase instance URL (e.g. https://metabase.example.com)")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication (recommended)")
	cmd.Flags().StringVar(&username, "username", "", "Username for session authentication (password prompted interactively)")

	return cmd
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/user/current")
			if err != nil {
				return fmt.Errorf("not authenticated: %w", err)
			}

			jsonFlag, _ := cmd.Flags().GetBool("json")
			if jsonFlag {
				return output.Render(resp)
			}

			var user struct {
				ID        int    `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}
			if err := json.Unmarshal(resp, &user); err != nil {
				return output.Render(resp)
			}

			cfg, _ := config.Load()
			host, _ := cfg.ActiveHost()
			fmt.Fprintf(cmd.OutOrStdout(), "Host:  %s\nUser:  %s %s (%s)\nID:    %d\n", host, user.FirstName, user.LastName, user.Email, user.ID)
			return nil
		},
	}
}

func newLogoutCmd() *cobra.Command {
	var host string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials for a host",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			target := host
			if target == "" {
				target, err = cfg.ActiveHost()
				if err != nil {
					return fmt.Errorf("no host specified and no default host configured")
				}
			}
			if _, ok := cfg.Hosts[target]; !ok {
				return fmt.Errorf("no credentials stored for %s", target)
			}
			delete(cfg.Hosts, target)
			if cfg.DefaultHost == target {
				cfg.DefaultHost = ""
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Logged out of %s\n", target)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "Host to log out of (defaults to current host)")
	return cmd
}
