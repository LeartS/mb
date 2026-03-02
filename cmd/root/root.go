package root

import (
	"fmt"

	"github.com/spf13/cobra"

	apiCmd "github.com/LeartS/mb/cmd/api"
	authCmd "github.com/LeartS/mb/cmd/auth"
	cardCmd "github.com/LeartS/mb/cmd/card"
	collectionCmd "github.com/LeartS/mb/cmd/collection"
	configCmd "github.com/LeartS/mb/cmd/config"
	dashboardCmd "github.com/LeartS/mb/cmd/dashboard"
	databaseCmd "github.com/LeartS/mb/cmd/database"
	datasetCmd "github.com/LeartS/mb/cmd/dataset"
	searchCmd "github.com/LeartS/mb/cmd/search"
	tableCmd "github.com/LeartS/mb/cmd/table"
)

// Set via ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var jsonOutput bool

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mb <command> <subcommand> [flags]",
		Short:         "Metabase CLI",
		Long:          "Work with Metabase from the command line.",
		Version:       fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	cmd.AddCommand(authCmd.NewCmd())
	cmd.AddCommand(configCmd.NewCmd())
	cmd.AddCommand(cardCmd.NewCmd())
	cmd.AddCommand(dashboardCmd.NewCmd())
	cmd.AddCommand(collectionCmd.NewCmd())
	cmd.AddCommand(databaseCmd.NewCmd())
	cmd.AddCommand(datasetCmd.NewCmd())
	cmd.AddCommand(tableCmd.NewCmd())
	cmd.AddCommand(searchCmd.NewCmd())
	cmd.AddCommand(apiCmd.NewCmd())

	return cmd
}
