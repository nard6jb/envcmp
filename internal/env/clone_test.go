package env

import (
	"testing"
)

func TestClone_NoDiff(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, warns := Clone(src, CloneOptions{})
	if len(warns) != 0 {
		t.Fatalf("expected no warnings, got %v", warns)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
	// ensure deep copy
	out["FOO"] = "changed"
	if src["FOO"] != "bar" {
		t.Error("Clone mutated original map")
	}
}

func TestClone_RenamePrefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "val"}
	out, warns := Clone(src, CloneOptions{OldPrefix: "APP_", NewPrefix: "SVC_"})
	if len(warns) != 0 {
		t.Fatalf("unexpected warnings: %v", warns)
	}
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %v", out["SVC_HOST"])
	}
	if out["SVC_PORT"] != "8080" {
		t.Errorf("expected SVC_PORT=8080, got %v", out["SVC_PORT"])
	}
	if out["OTHER"] != "val" {
		t.Errorf("expected OTHER=val, got %v", out["OTHER"])
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("old key APP_HOST should not exist after rename")
	}
}

func TestClone_DropKeys(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2", "SECRET": "s"}
	out, _ := Clone(src, CloneOptions{DropKeys: []string{"SECRET"}})
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should have been dropped")
	}
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestClone_OverrideKeys(t *testing.T) {
	src := map[string]string{"HOST": "old", "PORT": "80"}
	out, _ := Clone(src, CloneOptions{OverrideKeys: map[string]string{"HOST": "new", "EXTRA": "yes"}})
	if out["HOST"] != "new" {
		t.Errorf("expected HOST=new, got %s", out["HOST"])
	}
	if out["EXTRA"] != "yes" {
		t.Errorf("expected EXTRA=yes, got %s", out["EXTRA"])
	}
	if out["PORT"] != "80" {
		t.Errorf("expected PORT=80, got %s", out["PORT"])
	}
}

func TestClone_EmptySource(t *testing.T) {
	out, warns := Clone(map[string]string{}, CloneOptions{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
	if len(warns) != 0 {
		t.Errorf("expected no warnings, got %v", warns)
	}
}
