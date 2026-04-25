package env

import (
	"sort"
	"strings"
)

// GroupOptions controls how entries are grouped.
type GroupOptions struct {
	// PrefixSep is the separator used to detect prefix groups (default "_").
	PrefixSep string
	// MinGroupSize is the minimum number of keys to form a group (default 1).
	MinGroupSize int
}

// GroupResult holds a named group of key-value pairs.
type GroupResult struct {
	Prefix string
	Entries map[string]string
}

// Group partitions env entries by key prefix (e.g. DB_HOST, DB_PORT → "DB").
// Keys with no matching group are placed under the empty-string prefix.
func Group(env map[string]string, opts GroupOptions) []GroupResult {
	sep := opts.PrefixSep
	if sep == "" {
		sep = "_"
	}
	min := opts.MinGroupSize
	if min < 1 {
		min = 1
	}

	buckets := map[string]map[string]string{}
	for k, v := range env {
		prefix := ""
		if idx := strings.Index(k, sep); idx > 0 {
			prefix = k[:idx]
		}
		if buckets[prefix] == nil {
			buckets[prefix] = map[string]string{}
		}
		buckets[prefix][k] = v
	}

	// Merge small groups into the ungrouped bucket.
	for prefix, entries := range buckets {
		if prefix != "" && len(entries) < min {
			if buckets[""] == nil {
				buckets[""] = map[string]string{}
			}
			for k, v := range entries {
				buckets[""][k] = v
			}
			delete(buckets, prefix)
		}
	}

	prefixes := make([]string, 0, len(buckets))
	for p := range buckets {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	results := make([]GroupResult, 0, len(prefixes))
	for _, p := range prefixes {
		results = append(results, GroupResult{Prefix: p, Entries: buckets[p]})
	}
	return results
}
