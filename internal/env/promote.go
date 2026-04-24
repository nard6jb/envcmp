package env

import (
	"fmt"
	"sort"
)

// PromoteOptions controls how keys are promoted between environments.
type PromoteOptions struct {
	// OnlyKeys restricts promotion to these specific keys. If empty, all keys are promoted.
	OnlyKeys []string
	// Overwrite allows values in the destination to be overwritten.
	Overwrite bool
	// SkipEmpty skips keys whose source value is empty.
	SkipEmpty bool
}

// PromoteResult describes the outcome of a single key promotion.
type PromoteResult struct {
	Key      string
	OldValue string
	NewValue string
	Skipped  bool
	Reason   string
}

// Promote copies keys from src into dst according to opts.
// It returns a slice of PromoteResult describing what happened to each key.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, []PromoteResult) {
	allow := toSet(opts.OnlyKeys)

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []PromoteResult
	for _, k := range keys {
		if len(allow) > 0 && !allow[k] {
			continue
		}
		srcVal := src[k]
		if opts.SkipEmpty && srcVal == "" {
			results = append(results, PromoteResult{Key: k, Skipped: true, Reason: "empty value"})
			continue
		}
		existing, exists := dst[k]
		if exists && !opts.Overwrite {
			results = append(results, PromoteResult{Key: k, OldValue: existing, NewValue: srcVal, Skipped: true, Reason: fmt.Sprintf("already exists in destination")})
			continue
		}
		out[k] = srcVal
		results = append(results, PromoteResult{Key: k, OldValue: existing, NewValue: srcVal, Skipped: false})
	}
	return out, results
}
