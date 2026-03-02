package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/config"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage mb configuration",
	}
	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newListCmd())
	return cmd
}

var validKeys = map[string]string{
	"default_host": "The default Metabase instance URL",
}

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  "Valid keys: default_host",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			if _, ok := validKeys[key]; !ok {
				return fmt.Errorf("unknown config key %q; valid keys: default_host", key)
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			switch key {
			case "default_host":
				cfg.DefaultHost = value
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set %s = %s\n", key, value)
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			switch key {
			case "default_host":
				fmt.Fprintln(cmd.OutOrStdout(), cfg.DefaultHost)
			default:
				return fmt.Errorf("unknown config key %q", key)
			}
			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			jsonFlag, _ := cmd.Flags().GetBool("json")
			if jsonFlag {
				data, err := json.Marshal(cfg)
				if err != nil {
					return err
				}
				return output.Render(json.RawMessage(data))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "default_host: %s\n", cfg.DefaultHost)
			fmt.Fprintf(cmd.OutOrStdout(), "hosts:\n")
			for h, hcfg := range cfg.Hosts {
				authType := "api-key"
				if hcfg.APIKey == "" {
					authType = "session"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s)\n", h, authType)
			}
			return nil
		},
	}
}
