package env

import "strings"

// FilterOptions controls which keys are included in the filtered map.
type FilterOptions struct {
	// Prefix retains only keys that start with the given prefix.
	Prefix string

	// Keys retains only keys that appear in this list.
	// If empty, all keys are considered (subject to other options).
	Keys []string

	// Exclude removes keys that appear in this list.
	Exclude []string
}

// Filter returns a new map containing only the entries that satisfy opts.
// Keys in Exclude are always dropped, even if they match Prefix or Keys.
func Filter(env map[string]string, opts FilterOptions) map[string]string {
	excludeSet := toSet(opts.Exclude)
	keySet := toSet(opts.Keys)

	out := make(map[string]string, len(env))
	for k, v := range env {
		if excludeSet[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if len(keySet) > 0 && !keySet[k] {
			continue
		}
		out[k] = v
	}
	return out
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
