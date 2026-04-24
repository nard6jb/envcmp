package report

import (
	"fmt"
	"io"
	"os"

	"github.com/jasonuc/envcmp/internal/diff"
	"github.com/jasonuc/envcmp/internal/env"
)

// RenderSorted prints diff entries sorted according to opts.
// Color output is enabled when color is true.
func RenderSorted(entries []diff.Entry, opts env.SortOptions, color bool) {
	RenderSortedTo(os.Stdout, entries, opts, color)
}

// RenderSortedTo writes sorted diff entries to w.
func RenderSortedTo(w io.Writer, entries []diff.Entry, opts env.SortOptions, color bool) {
	sorted := env.SortEntries(entries, opts)

	if len(sorted) == 0 {
		if color {
			fmt.Fprintln(w, colorGreen+"No differences found."+colorReset)
		} else {
			fmt.Fprintln(w, "No differences found.")
		}
		return
	}

	for _, e := range sorted {
		switch {
		case e.LeftVal == "" && e.RightVal != "":
			printColored(w, color, colorGreen, fmt.Sprintf("+ %-30s %s", e.Key, e.RightVal))
		case e.LeftVal != "" && e.RightVal == "":
			printColored(w, color, colorRed, fmt.Sprintf("- %-30s %s", e.Key, e.LeftVal))
		case e.LeftVal != e.RightVal:
			printColored(w, color, colorYellow, fmt.Sprintf("~ %-30s %s → %s", e.Key, e.LeftVal, e.RightVal))
		default:
			fmt.Fprintf(w, "  %-30s %s\n", e.Key, e.LeftVal)
		}
	}
}
