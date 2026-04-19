package report

import (
	"fmt"
	"io"

	"github.com/yourusername/envcmp/internal/lint"
)

// RenderLint writes lint issues to w with optional colour.
// Returns true when there are no issues.
func RenderLint(w io.Writer, issues []lint.Issue, color bool) bool {
	if len(issues) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32m✔ no lint issues found\033[0m")
		} else {
			fmt.Fprintln(w, "✔ no lint issues found")
		}
		return true
	}

	for _, issue := range issues {
		if color {
			fmt.Fprintf(w, "\033[33m⚠ %s\033[0m\n", issue.String())
		} else {
			fmt.Fprintf(w, "⚠ %s\n", issue.String())
		}
	}
	return false
}
