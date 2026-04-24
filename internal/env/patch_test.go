package env

import (
	"testing"
)

var patchBase = map[string]string{
	"HOST": "localhost",
	"PORT": "5432",
	"DB":   "mydb",
}

func TestPatch_NoOps(t *testing.T) {
	out, results, err := Patch(patchBase, nil, PatchOptions{AllowNew: true, AllowDelete: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", out["HOST"])
	}
}

func TestPatch_ChangeValue(t *testing.T) {
	ops := []PatchOp{{Key: "PORT", Value: "3306"}}
	out, results, err := Patch(patchBase, ops, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "3306" {
		t.Errorf("expected PORT=3306, got %s", out["PORT"])
	}
	if !results[0].Changed {
		t.Error("expected result to be marked Changed")
	}
	if results[0].OldValue != "5432" {
		t.Errorf("expected OldValue=5432, got %s", results[0].OldValue)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	ops := []PatchOp{{Key: "DB", Delete: true}}
	out, results, err := Patch(patchBase, ops, PatchOptions{AllowDelete: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB"]; ok {
		t.Error("expected DB to be deleted")
	}
	if !results[0].Deleted {
		t.Error("expected result to be marked Deleted")
	}
}

func TestPatch_DeleteNotAllowed(t *testing.T) {
	ops := []PatchOp{{Key: "DB", Delete: true}}
	_, _, err := Patch(patchBase, ops, PatchOptions{AllowDelete: false})
	if err == nil {
		t.Error("expected error when delete not allowed")
	}
}

func TestPatch_AddNewKey(t *testing.T) {
	ops := []PatchOp{{Key: "TIMEOUT", Value: "30s"}}
	out, results, err := Patch(patchBase, ops, PatchOptions{AllowNew: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %s", out["TIMEOUT"])
	}
	if !results[0].Added {
		t.Error("expected result to be marked Added")
	}
}

func TestPatch_NewKeyNotAllowed_Skipped(t *testing.T) {
	ops := []PatchOp{{Key: "NEWKEY", Value: "val"}}
	_, results, err := Patch(patchBase, ops, PatchOptions{AllowNew: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Error("expected result to be Skipped")
	}
}

func TestPatch_ErrorOnMissing(t *testing.T) {
	ops := []PatchOp{{Key: "GHOST", Value: "val"}}
	_, _, err := Patch(patchBase, ops, PatchOptions{AllowNew: false, ErrorOnMissing: true})
	if err == nil {
		t.Error("expected error for missing key with ErrorOnMissing")
	}
}

func TestPatch_DoesNotMutateBase(t *testing.T) {
	ops := []PatchOp{{Key: "HOST", Value: "remotehost"}}
	_, _, _ = Patch(patchBase, ops, PatchOptions{})
	if patchBase["HOST"] != "localhost" {
		t.Error("Patch must not mutate the base map")
	}
}
