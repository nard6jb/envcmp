package report

import (
	"fmt"
	"io"
	"os"

	"github.com/subtlepseudonym/envcmp/internal/env"
)

// RenderUnique writes a human-readable summary of unique (exclusive) keys
// from a UniqueResult to stdout.
func RenderUnique(result env.UniqueResult, color bool) {
	RenderUniqueTo(os.Stdout, result, color)
}

// RenderUniqueTo writes the unique key report to the provided writer.
func RenderUniqueTo(w io.Writer, result env.UniqueResult, color bool) {
	if len(result.OnlyInLeft) == 0 && len(result.OnlyInRight) == 0 {
		if color {
			fmt.Fprintln(w, colorGreen+"✔ no unique keys — both files share the same key set"+colorReset)
		} else {
			fmt.Fprintln(w, "no unique keys — both files share the same key set")
		}
		return
	}

	for _, e := range result.OnlyInLeft {
		if color {
			fmt.Fprintf(w, "%s- [left only]  %s = %s%s\n", colorRed, e.Key, e.Left, colorReset)
		} else {
			fmt.Fprintf(w, "- [left only]  %s = %s\n", e.Key, e.Left)
		}
	}

	for _, e := range result.OnlyInRight {
		if color {
			fmt.Fprintf(w, "%s+ [right only] %s = %s%s\n", colorGreen, e.Key, e.Right, colorReset)
		} else {
			fmt.Fprintf(w, "+ [right only] %s = %s\n", e.Key, e.Right)
		}
	}

	total := len(result.OnlyInLeft) + len(result.OnlyInRight)
	if color {
		fmt.Fprintf(w, "%s%d unique key(s) found%s\n", colorYellow, total, colorReset)
	} else {
		fmt.Fprintf(w, "%d unique key(s) found\n", total)
	}
}
