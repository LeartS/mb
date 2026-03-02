package dashboard

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard <command>",
		Short: "Manage dashboards",
	}
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newAddCardCmd())
	cmd.AddCommand(newCopyCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/dashboard/")
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
		Short: "Get a dashboard by ID (includes cards)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/dashboard/" + args[0])
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
		collectionID int
		fromJSON     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a dashboard",
		Long: `Create a new dashboard.

Examples:
  mb dashboard create --name "Sales Overview" --collection 5
  mb dashboard create --from-json dashboard.json`,
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
				payload = map[string]any{
					"name":        name,
					"description": description,
					"parameters":  []any{},
				}
				if collectionID > 0 {
					payload["collection_id"] = collectionID
				}
			}

			resp, err := c.Post("/api/dashboard/", payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Dashboard name")
	cmd.Flags().StringVar(&description, "description", "", "Dashboard description")
	cmd.Flags().IntVar(&collectionID, "collection", 0, "Collection ID")
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "Path to JSON file with full dashboard definition")

	return cmd
}

func newUpdateCmd() *cobra.Command {
	var fromJSON string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a dashboard",
		Args:  cobra.ExactArgs(1),
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
			resp, err := c.Put("/api/dashboard/"+args[0], payload)
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
		Short: "Delete a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Delete("/api/dashboard/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newAddCardCmd() *cobra.Command {
	var fromJSON string

	cmd := &cobra.Command{
		Use:   "add-card <dashboard-id>",
		Short: "Set the cards on a dashboard",
		Long: `Replace the full set of dashcards on a dashboard. Provide the dashcard
definitions as a JSON file.

The JSON file should contain a "cards" array where each entry is a dashcard:

  {
    "cards": [
      {
        "id": -1,
        "card_id": 10,
        "row": 0,
        "col": 0,
        "size_x": 12,
        "size_y": 6,
        "parameter_mappings": [],
        "visualization_settings": {}
      }
    ]
  }

Use negative IDs (e.g. -1, -2) for new dashcards. The dashboard grid is 24 columns wide.

Example:
  mb dashboard add-card 42 --from-json cards.json`,
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
			resp, err := c.Put("/api/dashboard/"+args[0]+"/cards", payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "Path to JSON file with dashcard definitions")
	return cmd
}

func newCopyCmd() *cobra.Command {
	var (
		name         string
		description  string
		collectionID int
	)

	cmd := &cobra.Command{
		Use:   "copy <id>",
		Short: "Copy a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			payload := map[string]any{}
			if name != "" {
				payload["name"] = name
			}
			if description != "" {
				payload["description"] = description
			}
			if collectionID > 0 {
				payload["collection_id"] = collectionID
			}
			resp, err := c.Post("/api/dashboard/"+args[0]+"/copy", payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name for the copy")
	cmd.Flags().StringVar(&description, "description", "", "Description for the copy")
	cmd.Flags().IntVar(&collectionID, "collection", 0, "Collection ID for the copy")
	return cmd
}
