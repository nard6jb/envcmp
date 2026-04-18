package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
	if env.Entries["DB_HOST"].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", env.Entries["DB_HOST"].Value)
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nAPP_ENV=production\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Entries["SECRET"].Value != "my secret value" {
		t.Errorf("unexpected value: %q", env.Entries["SECRET"].Value)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
