package search

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	var modelFilter string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for cards, dashboards, collections, and more",
		Long: `Search across all Metabase objects.

Examples:
  mb search revenue
  mb search "monthly report" --model card
  mb search orders --model dashboard`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}
			q := strings.Join(args, " ")
			path := fmt.Sprintf("/api/search?q=%s", url.QueryEscape(q))
			if modelFilter != "" {
				path += "&models=" + url.QueryEscape(modelFilter)
			}
			resp, err := c.Get(path)
			if err != nil {
				return err
			}
			return output.Render(resp)
		},
	}

	cmd.Flags().StringVar(&modelFilter, "model", "", "Filter by model type (card, dashboard, collection, table, database)")
	return cmd
}
