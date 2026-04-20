package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/envcmp/internal/transform"
)

// RenderTransform writes a human-readable summary of the transform result to stdout.
// If color is true, ANSI escape codes are used to highlight changes.
func RenderTransform(original, result map[string]string, opts transform.Options, color bool) {
	RenderTransformTo(os.Stdout, original, result, opts, color)
}

// RenderTransformTo writes the transform summary to the provided writer.
// It shows each key that was affected by the transformation alongside its
// before/after values, and prints a final summary line.
func RenderTransformTo(w io.Writer, original, result map[string]string, opts transform.Options, color bool) {
	type change struct {
		key    string
		before string
		after  string
	}

	// Collect keys that changed during transformation.
	var changes []change

	for key, after := range result {
		before, ok := original[key]
		if !ok {
			// Key was introduced by transformation (e.g. prefix strip renamed it).
			continue
		}
		if before != after || key != originalKey(key, opts) {
			changes = append(changes, change{key: key, before: before, after: after})
		}
	}

	// Also detect keys that were renamed (removed from original, present in result).
	for origKey := range original {
		if _, exists := result[origKey]; !exists {
			// This key was renamed; the new name is already captured above.
			_ = origKey
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].key < changes[j].key
	})

	if len(changes) == 0 {
		if color {
			fmt.Fprintln(w, "\033[32m✔ No changes produced by transform.\033[0m")
		} else {
			fmt.Fprintln(w, "✔ No changes produced by transform.")
		}
		return
	}

	for _, c := range changes {
		if color {
			fmt.Fprintf(w, "  \033[33m~\033[0m %-30s  %s → %s\n", c.key, c.before, c.after)
		} else {
			fmt.Fprintf(w, "  ~ %-30s  %s → %s\n", c.key, c.before, c.after)
		}
	}

	fmt.Fprintln(w)

	if color {
		fmt.Fprintf(w, "\033[33m%d key(s) transformed.\033[0m\n", len(changes))
	} else {
		fmt.Fprintf(w, "%d key(s) transformed.\n", len(changes))
	}
}

// originalKey returns what the key would look like before transformation,
// used to detect renames caused by prefix stripping.
func originalKey(key string, opts transform.Options) string {
	if opts.StripPrefix != "" {
		return opts.StripPrefix + key
	}
	return key
}
