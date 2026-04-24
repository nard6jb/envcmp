package env

import (
	"fmt"
	"sort"
	"strings"
)

// CloneOptions controls how an env map is cloned and transformed.
type CloneOptions struct {
	// RenamePrefix replaces OldPrefix with NewPrefix on all matching keys.
	OldPrefix string
	NewPrefix string
	// OverrideKeys merges additional key=value pairs after cloning.
	OverrideKeys map[string]string
	// DropKeys removes specific keys from the cloned result.
	DropKeys []string
}

// Clone produces a deep copy of src, applying optional transformations.
func Clone(src map[string]string, opts CloneOptions) (map[string]string, []string) {
	out := make(map[string]string, len(src))
	var warnings []string

	dropSet := make(map[string]bool, len(opts.DropKeys))
	for _, k := range opts.DropKeys {
		dropSet[k] = true
	}

	for k, v := range src {
		if dropSet[k] {
			continue
		}
		newKey := k
		if opts.OldPrefix != "" && strings.HasPrefix(k, opts.OldPrefix) {
			newKey = opts.NewPrefix + strings.TrimPrefix(k, opts.OldPrefix)
		}
		if _, exists := out[newKey]; exists {
			warnings = append(warnings, fmt.Sprintf("key collision after rename: %s", newKey))
		}
		out[newKey] = v
	}

	for k, v := range opts.OverrideKeys {
		out[k] = v
	}

	sort.Strings(warnings)
	return out, warnings
}
