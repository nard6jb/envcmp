package env

import (
	"sort"

	"github.com/subtlepseudonym/envcmp/internal/diff"
)

// UniqueOptions controls the behaviour of the Unique function.
type UniqueOptions struct {
	// SortKeys sorts the resulting entries alphabetically by key.
	SortKeys bool
	// MaskSecrets replaces sensitive values with a masked placeholder.
	MaskSecrets bool
}

// UniqueResult holds the keys that appear in only one of the two maps.
type UniqueResult struct {
	// OnlyInLeft contains entries present in left but not right.
	OnlyInLeft []diff.Entry
	// OnlyInRight contains entries present in right but not left.
	OnlyInRight []diff.Entry
}

// Unique returns entries that are exclusive to each of the two provided maps.
// An entry is considered unique if its key does not appear in the other map at
// all, regardless of value.
func Unique(left, right map[string]string, opts UniqueOptions) UniqueResult {
	var result UniqueResult

	for k, v := range left {
		if _, ok := right[k]; !ok {
			e := diff.Entry{Key: k, Left: v}
			if opts.MaskSecrets {
				e = maskEntry(e)
			}
			result.OnlyInLeft = append(result.OnlyInLeft, e)
		}
	}

	for k, v := range right {
		if _, ok := left[k]; !ok {
			e := diff.Entry{Key: k, Right: v}
			if opts.MaskSecrets {
				e = maskEntry(e)
			}
			result.OnlyInRight = append(result.OnlyInRight, e)
		}
	}

	if opts.SortKeys {
		sort.Slice(result.OnlyInLeft, func(i, j int) bool {
			return result.OnlyInLeft[i].Key < result.OnlyInLeft[j].Key
		})
		sort.Slice(result.OnlyInRight, func(i, j int) bool {
			return result.OnlyInRight[i].Key < result.OnlyInRight[j].Key
		})
	}

	return result
}
