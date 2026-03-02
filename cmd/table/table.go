package table

import (
	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "table <command>",
		Short: "Manage tables",
	}
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newFieldsCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get table information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/table/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newFieldsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fields <id>",
		Short: "Get fields (columns) for a table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/table/" + args[0] + "/query_metadata")
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}
