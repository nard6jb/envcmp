package env

import "sort"

// FlattenOptions controls how a nested prefix map is flattened into a single env map.
type FlattenOptions struct {
	// Separator is placed between prefix levels when constructing keys.
	// Defaults to "_" if empty.
	Separator string

	// UppercaseKeys forces all resulting keys to uppercase.
	UppercaseKeys bool

	// Prefix is prepended to every resulting key.
	Prefix string
}

// FlattenResult holds a single flattened key-value pair.
type FlattenResult struct {
	Key      string
	Value    string
	Original string // the source key before transformation
}

// Flatten takes a map of maps (e.g. grouped env sections) and merges them
// into a single key→value map, applying separator and optional transforms.
// Duplicate keys from later groups overwrite earlier ones.
func Flatten(groups map[string]map[string]string, opts FlattenOptions) (map[string]string, []FlattenResult) {
	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	result := make(map[string]string)
	var entries []FlattenResult

	// Sort group names for deterministic output.
	groupNames := make([]string, 0, len(groups))
	for g := range groups {
		groupNames = append(groupNames, g)
	}
	sort.Strings(groupNames)

	for _, group := range groupNames {
		kvs := groups[group]
		keys := make([]string, 0, len(kvs))
		for k := range kvs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			original := k
			composed := group + sep + k
			if opts.Prefix != "" {
				composed = opts.Prefix + sep + composed
			}
			if opts.UppercaseKeys {
				composed = toUpperASCII(composed)
			}
			result[composed] = kvs[k]
			entries = append(entries, FlattenResult{
				Key:      composed,
				Value:    kvs[k],
				Original: original,
			})
		}
	}

	return result, entries
}

func toUpperASCII(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
