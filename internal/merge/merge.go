// Package merge provides functionality to merge multiple .env files
// with configurable conflict resolution strategies.
package merge

import "github.com/your-org/envcmp/internal/envfile"

// Strategy defines how conflicts are resolved when merging.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
)

// Result holds the merged environment and metadata about conflicts.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
}

// Conflict describes a key that appeared in more than one source file.
type Conflict struct {
	Key    string
	Values []string // one per source, in order
}

// Merge combines multiple parsed env maps using the given strategy.
// Source maps are processed in the order provided.
func Merge(sources []map[string]string, strategy Strategy) Result {
	merged := make(map[string]string)
	seen := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			seen[k] = append(seen[k], v)
		}
	}

	var conflicts []Conflict
	for k, vals := range seen {
		if len(vals) > 1 {
			conflicts = append(conflicts, Conflict{Key: k, Values: vals})
		}
		switch strategy {
		case StrategyLast:
			merged[k] = vals[len(vals)-1]
		default:
			merged[k] = vals[0]
		}
	}

	return Result{Merged: merged, Conflicts: conflicts}
}

// MergeFiles parses and merges the given file paths.
func MergeFiles(paths []string, strategy Strategy) (Result, error) {
	var sources []map[string]string
	for _, p := range paths {
		entries, err := envfile.Parse(p)
		if err != nil {
			return Result{}, err
		}
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		sources = append(sources, m)
	}
	return Merge(sources, strategy), nil
}
