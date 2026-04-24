package env

import (
	"fmt"

	"github.com/subtlepseudonym/envcmp/internal/merge"
)

// MergeOptions controls how environment maps are combined.
type MergeOptions struct {
	Strategy merge.Strategy
	MaskSecrets bool
}

// MergeResult holds the output of merging two env maps.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []MergeConflict
}

// MergeConflict describes a key present in both maps with differing values.
type MergeConflict struct {
	Key    string
	Left   string
	Right  string
	Chosen string
}

// Merge combines left and right env maps according to the provided options.
// It delegates conflict resolution to the merge package and surfaces
// per-key conflict details for reporting.
func Merge(left, right map[string]string, opts MergeOptions) (MergeResult, error) {
	if left == nil {
		return MergeResult{}, fmt.Errorf("env.Merge: left map must not be nil")
	}
	if right == nil {
		return MergeResult{}, fmt.Errorf("env.Merge: right map must not be nil")
	}

	merged, conflicts := merge.Merge(left, right, opts.Strategy)

	result := MergeResult{
		Merged:    merged,
		Conflicts: make([]MergeConflict, 0, len(conflicts)),
	}

	for _, c := range conflicts {
		lv := c.Left
		rv := c.Right
		chosen := c.Chosen
		if opts.MaskSecrets {
			lv = maskEntry(c.Key, lv)
			rv = maskEntry(c.Key, rv)
			chosen = maskEntry(c.Key, chosen)
		}
		result.Conflicts = append(result.Conflicts, MergeConflict{
			Key:    c.Key,
			Left:   lv,
			Right:  rv,
			Chosen: chosen,
		})
	}

	return result, nil
}
