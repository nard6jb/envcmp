package report_test

import (
	"strings"
	"testing"

	"github.com/your-org/envcmp/internal/report"
	"github.com/your-org/envcmp/internal/schema"
)

func TestRenderSchema_NoIssues(t *testing.T) {
	var buf strings.Builder
	report.RenderSchema(&buf, []schema.Issue{}, false)
	out := buf.String()
	if !strings.Contains(out, "no schema issues") {
		t.Errorf("expected no-issues message, got: %q", out)
	}
}

func TestRenderSchema_WithIssues(t *testing.T) {
	issues := []schema.Issue{
		{Key: "PORT", Message: "expected type int, got string"},
		{Key: "HOST", Message: "required key missing"},
	}
	var buf strings.Builder
	report.RenderSchema(&buf, issues, false)
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %q", out)
	}
	if !strings.Contains(out, "expected type int") {
		t.Errorf("expected type error message in output, got: %q", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %q", out)
	}
	if !strings.Contains(out, "required key missing") {
		t.Errorf("expected missing key message in output, got: %q", out)
	}
}

func TestRenderSchema_IssueCount(t *testing.T) {
	issues := []schema.Issue{
		{Key: "A", Message: "bad type"},
		{Key: "B", Message: "pattern mismatch"},
		{Key: "C", Message: "required key missing"},
	}
	var buf strings.Builder
	report.RenderSchema(&buf, issues, false)
	out := buf.String()
	if !strings.Contains(out, "3 schema issue") {
		t.Errorf("expected issue count summary, got: %q", out)
	}
}

func TestRenderSchema_ColorOutput(t *testing.T) {
	issues := []schema.Issue{
		{Key: "SECRET", Message: "required key missing"},
	}
	var buf strings.Builder
	report.RenderSchema(&buf, issues, true)
	out := buf.String()
	// ANSI escape codes should be present
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected ANSI color codes in output, got: %q", out)
	}
}

func TestRenderSchema_SortedKeys(t *testing.T) {
	issues := []schema.Issue{
		{Key: "ZEBRA", Message: "bad type"},
		{Key: "ALPHA", Message: "pattern mismatch"},
		{Key: "MANGO", Message: "required key missing"},
	}
	var buf strings.Builder
	report.RenderSchema(&buf, issues, false)
	out := buf.String()
	alphaIdx := strings.Index(out, "ALPHA")
	mangoIdx := strings.Index(out, "MANGO")
	zebraIdx := strings.Index(out, "ZEBRA")
	if alphaIdx > mangoIdx || mangoIdx > zebraIdx {
		t.Errorf("expected sorted output, got: %q", out)
	}
}
