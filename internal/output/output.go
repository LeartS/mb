package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"golang.org/x/term"
)

// IsTerminal returns true if stdout is a terminal.
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// JSON writes pretty-printed JSON to w.
func JSON(w io.Writer, data json.RawMessage) error {
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, data, "", "  "); err != nil {
		// Fall back to raw output if indenting fails.
		_, err := w.Write(data)
		return err
	}
	formatted.WriteByte('\n')
	_, err := w.Write(formatted.Bytes())
	return err
}

// Table writes rows as an aligned table to w.
// headers is the first row; rows contains the data.
func Table(w io.Writer, headers []string, rows [][]string) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	tw.Flush()
}

// Render outputs data as JSON (always, for now; table formatting for specific commands
// can call Table directly). This is the default output path.
func Render(data json.RawMessage) error {
	return JSON(os.Stdout, data)
}
