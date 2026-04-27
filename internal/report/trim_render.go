package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/jabuta/envcmp/internal/env"
)

// RenderTrim prints the results of a Trim operation to stdout.
func RenderTrim(results []env.TrimResult, color bool) {
	RenderTrimTo(os.Stdout, results, color)
}

// RenderTrimTo writes trim results to the given writer.
func RenderTrimTo(w io.Writer, results []env.TrimResult, color bool) {
	sorted := make([]env.TrimResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OriginalKey < sorted[j].OriginalKey
	})

	changed := 0
	for _, r := range sorted {
		if r.KeyChanged || r.ValueChanged {
			changed++
		}
	}

	if changed == 0 {
		if color {
			fmt.Fprintf(w, "\033[32m✔ no changes after trim (%d entries)\033[0m\n", len(results))
		} else {
			fmt.Fprintf(w, "✔ no changes after trim (%d entries)\n", len(results))
		}
		return
	}

	for _, r := range sorted {
		if !r.KeyChanged && !r.ValueChanged {
			continue
		}
		if r.KeyChanged && r.ValueChanged {
			if color {
				fmt.Fprintf(w, "\033[33m~ %s → %s = %q → %q\033[0m\n", r.OriginalKey, r.Key, r.OriginalValue, r.Value)
			} else {
				fmt.Fprintf(w, "~ %s → %s = %q → %q\n", r.OriginalKey, r.Key, r.OriginalValue, r.Value)
			}
		} else if r.KeyChanged {
			if color {
				fmt.Fprintf(w, "\033[33m~ key: %s → %s\033[0m\n", r.OriginalKey, r.Key)
			} else {
				fmt.Fprintf(w, "~ key: %s → %s\n", r.OriginalKey, r.Key)
			}
		} else {
			if color {
				fmt.Fprintf(w, "\033[33m~ %s value: %q → %q\033[0m\n", r.Key, r.OriginalValue, r.Value)
			} else {
				fmt.Fprintf(w, "~ %s value: %q → %q\n", r.Key, r.OriginalValue, r.Value)
			}
		}
	}

	fmt.Fprintf(w, "%d change(s) applied\n", changed)
}
