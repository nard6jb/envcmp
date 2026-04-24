package env

import (
	"testing"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}
	out, results := Promote(src, dst, PromoteOptions{})
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected output map: %v", out)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestPromote_OnlyKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, results := Promote(src, dst, PromoteOptions{OnlyKeys: []string{"A", "C"}})
	if _, ok := out["B"]; ok {
		t.Error("B should not be promoted")
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Errorf("unexpected output: %v", out)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestPromote_NoOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, results := Promote(src, dst, PromoteOptions{Overwrite: false})
	if out["A"] != "old" {
		t.Errorf("expected old value to be preserved, got %q", out["A"])
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected result to be skipped")
	}
}

func TestPromote_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, results := Promote(src, dst, PromoteOptions{Overwrite: true})
	if out["A"] != "new" {
		t.Errorf("expected new value, got %q", out["A"])
	}
	if len(results) != 1 || results[0].Skipped {
		t.Error("expected result to not be skipped")
	}
	if results[0].OldValue != "old" || results[0].NewValue != "new" {
		t.Errorf("unexpected result values: %+v", results[0])
	}
}

func TestPromote_SkipEmpty(t *testing.T) {
	src := map[string]string{"A": "", "B": "value"}
	dst := map[string]string{}
	out, results := Promote(src, dst, PromoteOptions{SkipEmpty: true})
	if _, ok := out["A"]; ok {
		t.Error("A should be skipped due to empty value")
	}
	if out["B"] != "value" {
		t.Errorf("expected B to be promoted, got %v", out)
	}
	skipped := 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		}
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped result, got %d", skipped)
	}
}

func TestPromote_EmptySource(t *testing.T) {
	out, results := Promote(map[string]string{}, map[string]string{"X": "1"}, PromoteOptions{})
	if out["X"] != "1" {
		t.Error("dst key should be preserved")
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}
