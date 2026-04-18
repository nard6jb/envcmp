package export_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envcmp/internal/diff"
	"github.com/user/envcmp/internal/export"
	"github.com/user/envcmp/internal/validate"
)

func TestWriteDiffJSON_Basic(t *testing.T) {
	result := map[string]diff.Entry{
		"FOO": {Status: diff.StatusChanged, Left: "a", Right: "b"},
		"BAR": {Status: diff.StatusMissingRight, Left: "x"},
	}
	var buf bytes.Buffer
	if err := export.WriteDiffJSON(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out export.DiffExport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "BAR" {
		t.Errorf("expected sorted first key BAR, got %s", out.Entries[0].Key)
	}
}

func TestWriteValidationJSON_Valid(t *testing.T) {
	result := validate.Result{Valid: true}
	var buf bytes.Buffer
	if err := export.WriteValidationJSON(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out export.ValidationExport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !out.Valid {
		t.Error("expected valid=true")
	}
}

func TestWriteValidationJSON_Missing(t *testing.T) {
	result := validate.Result{Valid: false, MissingKeys: []string{"SECRET"}}
	var buf bytes.Buffer
	_ = export.WriteValidationJSON(&buf, result)
	var out export.ValidationExport
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Valid {
		t.Error("expected valid=false")
	}
	if len(out.MissingKeys) != 1 || out.MissingKeys[0] != "SECRET" {
		t.Errorf("unexpected missing keys: %v", out.MissingKeys)
	}
}

func TestWriteDiffText_Output(t *testing.T) {
	result := map[string]diff.Entry{
		"KEY": {Status: diff.StatusEqual},
	}
	var buf bytes.Buffer
	if err := export.WriteDiffText(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text output")
	}
}
