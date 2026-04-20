package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/interpolate"
)

func TestRenderInterpolation_NoIssues(t *testing.T) {
	res := interpolate.Result{
		Resolved: map[string]string{"FOO": "bar"},
		Issues:   nil,
	}
	var buf bytes.Buffer
	RenderInterpolation(&buf, res, false)
	if !strings.Contains(buf.String(), "No interpolation issues") {
		t.Errorf("expected no-issues message, got: %s", buf.String())
	}
}

func TestRenderInterpolation_WithIssues(t *testing.T) {
	res := interpolate.Result{
		Resolved: map[string]string{"DSN": "postgres://${DB_USER}@localhost"},
		Issues: []interpolate.Issue{
			{Key: "DSN", Ref: "DB_USER", Message: "undefined variable reference: DB_USER"},
		},
	}
	var buf bytes.Buffer
	RenderInterpolation(&buf, res, false)
	out := buf.String()
	if !strings.Contains(out, "1 interpolation issue") {
		t.Errorf("expected issue count, got: %s", out)
	}
	if !strings.Contains(out, "DB_USER") {
		t.Errorf("expected DB_USER in output, got: %s", out)
	}
}

func TestRenderInterpolation_SortedOutput(t *testing.T) {
	res := interpolate.Result{
		Issues: []interpolate.Issue{
			{Key: "Z_KEY", Ref: "MISSING", Message: "undefined variable reference: MISSING"},
			{Key: "A_KEY", Ref: "GONE", Message: "undefined variable reference: GONE"},
		},
	}
	var buf bytes.Buffer
	RenderInterpolation(&buf, res, false)
	out := buf.String()
	if strings.Index(out, "A_KEY") > strings.Index(out, "Z_KEY") {
		t.Error("expected A_KEY before Z_KEY in sorted output")
	}
}

func TestRenderInterpolation_ColorNoIssues(t *testing.T) {
	res := interpolate.Result{Resolved: map[string]string{"X": "1"}}
	var buf bytes.Buffer
	RenderInterpolation(&buf, res, true)
	if !strings.Contains(buf.String(), "\033[32m") {
		t.Error("expected green ANSI code for no-issues color output")
	}
}

func TestRenderInterpolation_ColorWithIssues(t *testing.T) {
	res := interpolate.Result{
		Issues: []interpolate.Issue{
			{Key: "URL", Ref: "HOST", Message: "undefined variable reference: HOST"},
		},
	}
	var buf bytes.Buffer
	RenderInterpolation(&buf, res, true)
	if !strings.Contains(buf.String(), "\033[31m") {
		t.Error("expected red ANSI code for issues color output")
	}
}
