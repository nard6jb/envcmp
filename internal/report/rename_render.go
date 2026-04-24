package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envcmp/internal/env"
)

// RenderRename prints a human-readable summary of a rename operation to stdout.
func RenderRename(res env.RenameResult, color bool) {
	RenderRenameTo(os.Stdout, res, color)
}

// RenderRenameTo writes the rename summary to w.
func RenderRenameTo(w io.Writer, res env.RenameResult, color bool) {
	if len(res.Renamed) == 0 && len(res.Dropped) == 0 && len(res.Skipped) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32mNo keys renamed.\033[0m")
		} else {
			fmt.Fprintln(w, "No keys renamed.")
		}
		return
	}

	// Renamed keys
	oldKeys := make([]string, 0, len(res.Renamed))
	for k := range res.Renamed {
		oldKeys = append(oldKeys, k)
	}
	sort.Strings(oldKeys)
	for _, old := range oldKeys {
		new := res.Renamed[old]
		if color {
			fmt.Fprintf(w, "\033[33m~ %s -> %s\033[0m\n", old, new)
		} else {
			fmt.Fprintf(w, "~ %s -> %s\n", old, new)
		}
	}

	// Dropped keys
	sort.Strings(res.Dropped)
	for _, k := range res.Dropped {
		if color {
			fmt.Fprintf(w, "\033[31m- %s (dropped)\033[0m\n", k)
		} else {
			fmt.Fprintf(w, "- %s (dropped)\n", k)
		}
	}

	// Skipped keys
	sort.Strings(res.Skipped)
	for _, k := range res.Skipped {
		if color {
			fmt.Fprintf(w, "\033[90m? %s (skipped, target exists)\033[0m\n", k)
		} else {
			fmt.Fprintf(w, "? %s (skipped, target exists)\n", k)
		}
	}
}
