package env

import (
	"testing"

	"github.com/jasonuc/envcmp/internal/diff"
)

func entries() []diff.Entry {
	return []diff.Entry{
		{Key: "ZEBRA", LeftVal: "1", RightVal: "1"},
		{Key: "ALPHA", LeftVal: "", RightVal: "new"},
		{Key: "MANGO", LeftVal: "old", RightVal: "new"},
		{Key: "BETA", LeftVal: "gone", RightVal: ""},
	}
}

func TestSortEntries_ByKeyAsc(t *testing.T) {
	out := SortEntries(entries(), SortOptions{ByKey: true})
	keys := []string{out[0].Key, out[1].Key, out[2].Key, out[3].Key}
	want := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("pos %d: got %q, want %q", i, keys[i], k)
		}
	}
}

func TestSortEntries_ByKeyDesc(t *testing.T) {
	out := SortEntries(entries(), SortOptions{ByKey: true, Reverse: true})
	if out[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %q", out[0].Key)
	}
	if out[3].Key != "ALPHA" {
		t.Errorf("expected ALPHA last, got %q", out[3].Key)
	}
}

func TestSortEntries_GroupByStatus(t *testing.T) {
	out := SortEntries(entries(), SortOptions{GroupByStatus: true})
	// rank 0 = missing-left (ALPHA), rank 1 = missing-right (BETA),
	// rank 2 = changed (MANGO),      rank 3 = same (ZEBRA)
	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if out[i].Key != k {
			t.Errorf("pos %d: got %q, want %q", i, out[i].Key, k)
		}
	}
}

func TestSortEntries_GroupByStatusReverse(t *testing.T) {
	out := SortEntries(entries(), SortOptions{GroupByStatus: true, Reverse: true})
	// rank order reversed: same(3) > changed(2) > missing-right(1) > missing-left(0)
	if out[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first in reverse group sort, got %q", out[0].Key)
	}
}

func TestSortEntries_DoesNotMutateInput(t *testing.T) {
	in := entries()
	orig := in[0].Key
	_ = SortEntries(in, SortOptions{ByKey: true})
	if in[0].Key != orig {
		t.Error("SortEntries must not mutate the input slice")
	}
}
