package report

import (
	"fmt"
	"io"
	"os"

	"github.com/subtlepseudonym/envcmp/internal/env"
)

// RenderEnvDiff prints a formatted diff result produced by env.Diff.
// When color is true ANSI escape codes are used for terminal output.
func RenderEnvDiff(result env.DiffResult, color bool) {
	RenderEnvDiffTo(os.Stdout, result, color)
}

// RenderEnvDiffTo writes the formatted diff to w.
func RenderEnvDiffTo(w io.Writer, result env.DiffResult, color bool) {
	if !result.HasDiff {
		if color {
			fmt.Fprintln(w, "\033[32m✔ no differences found\033[0m")
		} else {
			fmt.Fprintln(w, "✔ no differences found")
		}
		return
	}

	for _, e := range result.Entries {
		line := formatEntry(e, color)
		fmt.Fprintln(w, line)
	}
}
