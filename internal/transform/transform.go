// Package transform provides utilities for transforming env maps,
// such as renaming keys, uppercasing, or applying prefix/suffix changes.
package transform

import "strings"

// Options controls how the transformation is applied.
type Options struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// LowercaseKeys converts all keys to lowercase.
	LowercaseKeys bool
	// StripPrefix removes a leading prefix from each key, if present.
	StripPrefix string
	// AddPrefix prepends a string to each key.
	AddPrefix string
	// RenameMap maps old key names to new key names.
	RenameMap map[string]string
}

// Apply returns a new map with keys transformed according to opts.
// If a key appears in RenameMap, the rename takes precedence over
// prefix/case transformations.
func Apply(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		newKey := transformKey(k, opts)
		result[newKey] = v
	}
	return result
}

func transformKey(key string, opts Options) string {
	// Explicit rename takes highest priority.
	if opts.RenameMap != nil {
		if renamed, ok := opts.RenameMap[key]; ok {
			return renamed
		}
	}

	result := key

	if opts.StripPrefix != "" {
		result = strings.TrimPrefix(result, opts.StripPrefix)
	}

	if opts.AddPrefix != "" {
		result = opts.AddPrefix + result
	}

	if opts.UppercaseKeys {
		result = strings.ToUpper(result)
	} else if opts.LowercaseKeys {
		result = strings.ToLower(result)
	}

	return result
}
