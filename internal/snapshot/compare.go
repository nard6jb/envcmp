package snapshot

import (
	"fmt"

	"github.com/user/envcmp/internal/diff"
)

// CompareResult holds the outcome of comparing a live env map against a snapshot.
type CompareResult struct {
	SnapshotFile string
	Diffs        []diff.Entry
	HasDiff      bool
}

// CompareAgainstSnapshot loads a snapshot from path and compares it against
// the provided live environment map. Returns a CompareResult or an error.
func CompareAgainstSnapshot(snapshotPath string, live map[string]string) (CompareResult, error) {
	snap, err := Load(snapshotPath)
	if err != nil {
		return CompareResult{}, fmt.Errorf("loading snapshot %q: %w", snapshotPath, err)
	}

	entries := diff.Compare(snap, live)
	return CompareResult{
		SnapshotFile: snapshotPath,
		Diffs:        entries,
		HasDiff:      diff.HasDiff(entries),
	}, nil
}
