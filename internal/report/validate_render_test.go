package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your/envcmp/internal/report"
	"github.com/your/envcmp/internal/validate"
)

func TestRenderValidation_Valid(t *testing.T) {
	var buf bytes.Buffer
	r := validate.Result{Valid: true}
	report.RenderValidation(&buf, r, ".env")
	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected valid message, got: %s", buf.String())
	}
}

func TestRenderValidation_MissingKeys(t *testing.T) {
	var buf bytes.Buffer
	r := validate.Result{
		Valid:       false,
		MissingKeys: []string{"SECRET_KEY", "DB_URL"},
	}
	report.RenderValidation(&buf, r, ".env.prod")
	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Error("expected SECRET_KEY in output")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
	if !strings.Contains(out, "invalid") {
		t.Error("expected invalid summary")
	}
}

func TestRenderValidation_ExtraKeys(t *testing.T) {
	var buf bytes.Buffer
	r := validate.Result{
		Valid:     true,
		ExtraKeys: []string{"UNUSED_VAR"},
	}
	report.RenderValidation(&buf, r, ".env.staging")
	out := buf.String()
	if !strings.Contains(out, "UNUSED_VAR") {
		t.Error("expected UNUSED_VAR in output")
	}
	if !strings.Contains(out, "extra") {
		t.Error("expected extra keys warning")
	}
}
