package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/diff"
)

func TestRenderSnapshot_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	RenderSnapshot(nil, "prod.snap", false, &buf)
	out := buf.String()
	if !strings.Contains(out, "✔") {
		t.Errorf("expected success marker, got: %s", out)
	}
	if !strings.Contains(out, "prod.snap") {
		t.Errorf("expected snapshot path in output, got: %s", out)
	}
}

func TestRenderSnapshot_NoDiffColor(t *testing.T) {
	var buf bytes.Buffer
	RenderSnapshot(nil, "prod.snap", true, &buf)
	out := buf.String()
	if !strings.Contains(out, "\033[32m") {
		t.Errorf("expected green color code, got: %s", out)
	}
}

func TestRenderSnapshot_WithDiff(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Changed, LeftValue: "old-host", RightValue: "new-host"},
		{Key: "REMOVED_KEY", Status: diff.MissingInRight, LeftValue: "val"},
		{Key: "NEW_KEY", Status: diff.MissingInLeft, RightValue: "newval"},
	}
	var buf bytes.Buffer
	RenderSnapshot(results, "staging.snap", false, &buf)
	out := buf.String()

	if !strings.Contains(out, "⚠") {
		t.Errorf("expected warning marker, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' label in output")
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' label in output")
	}
	if !strings.Contains(out, "changed") {
		t.Errorf("expected 'changed' label in output")
	}
}

func TestRenderSnapshot_SortedOutput(t *testing.T) {
	results := []diff.Result{
		{Key: "Z_KEY", Status: diff.Changed, LeftValue: "a", RightValue: "b"},
		{Key: "A_KEY", Status: diff.Changed, LeftValue: "c", RightValue: "d"},
	}
	var buf bytes.Buffer
	RenderSnapshot(results, "snap", false, &buf)
	out := buf.String()

	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}
