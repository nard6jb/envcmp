package env

import (
	"sort"

	"github.com/jasonlovesdoggo/envcmp/internal/mask"
)

// SubtractOptions controls the behaviour of Subtract.
type SubtractOptions struct {
	// MaskSecrets replaces sensitive values with "***" in the returned entries.
	MaskSecrets bool
}

// SubtractEntry represents a key/value pair that was present in left but
// absent from right (or whose value differed, depending on usage).
type SubtractEntry struct {
	Key   string
	Value string
}

// Subtract returns all entries that exist in left but are NOT present in
// right (keyed by key name). The result is sorted alphabetically by key.
func Subtract(left, right map[string]string, opts SubtractOptions) []SubtractEntry {
	rightKeys := toSubtractSet(right)

	var result []SubtractEntry
	for k, v := range left {
		if rightKeys[k] {
			continue
		}
		if opts.MaskSecrets && mask.IsSensitive(k) {
			v = mask.MaskValue(v)
		}
		result = append(result, SubtractEntry{Key: k, Value: v})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
}

func toSubtractSet(m map[string]string) map[string]bool {
	s := make(map[string]bool, len(m))
	for k := range m {
		s[k] = true
	}
	return s
}
