package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/evanw/envcmp/internal/merge"
)

// RenderMerge writes a human-readable merge result to stdout.
func RenderMerge(result merge.Result, color bool) {
	RenderMergeTo(os.Stdout, result, color)
}

// RenderMergeTo writes a human-readable merge result to the given writer.
func RenderMergeTo(w io.Writer, result merge.Result, color bool) {
	if len(result.Conflicts) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32m✔ Merge complete — no conflicts\033[0m")
		} else {
			fmt.Fprintln(w, "✔ Merge complete — no conflicts")
		}
		return
	}

	keys := make([]string, 0, len(result.Conflicts))
	for k := range result.Conflicts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if color {
		fmt.Fprintln(w, "\033[33m⚠ Merge completed with conflicts:\033[0m")
	} else {
		fmt.Fprintln(w, "⚠ Merge completed with conflicts:")
	}

	for _, k := range keys {
		vals := result.Conflicts[k]
		if color {
			fmt.Fprintf(w, "  \033[33m~ %s\033[0m (kept: %q, discarded: %q)\n", k, vals[0], vals[1])
		} else {
			fmt.Fprintf(w, "  ~ %s (kept: %q, discarded: %q)\n", k, vals[0], vals[1])
		}
	}
}
