package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/envcmp/internal/lint"
	"github.com/yourusername/envcmp/internal/report"
)

func TestRenderLint_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	ok := report.RenderLint(&buf, nil, false)
	if !ok {
		t.Error("expected ok=true when no issues")
	}
	if !strings.Contains(buf.String(), "no lint issues") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRenderLint_WithIssues(t *testing.T) {
	issues := []lint.Issue{
		{Line: 1, Key: "FOO", Message: "empty value"},
		{Line: 2, Key: "BAR", Message: "duplicate key (first seen at line 1)"},
	}
	var buf bytes.Buffer
	ok := report.RenderLint(&buf, issues, false)
	if ok {
		t.Error("expected ok=false when issues present")
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Error("expected FOO in output")
	}
	if !strings.Contains(out, "BAR") {
		t.Error("expected BAR in output")
	}
}

func TestRenderLint_ColorOutput(t *testing.T) {
	issues := []lint.Issue{{Line: 1, Key: "X", Message: "empty value"}}
	var buf bytes.Buffer
	report.RenderLint(&buf, issues, true)
	if !strings.Contains(buf.String(), "\033[") {
		t.Error("expected ANSI codes in color output")
	}
}

func TestRenderLint_NoIssuesColor(t *testing.T) {
	var buf bytes.Buffer
	report.RenderLint(&buf, []lint.Issue{}, true)
	if !strings.Contains(buf.String(), "\033[") {
		t.Error("expected ANSI codes in color ok output")
	}
}
