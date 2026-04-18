// Package filter provides key filtering utilities for env maps.
package filter

import "strings"

// Options holds filtering configuration.
type Options struct {
	Prefix  string
	Keys    []string
	Exclude []string
}

// Apply filters an env map according to the given Options.
// If Prefix is set, only keys with that prefix are retained.
// If Keys is non-empty, only those keys are retained.
// Exclude removes specific keys from the result.
func Apply(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))

	allowSet := toSet(opts.Keys)
	excludeSet := toSet(opts.Exclude)

	for k, v := range env {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if len(allowSet) > 0 {
			if _, ok := allowSet[k]; !ok {
				continue
			}
		}
		if _, excluded := excludeSet[k]; excluded {
			continue
		}
		result[k] = v
	}
	return result
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
