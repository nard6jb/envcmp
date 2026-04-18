package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/diff"
	"github.com/user/envcmp/internal/report"
)

var noColorOpts = report.Options{Format: report.FormatText, NoColor: true}

func TestRender_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	report.Render(&buf, []diff.Result{}, noColorOpts)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestRender_MissingRight(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Missing, LeftValue: "localhost"},
	}
	report.Render(&buf, results, noColorOpts)
	out := buf.String()
	if !strings.Contains(out, "MISSING RIGHT") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestRender_MissingLeft(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "API_KEY", Status: diff.Extra, RightValue: "abc"},
	}
	report.Render(&buf, results, noColorOpts)
	out := buf.String()
	if !strings.Contains(out, "MISSING LEFT") || !strings.Contains(out, "API_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestRender_Changed(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "PORT", Status: diff.Changed, LeftValue: "3000", RightValue: "8080"},
	}
	report.Render(&buf, results, noColorOpts)
	out := buf.String()
	if !strings.Contains(out, "CHANGED") || !strings.Contains(out, "PORT") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "3000") || !strings.Contains(out, "8080") {
		t.Errorf("expected both values in output: %s", out)
	}
}

func TestRender_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "Z_VAR", Status: diff.Missing},
		{Key: "A_VAR", Status: diff.Extra},
	}
	report.Render(&buf, results, noColorOpts)
	out := buf.String()
	if strings.Index(out, "A_VAR") > strings.Index(out, "Z_VAR") {
		t.Error("expected output sorted alphabetically by key")
	}
}
