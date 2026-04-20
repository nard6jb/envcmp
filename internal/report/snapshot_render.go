package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envcmp/internal/diff"
)

// RenderSnapshot prints the result of comparing an env file against a snapshot.
// It reuses the same diff.Result type used by the standard diff report.
func RenderSnapshot(results []diff.Result, snapshotPath string, color bool, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	if len(results) == 0 {
		if color {
			fmt.Fprintf(w, "\033[32m✔ env matches snapshot: %s\033[0m\n", snapshotPath)
		} else {
			fmt.Fprintf(w, "✔ env matches snapshot: %s\n", snapshotPath)
		}
		return
	}

	if color {
		fmt.Fprintf(w, "\033[33m⚠ env differs from snapshot: %s\033[0m\n", snapshotPath)
	} else {
		fmt.Fprintf(w, "⚠ env differs from snapshot: %s\n", snapshotPath)
	}

	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		switch r.Status {
		case diff.MissingInRight:
			printColored(w, color, "\033[31m", fmt.Sprintf("  - removed:  %s\n", r.Key))
		case diff.MissingInLeft:
			printColored(w, color, "\033[32m", fmt.Sprintf("  + added:    %s = %s\n", r.Key, r.RightValue))
		case diff.Changed:
			printColored(w, color, "\033[33m", fmt.Sprintf("  ~ changed:  %s: %s → %s\n", r.Key, r.LeftValue, r.RightValue))
		}
	}
}
