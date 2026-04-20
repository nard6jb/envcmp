package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envcmp/internal/schema"
)

// RenderSchema writes a human-readable schema validation report to w.
// If color is true, ANSI color codes are used.
func RenderSchema(w io.Writer, issues []schema.Issue, color bool) {
	if len(issues) == 0 {
		if color {
			fmt.Fprintln(w, green("✔ schema validation passed"))
		} else {
			fmt.Fprintln(w, "✔ schema validation passed")
		}
		return
	}

	// Sort by key for deterministic output.
	sorted := make([]schema.Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	if color {
		fmt.Fprintln(w, red(fmt.Sprintf("✖ schema validation failed (%d issue(s)):", len(issues))))
	} else {
		fmt.Fprintf(w, "✖ schema validation failed (%d issue(s)):\n", len(issues))
	}

	for _, issue := range sorted {
		line := fmt.Sprintf("  [%s] %s", issue.Key, issue.Message)
		if color {
			fmt.Fprintln(w, yellow(line))
		} else {
			fmt.Fprintln(w, line)
		}
	}
}
