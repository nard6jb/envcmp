package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jasonlovesdoggo/envcmp/internal/env"
)

const (
	colorPivotSame    = "\033[32m" // green
	colorPivotDiff    = "\033[33m" // yellow
	colorPivotMissing = "\033[31m" // red
	colorReset        = "\033[0m"
)

// RenderPivot writes a pivot table to stdout.
func RenderPivot(labels []string, rows []env.PivotRow, color bool) {
	RenderPivotTo(os.Stdout, labels, rows, color)
}

// RenderPivotTo writes a pivot table to w.
func RenderPivotTo(w io.Writer, labels []string, rows []env.PivotRow, color bool) {
	if len(rows) == 0 {
		fmt.Fprintln(w, "no keys found")
		return
	}

	// compute column widths
	keyWidth := len("KEY")
	for _, r := range rows {
		if len(r.Key) > keyWidth {
			keyWidth = len(r.Key)
		}
	}
	colWidths := make([]int, len(labels))
	for i, l := range labels {
		colWidths[i] = len(l)
	}
	for _, r := range rows {
		for i, v := range r.Values {
			if len(v) > colWidths[i] {
				colWidths[i] = len(v)
			}
		}
	}

	// header
	header := fmt.Sprintf("%-*s", keyWidth, "KEY")
	for i, l := range labels {
		header += fmt.Sprintf("  %-*s", colWidths[i], l)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	// rows
	for _, r := range rows {
		line := fmt.Sprintf("%-*s", keyWidth, r.Key)
		for i, v := range r.Values {
			cell := v
			if cell == "" {
				if color {
					cell = colorPivotMissing + "(missing)" + colorReset
				} else {
					cell = "(missing)"
				}
				line += fmt.Sprintf("  %-*s", colWidths[i], cell)
			} else {
				if color {
					if r.Same {
						line += "  " + colorPivotSame + fmt.Sprintf("%-*s", colWidths[i], cell) + colorReset
					} else {
						line += "  " + colorPivotDiff + fmt.Sprintf("%-*s", colWidths[i], cell) + colorReset
					}
				} else {
					line += fmt.Sprintf("  %-*s", colWidths[i], cell)
				}
			}
		}
		fmt.Fprintln(w, line)
	}
}
