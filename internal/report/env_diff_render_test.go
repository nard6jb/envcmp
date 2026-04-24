package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/subtlepseudonym/envcmp/internal/diff"
	"github.com/subtlepseudonym/envcmp/internal/env"
)

func TestRenderEnvDiff_NoDiff(t *testing.T) {
	result := env.DiffResult{HasDiff: false, Entries: nil}
	var buf bytes.Buffer
	RenderEnvDiffTo(&buf, result, false)
	if !strings.Contains(buf.String(), "no differences") {
		t.Errorf("expected no-diff message, got %q", buf.String())
	}
}

func TestRenderEnvDiff_NoDiffColor(t *testing.T) {
	result := env.DiffResult{HasDiff: false}
	var buf bytes.Buffer
	RenderEnvDiffTo(&buf, result, true)
	if !strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected ANSI codes in color output")
	}
}

func TestRenderEnvDiff_WithChanges(t *testing.T) {
	result := env.DiffResult{
		HasDiff: true,
		Entries: []diff.Entry{
			{Key: "FOO", Left: "old", Right: "new", Status: diff.Changed},
		},
	}
	var buf bytes.Buffer
	RenderEnvDiffTo(&buf, result, false)
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output, got %q", out)
	}
}

func TestRenderEnvDiff_MissingRight(t *testing.T) {
	result := env.DiffResult{
		HasDiff: true,
		Entries: []diff.Entry{
			{Key: "BAR", Left: "val", Right: "", Status: diff.MissingRight},
		},
	}
	var buf bytes.Buffer
	RenderEnvDiffTo(&buf, result, false)
	if !strings.Contains(buf.String(), "BAR") {
		t.Errorf("expected BAR in output")
	}
}

func TestRenderEnvDiff_MissingLeft(t *testing.T) {
	result := env.DiffResult{
		HasDiff: true,
		Entries: []diff.Entry{
			{Key: "BAZ", Left: "", Right: "val", Status: diff.MissingLeft},
		},
	}
	var buf bytes.Buffer
	RenderEnvDiffTo(&buf, result, false)
	if !strings.Contains(buf.String(), "BAZ") {
		t.Errorf("expected BAZ in output")
	}
}
