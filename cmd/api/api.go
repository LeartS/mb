package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/LeartS/mb/internal/client"
	"github.com/LeartS/mb/internal/output"
)

func NewCmd() *cobra.Command {
	var inputFile string

	cmd := &cobra.Command{
		Use:   "api <method> <path> [body]",
		Short: "Make a raw authenticated API request",
		Long: `Make an authenticated API request to any Metabase endpoint.
This is an escape hatch for endpoints not covered by dedicated commands.

The path should start with /api/. The body can be provided as an inline
argument or via --input-file.

Examples:
  mb api GET /api/user/current
  mb api POST /api/card/ '{"name":"test","dataset_query":...}'
  mb api PUT /api/card/42 --input-file payload.json
  mb api DELETE /api/card/42`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New()
			if err != nil {
				return err
			}

			method := strings.ToUpper(args[0])
			path := args[1]

			var body []byte
			if inputFile != "" {
				body, err = os.ReadFile(inputFile)
				if err != nil {
					return fmt.Errorf("reading input file: %w", err)
				}
			} else if len(args) > 2 {
				body = []byte(args[2])
			}

			data, status, err := c.DoRaw(method, path, body)
			if err != nil {
				return err
			}

			if status < 200 || status >= 300 {
				fmt.Fprintf(os.Stderr, "HTTP %d\n", status)
				os.Stdout.Write(data)
				os.Stdout.Write([]byte("\n"))
				os.Exit(1)
			}

			return output.Render(data)
		},
	}

	cmd.Flags().StringVar(&inputFile, "input-file", "", "Read request body from file")
	return cmd
}
