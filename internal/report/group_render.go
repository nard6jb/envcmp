package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envcmp/internal/env"
)

// RenderGroup prints grouped env entries to stdout.
func RenderGroup(groups []env.GroupResult, color bool) {
	RenderGroupTo(os.Stdout, groups, color)
}

// RenderGroupTo writes grouped env entries to w.
func RenderGroupTo(w io.Writer, groups []env.GroupResult, color bool) {
	if len(groups) == 0 {
		fmt.Fprintln(w, "no entries")
		return
	}

	for _, g := range groups {
		header := g.Prefix
		if header == "" {
			header = "(ungrouped)"
		}
		if color {
			fmt.Fprintf(w, "\033[1;34m[%s]\033[0m\n", header)
		} else {
			fmt.Fprintf(w, "[%s]\n", header)
		}

		keys := make([]string, 0, len(g.Entries))
		for k := range g.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if color {
				fmt.Fprintf(w, "  \033[32m%s\033[0m=%s\n", k, g.Entries[k])
			} else {
				fmt.Fprintf(w, "  %s=%s\n", k, g.Entries[k])
			}
		}
	}
}
