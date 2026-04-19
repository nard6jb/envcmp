package snapshot

import (
	"fmt"
	"sort"

	"github.com/yourusername/envcmp/internal/diff"
)

// DiffResult holds the result of comparing current env against a snapshot.
type DiffResult struct {
	SnapshotFile string
	EnvFile      string
	Diffs        []diff.Entry
	HasDiff      bool
}

// CompareEnvAgainstSnapshot loads a snapshot from path and compares it
// against the provided env map, returning a structured DiffResult.
func CompareEnvAgainstSnapshot(snapshotPath string, current map[string]string) (*DiffResult, error) {
	snap, err := Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: load %q: %w", snapshotPath, err)
	}

	entries := diff.Compare(snap, current)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &DiffResult{
		SnapshotFile: snapshotPath,
		Diffs:        entries,
		HasDiff:      diff.HasDiff(entries),
	}, nil
}
