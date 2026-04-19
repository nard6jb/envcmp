package snapshot_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourusername/envcmp/internal/snapshot"
)

func writeSnapshot(t *testing.T, data map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp("", "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestCompareEnvAgainstSnapshot_NoDiff(t *testing.T) {
	snap := map[string]string{"APP_ENV": "prod", "PORT": "8080"}
	path := writeSnapshot(t, snap)

	result, err := snapshot.CompareEnvAgainstSnapshot(path, snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDiff {
		t.Error("expected no diff")
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(result.Diffs))
	}
}

func TestCompareEnvAgainstSnapshot_DetectsDiff(t *testing.T) {
	snap := map[string]string{"APP_ENV": "prod", "PORT": "8080"}
	path := writeSnapshot(t, snap)

	current := map[string]string{"APP_ENV": "staging", "PORT": "8080", "NEW_KEY": "value"}
	result, err := snapshot.CompareEnvAgainstSnapshot(path, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDiff {
		t.Error("expected diff")
	}
	if len(result.Diffs) == 0 {
		t.Error("expected non-empty diffs")
	}
}

func TestCompareEnvAgainstSnapshot_FileNotFound(t *testing.T) {
	_, err := snapshot.CompareEnvAgainstSnapshot("/nonexistent/snap.json", map[string]string{})
	if err == nil {
		t.Error("expected error for missing snapshot file")
	}
}
