package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/env"
)

func TestRenderPatch_NoResults(t *testing.T) {
	var buf bytes.Buffer
	RenderPatchTo(&buf, nil, false)
	if !strings.Contains(buf.String(), "No patch operations") {
		t.Errorf("expected no-op message, got: %s", buf.String())
	}
}

func TestRenderPatch_Added(t *testing.T) {
	results := []env.PatchResult{
		{Key: "TIMEOUT", NewValue: "30s", Added: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "ADD") {
		t.Errorf("expected ADD in output, got: %s", out)
	}
	if !strings.Contains(out, "TIMEOUT") {
		t.Errorf("expected key TIMEOUT in output, got: %s", out)
	}
}

func TestRenderPatch_Deleted(t *testing.T) {
	results := []env.PatchResult{
		{Key: "DB", OldValue: "mydb", Deleted: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "DELETE") {
		t.Errorf("expected DELETE in output, got: %s", out)
	}
	if !strings.Contains(out, "mydb") {
		t.Errorf("expected old value in output, got: %s", out)
	}
}

func TestRenderPatch_Changed(t *testing.T) {
	results := []env.PatchResult{
		{Key: "PORT", OldValue: "5432", NewValue: "3306", Changed: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "CHANGE") {
		t.Errorf("expected CHANGE in output, got: %s", out)
	}
	if !strings.Contains(out, "5432") || !strings.Contains(out, "3306") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}

func TestRenderPatch_Skipped(t *testing.T) {
	results := []env.PatchResult{
		{Key: "GHOST", Skipped: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "SKIP") {
		t.Errorf("expected SKIP in output, got: %s", out)
	}
}

func TestRenderPatch_SortedOutput(t *testing.T) {
	results := []env.PatchResult{
		{Key: "ZEBRA", NewValue: "z", Added: true},
		{Key: "ALPHA", NewValue: "a", Added: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, false)
	out := buf.String()
	alphaIdx := strings.Index(out, "ALPHA")
	zebraIdx := strings.Index(out, "ZEBRA")
	if alphaIdx > zebraIdx {
		t.Error("expected ALPHA to appear before ZEBRA in sorted output")
	}
}

func TestRenderPatch_ColorAdded(t *testing.T) {
	results := []env.PatchResult{
		{Key: "NEW", NewValue: "val", Added: true},
	}
	var buf bytes.Buffer
	RenderPatchTo(&buf, results, true)
	if !strings.Contains(buf.String(), "\x1b[") {
		t.Error("expected ANSI color codes in color output")
	}
}
