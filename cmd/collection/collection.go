package collection

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection <command>",
		Short: "Manage collections",
	}
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newTreeCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newItemsCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/collection/")
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newTreeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tree",
		Short: "Show collection hierarchy as a tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/collection/tree")
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
		Short: "Get a collection by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			resp, err := c.Get("/api/collection/" + args[0])
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}
}

func newCreateCmd() *cobra.Command {
	var (
		name     string
		parentID int
		color    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Long: `Create a new collection.

Examples:
  mb collection create --name "Analytics"
  mb collection create --name "Team Reports" --parent 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			payload := map[string]any{
				"name": name,
			}
			if parentID > 0 {
				payload["parent_id"] = parentID
			}
			if color != "" {
				payload["color"] = color
			}
			resp, err := c.Post("/api/collection/", payload)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Collection name")
	cmd.Flags().IntVar(&parentID, "parent", 0, "Parent collection ID")
	cmd.Flags().StringVar(&color, "color", "", "Collection color (hex, e.g. #509EE3)")

	return cmd
}

func newItemsCmd() *cobra.Command {
	var modelFilter string

	cmd := &cobra.Command{
		Use:   "items <id>",
		Short: "List items in a collection",
		Long: `List items in a collection. Use "root" as the ID for the root collection.

Examples:
  mb collection items 5
  mb collection items root
  mb collection items 5 --model card`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			path := "/api/collection/" + args[0] + "/items"
			if modelFilter != "" {
				path += "?models=" + modelFilter
			}
			resp, err := c.Get(path)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&modelFilter, "model", "", "Filter by model type (card, dashboard, collection, pulse)")
	return cmd
}
