package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/subtlepseudonym/envcmp/internal/diff"
)

func TestUnique_NoOptions(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2", "C": "3"}

	result := Unique(left, right, UniqueOptions{})

	if len(result.OnlyInLeft) != 1 || result.OnlyInLeft[0].Key != "A" {
		t.Errorf("expected OnlyInLeft=[A], got %v", result.OnlyInLeft)
	}
	if len(result.OnlyInRight) != 1 || result.OnlyInRight[0].Key != "C" {
		t.Errorf("expected OnlyInRight=[C], got %v", result.OnlyInRight)
	}
}

func TestUnique_SortKeys(t *testing.T) {
	left := map[string]string{"Z": "z", "A": "a", "M": "m"}
	right := map[string]string{}

	result := Unique(left, right, UniqueOptions{SortKeys: true})

	want := []string{"A", "M", "Z"}
	for i, e := range result.OnlyInLeft {
		if e.Key != want[i] {
			t.Errorf("index %d: want key %q, got %q", i, want[i], e.Key)
		}
	}
}

func TestUnique_EmptyMaps(t *testing.T) {
	result := Unique(map[string]string{}, map[string]string{}, UniqueOptions{})

	if len(result.OnlyInLeft) != 0 || len(result.OnlyInRight) != 0 {
		t.Errorf("expected empty results for empty maps")
	}
}

func TestUnique_MaskSecrets(t *testing.T) {
	left := map[string]string{"SECRET_KEY": "hunter2"}
	right := map[string]string{}

	result := Unique(left, right, UniqueOptions{MaskSecrets: true})

	if len(result.OnlyInLeft) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.OnlyInLeft))
	}
	if result.OnlyInLeft[0].Left == "hunter2" {
		t.Errorf("expected secret value to be masked")
	}
}

func TestUnique_IdenticalMaps(t *testing.T) {
	m := map[string]string{"X": "1", "Y": "2"}
	result := Unique(m, m, UniqueOptions{})

	if len(result.OnlyInLeft) != 0 || len(result.OnlyInRight) != 0 {
		t.Errorf("expected no unique entries for identical maps")
	}
}

func TestUnique_ValuesPreserved(t *testing.T) {
	left := map[string]string{"ONLY_LEFT": "hello"}
	right := map[string]string{"ONLY_RIGHT": "world"}

	result := Unique(left, right, UniqueOptions{})

	wantLeft := []diff.Entry{{Key: "ONLY_LEFT", Left: "hello"}}
	wantRight := []diff.Entry{{Key: "ONLY_RIGHT", Right: "world"}}

	if diff := cmp.Diff(wantLeft, result.OnlyInLeft); diff != "" {
		t.Errorf("OnlyInLeft mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantRight, result.OnlyInRight); diff != "" {
		t.Errorf("OnlyInRight mismatch (-want +got):\n%s", diff)
	}
}
