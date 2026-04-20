// Package env provides utilities for expanding environment variable
// references from a base map into a target map, supporting both
// OS-level fallback and strict-only resolution modes.
package env

import (
	"fmt"
	"os"
	"strings"
)

// ExpandOptions controls how variable expansion is performed.
type ExpandOptions struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// StrictMode causes an error if any reference cannot be resolved.
	StrictMode bool
}

// ExpandResult holds the expanded map and any unresolved keys.
type ExpandResult struct {
	Expanded   map[string]string
	Unresolved []string
}

// Expand resolves ${VAR} and $VAR references in values of src using base as
// the lookup table. If opts.FallbackToOS is true, missing keys are looked up
// in the process environment. If opts.StrictMode is true, any unresolved
// reference returns an error.
func Expand(src map[string]string, base map[string]string, opts ExpandOptions) (ExpandResult, error) {
	result := make(map[string]string, len(src))
	var unresolved []string

	lookup := func(key string) (string, bool) {
		if v, ok := base[key]; ok {
			return v, true
		}
		if v, ok := src[key]; ok {
			return v, true
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(key); ok {
				return v, true
			}
		}
		return "", false
	}

	for k, v := range src {
		expanded, missing := expandValue(v, lookup)
		result[k] = expanded
		unresolved = append(unresolved, missing...)
	}

	if opts.StrictMode && len(unresolved) > 0 {
		return ExpandResult{}, fmt.Errorf("unresolved references: %s", strings.Join(dedupe(unresolved), ", "))
	}

	return ExpandResult{Expanded: result, Unresolved: dedupe(unresolved)}, nil
}

func expandValue(val string, lookup func(string) (string, bool)) (string, []string) {
	var missing []string
	result := os.Expand(val, func(key string) string {
		if v, ok := lookup(key); ok {
			return v
		}
		missing = append(missing, key)
		return ""
	})
	return result, missing
}

func dedupe(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	out := s[:0:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
