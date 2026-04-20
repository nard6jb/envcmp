package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envcmp/internal/interpolate"
)

// RenderInterpolation writes a human-readable summary of interpolation results
// to w. If color is true ANSI escape codes are used.
func RenderInterpolation(w io.Writer, res interpolate.Result, color bool) {
	if len(res.Issues) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32m✔ No interpolation issues found.\033[0m")
		} else {
			fmt.Fprintln(w, "✔ No interpolation issues found.")
		}
		return
	}

	// Group issues by key for stable, readable output.
	byKey := make(map[string][]interpolate.Issue)
	for _, iss := range res.Issues {
		byKey[iss.Key] = append(byKey[iss.Key], iss)
	}

	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if color {
		fmt.Fprintf(w, "\033[31m✖ %d interpolation issue(s) found:\033[0m\n", len(res.Issues))
	} else {
		fmt.Fprintf(w, "✖ %d interpolation issue(s) found:\n", len(res.Issues))
	}

	for _, k := range keys {
		for _, iss := range byKey[k] {
			if color {
				fmt.Fprintf(w, "  \033[33m%-20s\033[0m → %s\n", k, iss.Message)
			} else {
				fmt.Fprintf(w, "  %-20s → %s\n", k, iss.Message)
			}
		}
	}
}
