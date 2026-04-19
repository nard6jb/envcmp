package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/evanw/envcmp/internal/merge"
)

func TestRenderMerge_NoConflicts(t *testing.T) {
	result := merge.Result{
		Merged:    map[string]string{"FOO": "bar"},
		Conflicts: map[string][2]string{},
	}
	var buf bytes.Buffer
	RenderMergeTo(&buf, result, false)
	if !strings.Contains(buf.String(), "no conflicts") {
		t.Errorf("expected no-conflict message, got: %s", buf.String())
	}
}

func TestRenderMerge_WithConflicts(t *testing.T) {
	result := merge.Result{
		Merged: map[string]string{"FOO": "first", "BAR": "baz"},
		Conflicts: map[string][2]string{
			"FOO": {"first", "second"},
		},
	}
	var buf bytes.Buffer
	RenderMergeTo(&buf, result, false)
	out := buf.String()
	if !strings.Contains(out, "conflicts") {
		t.Errorf("expected conflicts header, got: %s", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "first") {
		t.Errorf("expected kept value in output, got: %s", out)
	}
	if !strings.Contains(out, "second") {
		t.Errorf("expected discarded value in output, got: %s", out)
	}
}

func TestRenderMerge_ColorNoConflicts(t *testing.T) {
	result := merge.Result{
		Merged:    map[string]string{},
		Conflicts: map[string][2]string{},
	}
	var buf bytes.Buffer
	RenderMergeTo(&buf, result, true)
	if !strings.Contains(buf.String(), "\033[32m") {
		t.Errorf("expected green color code in output")
	}
}

func TestRenderMerge_ColorWithConflicts(t *testing.T) {
	result := merge.Result{
		Merged: map[string]string{"KEY": "a"},
		Conflicts: map[string][2]string{
			"KEY": {"a", "b"},
		},
	}
	var buf bytes.Buffer
	RenderMergeTo(&buf, result, true)
	if !strings.Contains(buf.String(), "\033[33m") {
		t.Errorf("expected yellow color code in output")
	}
}
