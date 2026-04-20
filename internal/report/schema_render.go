package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/your-org/envcmp/internal/schema"
)

const (
	colorRed    = "\x1b[31m"
	colorYellow = "\x1b[33m"
	colorReset  = "\x1b[0m"
	colorBold   = "\x1b[1m"
)

// RenderSchema writes a human-readable schema validation report to w.
// If color is true, ANSI escape codes are used for highlighting.
func RenderSchema(w io.Writer, issues []schema.Issue, color bool) {
	if len(issues) == 0 {
		if color {
			fmt.Fprintf(w, "%s✔ no schema issues found%s\n", "\x1b[32m", colorReset)
		} else {
			fmt.Fprintln(w, "✔ no schema issues found")
		}
		return
	}

	// Sort issues by key for deterministic output
	sorted := make([]schema.Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	word := "issues"
	if len(sorted) == 1 {
		word = "issue"
	}

	if color {
		fmt.Fprintf(w, "%s%d schema %s detected:%s\n", colorBold, len(sorted), word, colorReset)
	} else {
		fmt.Fprintf(w, "%d schema %s detected:\n", len(sorted), word)
	}

	for _, issue := range sorted {
		if color {
			fmt.Fprintf(w, "  %s%-24s%s %s%s%s\n",
				colorYellow, issue.Key, colorReset,
				colorRed, issue.Message, colorReset,
			)
		} else {
			fmt.Fprintf(w, "  %-24s %s\n", issue.Key, issue.Message)
		}
	}
}
