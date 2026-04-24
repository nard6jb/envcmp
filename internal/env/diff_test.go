package env

import (
	"testing"

	"github.com/subtlepseudonym/envcmp/internal/diff"
)

func TestDiff_NoDiff(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}

	result := Diff(left, right, DiffOptions{})

	if result.HasDiff {
		t.Errorf("expected no diff, got diff")
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	left := map[string]string{"A": "old", "B": "same"}
	right := map[string]string{"A": "new", "B": "same"}

	result := Diff(left, right, DiffOptions{})

	if !result.HasDiff {
		t.Errorf("expected diff, got none")
	}
	if len(result.Entries) != 1 || result.Entries[0].Key != "A" {
		t.Errorf("expected entry for key A")
	}
}

func TestDiff_IgnoreKeys(t *testing.T) {
	left := map[string]string{"A": "1", "IGNORE": "x"}
	right := map[string]string{"A": "1", "IGNORE": "y"}

	opts := DiffOptions{IgnoreKeys: map[string]struct{}{"IGNORE": {}}}
	result := Diff(left, right, opts)

	if result.HasDiff {
		t.Errorf("expected no diff after ignoring key")
	}
}

func TestDiff_MaskSecrets(t *testing.T) {
	left := map[string]string{"API_SECRET": "abc123"}
	right := map[string]string{"API_SECRET": "xyz789"}

	result := Diff(left, right, DiffOptions{MaskSecrets: true})

	if !result.HasDiff {
		t.Fatalf("expected diff")
	}
	e := result.Entries[0]
	if e.Left == "abc123" || e.Right == "xyz789" {
		t.Errorf("expected values to be masked, got left=%q right=%q", e.Left, e.Right)
	}
}

func TestDiff_SortedEntries(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "1", "M": "1"}
	right := map[string]string{"Z": "2", "A": "2", "M": "2"}

	result := Diff(left, right, DiffOptions{})

	keys := make([]string, len(result.Entries))
	for i, e := range result.Entries {
		keys[i] = e.Key
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %q got %q", i, k, keys[i])
		}
	}
}

func TestDiff_MissingInRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}

	result := Diff(left, right, DiffOptions{})

	if !result.HasDiff {
		t.Errorf("expected diff for missing key")
	}
	var found bool
	for _, e := range result.Entries {
		if e.Key == "B" && e.Status == diff.MissingRight {
			found = true
		}
	}
	if !found {
		t.Errorf("expected MissingRight entry for key B")
	}
}
