package database

import (
	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database <command>",
		Aliases: []string{"db"},
		Short:   "Manage databases",
	}
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newMetadataCmd())
	cmd.AddCommand(newSchemasCmd())
	cmd.AddCommand(newTablesCmd())
	cmd.AddCommand(newSyncCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/database/")
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get database by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/database/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newMetadataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metadata <id>",
		Short: "Get full metadata for a database (tables, fields, types)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/database/" + args[0] + "/metadata")
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newSchemasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schemas <id>",
		Short: "List schemas in a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/database/" + args[0] + "/schemas")
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newTablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tables <database-id> <schema>",
		Short: "List tables in a database schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/database/" + args[0] + "/schema/" + args[1])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync <id>",
		Short: "Trigger a schema sync for a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Post("/api/database/"+args[0]+"/sync_schema", map[string]any{})
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}
