package report

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/env"
)

func TestRenderGroup_NoEntries(t *testing.T) {
	var buf strings.Builder
	RenderGroupTo(&buf, []env.GroupResult{}, false)
	if !strings.Contains(buf.String(), "no entries") {
		t.Errorf("expected 'no entries', got %q", buf.String())
	}
}

func TestRenderGroup_SingleGroup(t *testing.T) {
	groups := []env.GroupResult{
		{Prefix: "DB", Entries: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}},
	}
	var buf strings.Builder
	RenderGroupTo(&buf, groups, false)
	out := buf.String()
	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] header")
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Error("expected DB_HOST=localhost")
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Error("expected DB_PORT=5432")
	}
}

func TestRenderGroup_Ungrouped(t *testing.T) {
	groups := []env.GroupResult{
		{Prefix: "", Entries: map[string]string{"HOST": "localhost"}},
	}
	var buf strings.Builder
	RenderGroupTo(&buf, groups, false)
	if !strings.Contains(buf.String(), "(ungrouped)") {
		t.Error("expected (ungrouped) label")
	}
}

func TestRenderGroup_ColorOutput(t *testing.T) {
	groups := []env.GroupResult{
		{Prefix: "APP", Entries: map[string]string{"APP_NAME": "myapp"}},
	}
	var buf strings.Builder
	RenderGroupTo(&buf, groups, true)
	out := buf.String()
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in color output")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in output")
	}
}

func TestRenderGroup_SortedKeys(t *testing.T) {
	groups := []env.GroupResult{
		{Prefix: "X", Entries: map[string]string{"X_Z": "3", "X_A": "1", "X_M": "2"}},
	}
	var buf strings.Builder
	RenderGroupTo(&buf, groups, false)
	out := buf.String()
	idxA := strings.Index(out, "X_A")
	idxM := strings.Index(out, "X_M")
	idxZ := strings.Index(out, "X_Z")
	if !(idxA < idxM && idxM < idxZ) {
		t.Error("expected keys sorted alphabetically")
	}
}
