package merge

import "testing"

func maps(pairs ...map[string]string) []map[string]string { return pairs }

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3"}
	res := Merge(maps(a, b), StrategyFirst)
	if res.Merged["A"] != "1" || res.Merged["C"] != "3" {
		t.Fatalf("unexpected merged map: %v", res.Merged)
	}
	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res := Merge(maps(a, b), StrategyFirst)
	if res.Merged["KEY"] != "first" {
		t.Fatalf("expected 'first', got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res := Merge(maps(a, b), StrategyLast)
	if res.Merged["KEY"] != "second" {
		t.Fatalf("expected 'second', got %q", res.Merged["KEY"])
	}
}

func TestMerge_MultipleConflicts(t *testing.T) {
	a := map[string]string{"X": "a", "Y": "a"}
	b := map[string]string{"X": "b", "Y": "b"}
	c := map[string]string{"X": "c"}
	res := Merge(maps(a, b, c), StrategyLast)
	if res.Merged["X"] != "c" {
		t.Fatalf("expected 'c', got %q", res.Merged["X"])
	}
	if len(res.Conflicts) != 2 {
		t.Fatalf("expected 2 conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res := Merge([]map[string]string{}, StrategyFirst)
	if len(res.Merged) != 0 {
		t.Fatal("expected empty merged map")
	}
}
