package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/subtlepseudonym/envcmp/internal/env"
)

// RenderEnvMerge prints the result of an env.Merge operation to stdout.
func RenderEnvMerge(result env.MergeResult, color bool) {
	RenderEnvMergeTo(os.Stdout, result, color)
}

// RenderEnvMergeTo writes the merge result to the provided writer.
func RenderEnvMergeTo(w io.Writer, result env.MergeResult, color bool) {
	if len(result.Conflicts) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32m✔ merge complete — no conflicts\033[0m")
		} else {
			fmt.Fprintln(w, "✔ merge complete — no conflicts")
		}
		return
	}

	fmt.Fprintf(w, "merge complete — %d conflict(s) resolved\n", len(result.Conflicts))

	conflicts := make([]env.MergeConflict, len(result.Conflicts))
	copy(conflicts, result.Conflicts)
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	for _, c := range conflicts {
		if color {
			fmt.Fprintf(w, "  \033[33m~ %s\033[0m\n", c.Key)
			fmt.Fprintf(w, "      left:   \033[31m%s\033[0m\n", c.Left)
			fmt.Fprintf(w, "      right:  \033[31m%s\033[0m\n", c.Right)
			fmt.Fprintf(w, "      chosen: \033[32m%s\033[0m\n", c.Chosen)
		} else {
			fmt.Fprintf(w, "  ~ %s\n", c.Key)
			fmt.Fprintf(w, "      left:   %s\n", c.Left)
			fmt.Fprintf(w, "      right:  %s\n", c.Right)
			fmt.Fprintf(w, "      chosen: %s\n", c.Chosen)
		}
	}
}
