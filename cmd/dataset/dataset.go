package dataset

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataset <command>",
		Short: "Run ad-hoc queries",
	}
	cmd.AddCommand(newQueryCmd())
	return cmd
}

func newQueryCmd() *cobra.Command {
	var (
		databaseID   int
		nativeQuery  string
		exportFormat string
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Execute an ad-hoc native SQL query",
		Long: `Execute an ad-hoc SQL query without saving it as a card.

Examples:
  mb dataset query --database 1 --native-query "SELECT count(*) FROM orders"
  mb dataset query --database 1 --native-query "SELECT * FROM orders LIMIT 10" --format csv`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			if databaseID == 0 {
				return fmt.Errorf("--database is required")
			}
			if nativeQuery == "" {
				return fmt.Errorf("--native-query is required")
			}

			payload := map[string]any{
				"database": databaseID,
				"type":     "native",
				"native": map[string]any{
					"query":         nativeQuery,
					"template-tags": map[string]any{},
				},
			}

			var path string
			if exportFormat != "" {
				path = "/api/dataset/" + exportFormat
			} else {
				path = "/api/dataset/"
			}

			resp, err := c.Post(path, payload)
			if err != nil {
				return err
			}

			if exportFormat != "" {
				os.Stdout.Write(resp)
				return nil
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().IntVar(&databaseID, "database", 0, "Database ID")
	cmd.Flags().StringVar(&nativeQuery, "native-query", "", "SQL query to execute")
	cmd.Flags().StringVar(&exportFormat, "format", "", "Export format: csv, json, xlsx (omit for default JSON)")

	return cmd
}
