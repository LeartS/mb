package card

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "card <command>",
		Aliases: []string{"question"},
		Short:   "Manage saved questions (cards)",
	}
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newQueryCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved questions",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/card/")
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
		Short: "Get a saved question by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/card/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newCreateCmd() *cobra.Command {
	var (
		name         string
		description  string
		display      string
		collectionID int
		databaseID   int
		nativeQuery  string
		fromJSON     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a saved question",
		Long: `Create a new saved question (card). Use --native-query for SQL questions
or --from-json for full control over the card definition.

Examples:
  mb card create --name "Revenue by Month" --database 1 --native-query "SELECT date_trunc('month', created_at) AS month, SUM(total) FROM orders GROUP BY 1"
  mb card create --from-json card.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}

			var payload map[string]any

			if fromJSON != "" {
				data, err := os.ReadFile(fromJSON)
				if err != nil {
					return fmt.Errorf("reading JSON file: %w", err)
				}
				if err := json.Unmarshal(data, &payload); err != nil {
					return fmt.Errorf("parsing JSON file: %w", err)
				}
			} else {
				if name == "" {
					return fmt.Errorf("--name is required")
				}
				if nativeQuery == "" {
					return fmt.Errorf("--native-query is required (or use --from-json)")
				}
				if databaseID == 0 {
					return fmt.Errorf("--database is required")
				}
				payload = map[string]any{
					"name":        name,
					"description": description,
					"display":     display,
					"dataset_query": map[string]any{
						"type":     "native",
						"database": databaseID,
						"native": map[string]any{
							"query":         nativeQuery,
							"template-tags": map[string]any{},
						},
					},
					"visualization_settings": map[string]any{},
				}
				if collectionID > 0 {
					payload["collection_id"] = collectionID
				}
			}

			resp, err := c.Post("/api/card/", payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Question name")
	cmd.Flags().StringVar(&description, "description", "", "Question description")
	cmd.Flags().StringVar(&display, "display", "table", "Visualization type (table, bar, line, pie, scalar, area, combo, pivot, funnel, map, row)")
	cmd.Flags().IntVar(&collectionID, "collection", 0, "Collection ID to save in")
	cmd.Flags().IntVar(&databaseID, "database", 0, "Database ID for the query")
	cmd.Flags().StringVar(&nativeQuery, "native-query", "", "Native SQL query")
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "Path to JSON file with full card definition")

	return cmd
}

func newUpdateCmd() *cobra.Command {
	var fromJSON string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a saved question",
		Long: `Update an existing card. Pass fields to update as a JSON file.

Example:
  mb card update 42 --from-json updates.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			if fromJSON == "" {
				return fmt.Errorf("--from-json is required")
			}
			data, err := os.ReadFile(fromJSON)
			if err != nil {
				return fmt.Errorf("reading JSON file: %w", err)
			}
			var payload map[string]any
			if err := json.Unmarshal(data, &payload); err != nil {
				return fmt.Errorf("parsing JSON file: %w", err)
			}
			resp, err := c.Put("/api/card/"+args[0], payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "Path to JSON file with update fields")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Archive a saved question",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Delete("/api/card/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newQueryCmd() *cobra.Command {
	var exportFormat string

	cmd := &cobra.Command{
		Use:   "query <id>",
		Short: "Execute a saved question and return results",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			var path string
			if exportFormat != "" {
				path = "/api/card/" + args[0] + "/query/" + exportFormat
			} else {
				path = "/api/card/" + args[0] + "/query"
			}
			resp, err := c.Post(path, map[string]any{})
			if err != nil {
				return err
			}
			if exportFormat != "" {
				// Raw output for CSV/XLSX.
				os.Stdout.Write(resp)
				return nil
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&exportFormat, "format", "", "Export format: csv, json, xlsx (omit for default JSON response)")
	return cmd
}
