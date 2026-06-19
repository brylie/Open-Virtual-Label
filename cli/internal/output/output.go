package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

var quiet bool
var jsonMode bool

func SetQuiet(q bool) { quiet = q }
func SetJSON(j bool)  { jsonMode = j }

// Info prints to stderr unless quiet mode is on.
func Info(format string, args ...any) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// Print prints to stdout unless quiet mode is on.
func Print(format string, args ...any) {
	if !quiet {
		fmt.Printf(format+"\n", args...)
	}
}

// Success prints a checkmark-prefixed success message.
func Success(format string, args ...any) {
	if !quiet {
		fmt.Printf("✓ "+format+"\n", args...)
	}
}

// Fail prints a cross-prefixed failure message.
func Fail(format string, args ...any) {
	fmt.Printf("✗ "+format+"\n", args...)
}

// Error prints an error to stderr regardless of quiet mode.
func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// JSON marshals v and prints it to stdout.
func JSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// IsJSON returns true if JSON output mode is active.
func IsJSON() bool { return jsonMode }

// Table prints a formatted table.
func Table(headers []string, rows [][]string) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(headers)
	t.SetBorder(false)
	t.SetColumnSeparator("  ")
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeaderLine(false)
	for _, row := range rows {
		t.Append(row)
	}
	t.Render()
}
