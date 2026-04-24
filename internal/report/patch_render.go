package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envcmp/internal/env"
)

// RenderPatch writes a human-readable summary of patch results to stdout.
func RenderPatch(results []env.PatchResult, color bool) {
	RenderPatchTo(os.Stdout, results, color)
}

// RenderPatchTo writes patch results to the provided writer.
func RenderPatchTo(w io.Writer, results []env.PatchResult, color bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No patch operations applied.")
		return
	}

	sorted := make([]env.PatchResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		switch {
		case r.Skipped:
			line := fmt.Sprintf("~ SKIP    %s", r.Key)
			if color {
				fmt.Fprintln(w, colorize("33", line))
			} else {
				fmt.Fprintln(w, line)
			}
		case r.Deleted:
			line := fmt.Sprintf("- DELETE  %s (was %q)", r.Key, r.OldValue)
			if color {
				fmt.Fprintln(w, colorize("31", line))
			} else {
				fmt.Fprintln(w, line)
			}
		case r.Added:
			line := fmt.Sprintf("+ ADD     %s=%q", r.Key, r.NewValue)
			if color {
				fmt.Fprintln(w, colorize("32", line))
			} else {
				fmt.Fprintln(w, line)
			}
		case r.Changed:
			line := fmt.Sprintf("~ CHANGE  %s: %q -> %q", r.Key, r.OldValue, r.NewValue)
			if color {
				fmt.Fprintln(w, colorize("34", line))
			} else {
				fmt.Fprintln(w, line)
			}
		default:
			fmt.Fprintf(w, "  NOOP    %s\n", r.Key)
		}
	}
}
