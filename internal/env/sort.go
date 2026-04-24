package env

import (
	"sort"

	"github.com/jasonuc/envcmp/internal/diff"
)

// SortOptions controls how entries are sorted.
type SortOptions struct {
	// ByKey sorts entries alphabetically by key (default behaviour).
	ByKey bool
	// Reverse reverses the sort order.
	Reverse bool
	// GroupByStatus groups diff entries: missing-left, missing-right, changed, same.
	GroupByStatus bool
}

// SortEntries returns a new slice of diff.Entry values sorted according to opts.
func SortEntries(entries []diff.Entry, opts SortOptions) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	copy(out, entries)

	if opts.GroupByStatus {
		sort.SliceStable(out, func(i, j int) bool {
			si := statusRank(out[i])
			sj := statusRank(out[j])
			if si != sj {
				if opts.Reverse {
					return si > sj
				}
				return si < sj
			}
			return keyLess(out[i].Key, out[j].Key, opts.Reverse)
		})
		return out
	}

	sort.SliceStable(out, func(i, j int) bool {
		return keyLess(out[i].Key, out[j].Key, opts.Reverse)
	})
	return out
}

func keyLess(a, b string, reverse bool) bool {
	if reverse {
		return a > b
	}
	return a < b
}

// statusRank assigns a numeric rank to a diff entry for group-by-status sorting.
func statusRank(e diff.Entry) int {
	switch {
	case e.LeftVal == "" && e.RightVal != "":
		return 0 // missing in left
	case e.LeftVal != "" && e.RightVal == "":
		return 1 // missing in right
	case e.LeftVal != e.RightVal:
		return 2 // changed
	default:
		return 3 // same
	}
}
