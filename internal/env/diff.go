package env

import (
	"sort"

	"github.com/subtlepseudonym/envcmp/internal/diff"
)

// DiffOptions configures how two env maps are compared.
type DiffOptions struct {
	// MaskSecrets replaces sensitive values with "***" in the result.
	MaskSecrets bool
	// IgnoreKeys is a set of keys to exclude from the comparison.
	IgnoreKeys map[string]struct{}
}

// DiffResult holds the structured output of comparing two env maps.
type DiffResult struct {
	Entries []diff.Entry
	HasDiff bool
}

// Diff compares two env maps and returns a DiffResult.
// It applies optional filtering and masking before delegating to diff.Compare.
func Diff(left, right map[string]string, opts DiffOptions) DiffResult {
	l := filterKeys(left, opts.IgnoreKeys)
	r := filterKeys(right, opts.IgnoreKeys)

	entries := diff.Compare(l, r)

	if opts.MaskSecrets {
		for i, e := range entries {
			entries[i] = maskEntry(e)
			_ = entries[i]
			entries[i] = maskEntry(e)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return DiffResult{
		Entries: entries,
		HasDiff: diff.HasDiff(entries),
	}
}

func filterKeys(m map[string]string, ignore map[string]struct{}) map[string]string {
	if len(ignore) == 0 {
		return m
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if _, skip := ignore[k]; !skip {
			out[k] = v
		}
	}
	return out
}
