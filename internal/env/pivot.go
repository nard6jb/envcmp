package env

import (
	"sort"

	"github.com/jasonlovesdoggo/envcmp/internal/envfile"
)

// PivotOptions controls the behaviour of Pivot.
type PivotOptions struct {
	// Labels are the column headers corresponding to each input map.
	// If empty, "env0", "env1", … are used.
	Labels []string
}

// PivotRow represents a single key compared across multiple environments.
type PivotRow struct {
	Key    string
	Values []string // one value per input map; empty string when key is absent
	Same   bool     // true when all present values are identical
}

// Pivot takes N env maps and returns a table where each row represents one
// unique key and its value in every environment.
func Pivot(maps []map[string]envfile.Entry, opts PivotOptions) ([]string, []PivotRow) {
	labels := make([]string, len(maps))
	for i := range maps {
		if i < len(opts.Labels) && opts.Labels[i] != "" {
			labels[i] = opts.Labels[i]
		} else {
			labels[i] = fmt.Sprintf("env%d", i)
		}
	}

	// collect all unique keys
	keySet := map[string]struct{}{}
	for _, m := range maps {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]PivotRow, 0, len(keys))
	for _, k := range keys {
		row := PivotRow{Key: k, Values: make([]string, len(maps))}
		first := ""
		presentCount := 0
		for i, m := range maps {
			if e, ok := m[k]; ok {
				row.Values[i] = e.Value
				if presentCount == 0 {
					first = e.Value
				}
				presentCount++
			}
		}
		row.Same = presentCount == len(maps) && allEqual(row.Values, first)
		rows = append(rows, row)
	}
	return labels, rows
}

func allEqual(vals []string, ref string) bool {
	for _, v := range vals {
		if v != ref {
			return false
		}
	}
	return true
}
