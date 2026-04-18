// Package report formats and renders diff results for CLI output.
package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envcmp/internal/diff"
)

// Format controls how the report is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures report rendering.
type Options struct {
	Format  Format
	Masked  bool
	NoColor bool
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorYellow = "\033[33m"
	colorReset = "\033[0m"
)

// Render writes a human-readable diff report to w.
func Render(w io.Writer, results []diff.Result, opts Options) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			printColored(w, colorRed, opts.NoColor, fmt.Sprintf("- [MISSING RIGHT] %s", r.Key))
		case diff.Extra:
			printColored(w, colorGreen, opts.NoColor, fmt.Sprintf("+ [MISSING LEFT]  %s", r.Key))
		case diff.Changed:
			printColored(w, colorYellow, opts.NoColor, fmt.Sprintf("~ [CHANGED]       %s: %q => %q", r.Key, r.LeftValue, r.RightValue))
		}
	}
}

func printColored(w io.Writer, color string, noColor bool, msg string) {
	if noColor {
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintf(w, "%s%s%s\n", color, msg, colorReset)
}
