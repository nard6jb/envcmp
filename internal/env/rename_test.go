package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRename_NoOptions(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Rename(src, RenameOptions{})
	if diff := cmp.Diff(src, res.Out); diff != "" {
		t.Errorf("unexpected diff: %s", diff)
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
}

func TestRename_BasicMapping(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value", "KEEP": "yes"}
	opts := RenameOptions{Mapping: map[string]string{"OLD_KEY": "NEW_KEY"}}
	res := Rename(src, opts)

	if res.Out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", res.Out["NEW_KEY"])
	}
	if _, exists := res.Out["OLD_KEY"]; exists {
		t.Error("OLD_KEY should not exist in output")
	}
	if res.Out["KEEP"] != "yes" {
		t.Error("KEEP should be preserved")
	}
	if res.Renamed["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected Renamed[OLD_KEY]=NEW_KEY, got %q", res.Renamed["OLD_KEY"])
	}
}

func TestRename_DropUnmapped(t *testing.T) {
	src := map[string]string{"OLD": "v", "EXTRA": "drop_me"}
	opts := RenameOptions{
		Mapping:      map[string]string{"OLD": "NEW"},
		DropUnmapped: true,
	}
	res := Rename(src, opts)

	if _, exists := res.Out["EXTRA"]; exists {
		t.Error("EXTRA should have been dropped")
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "EXTRA" {
		t.Errorf("expected Dropped=[EXTRA], got %v", res.Dropped)
	}
	if res.Out["NEW"] != "v" {
		t.Errorf("expected NEW=v, got %q", res.Out["NEW"])
	}
}

func TestRename_SkipsConflict(t *testing.T) {
	// NEW_KEY already exists; rename should be skipped.
	src := map[string]string{"OLD": "old_val", "NEW": "existing"}
	opts := RenameOptions{Mapping: map[string]string{"OLD": "NEW"}}
	res := Rename(src, opts)

	if res.Out["NEW"] != "existing" {
		t.Errorf("expected existing value preserved, got %q", res.Out["NEW"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "OLD" {
		t.Errorf("expected Skipped=[OLD], got %v", res.Skipped)
	}
}

func TestRename_MissingSourceKey(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	opts := RenameOptions{Mapping: map[string]string{"NONEXISTENT": "TARGET"}}
	res := Rename(src, opts)

	if _, exists := res.Out["TARGET"]; exists {
		t.Error("TARGET should not be created for missing source key")
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
}
