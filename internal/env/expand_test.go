package env

import (
	"os"
	"testing"
)

func TestExpand_NoReferences(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Expand(src, nil, ExpandOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expanded["FOO"] != "bar" || res.Expanded["BAZ"] != "qux" {
		t.Errorf("expected values unchanged, got %v", res.Expanded)
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestExpand_ResolvesFromBase(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	src := map[string]string{"DSN": "postgres://${HOST}:${PORT}/db"}
	res, err := Expand(src, base, ExpandOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://localhost:5432/db"
	if res.Expanded["DSN"] != want {
		t.Errorf("got %q, want %q", res.Expanded["DSN"], want)
	}
}

func TestExpand_UnresolvedCollected(t *testing.T) {
	src := map[string]string{"URL": "http://${MISSING_HOST}/path"}
	res, err := Expand(src, nil, ExpandOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "MISSING_HOST" {
		t.Errorf("expected [MISSING_HOST], got %v", res.Unresolved)
	}
}

func TestExpand_StrictModeErrors(t *testing.T) {
	src := map[string]string{"URL": "http://${GHOST}/path"}
	_, err := Expand(src, nil, ExpandOptions{StrictMode: true})
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	os.Setenv("ENVCMP_TEST_OS_VAR", "from-os")
	defer os.Unsetenv("ENVCMP_TEST_OS_VAR")

	src := map[string]string{"VAL": "${ENVCMP_TEST_OS_VAR}"}
	res, err := Expand(src, nil, ExpandOptions{FallbackToOS: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expanded["VAL"] != "from-os" {
		t.Errorf("got %q, want %q", res.Expanded["VAL"], "from-os")
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestExpand_DeduplicatesUnresolved(t *testing.T) {
	src := map[string]string{
		"A": "${GHOST}",
		"B": "${GHOST}",
	}
	res, err := Expand(src, nil, ExpandOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) != 1 {
		t.Errorf("expected 1 unique unresolved, got %v", res.Unresolved)
	}
}
