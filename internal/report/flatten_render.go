package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envcmp/internal/env"
)

// RenderFlatten prints flattened key-value entries to stdout.
func RenderFlatten(entries []env.FlattenResult, color bool) {
	RenderFlattenTo(os.Stdout, entries, color)
}

// RenderFlattenTo writes flattened key-value entries to the given writer.
func RenderFlattenTo(w io.Writer, entries []env.FlattenResult, color bool) {
	if len(entries) == 0 {
		if color {
			fmt.Fprintln(w, colorGreen+"No entries to display."+colorReset)
		} else {
			fmt.Fprintln(w, "No entries to display.")
		}
		return
	}

	// Sort for deterministic output.
	sorted := make([]env.FlattenResult, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, e := range sorted {
		if color {
			fmt.Fprintf(w, "%s%s%s=%s\n", colorCyan, e.Key, colorReset, e.Value)
		} else {
			fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
		}
	}
}
