// Package env provides utilities for expanding and normalizing environment variable maps.
package env

import (
	"strings"
)

// NormalizeOptions controls how normalization is applied to an env map.
type NormalizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// RemoveEmpty drops entries with empty values after trimming.
	RemoveEmpty bool
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// LowercaseKeys converts all keys to lowercase (takes precedence over UppercaseKeys).
	LowercaseKeys bool
}

// Normalize applies the given options to a copy of the input map and returns
// the normalized result. The original map is never mutated.
func Normalize(input map[string]string, opts NormalizeOptions) map[string]string {
	out := make(map[string]string, len(input))

	for k, v := range input {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}

		if opts.RemoveEmpty && v == "" {
			continue
		}

		newKey := k
		switch {
		case opts.LowercaseKeys:
			newKey = strings.ToLower(k)
		case opts.UppercaseKeys:
			newKey = strings.ToUpper(k)
		}

		out[newKey] = v
	}

	return out
}
