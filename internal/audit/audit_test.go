package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envcmp/internal/audit"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestAppend_CreatesFile(t *testing.T) {
	p := tempLog(t)
	e := audit.Entry{Command: "diff", Files: []string{"a.env", "b.env"}}
	if err := audit.Append(p, e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("log file not created")
	}
}

func TestAppend_AccumulatesEntries(t *testing.T) {
	p := tempLog(t)
	for i := 0; i < 3; i++ {
		if err := audit.Append(p, audit.Entry{Command: "diff"}); err != nil {
			t.Fatalf("append %d failed: %v", i, err)
		}
	}
	l, err := audit.LoadLog(p)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(l.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(l.Entries))
	}
}

func TestAppend_SetsTimestamp(t *testing.T) {
	p := tempLog(t)
	before := time.Now().UTC()
	_ = audit.Append(p, audit.Entry{Command: "lint"})
	l, _ := audit.LoadLog(p)
	if l.Entries[0].Timestamp.Before(before) {
		t.Error("timestamp not set correctly")
	}
}

func TestLoadLog_FileNotFound(t *testing.T) {
	_, err := audit.LoadLog("/nonexistent/path/audit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAppend_StoresFields(t *testing.T) {
	p := tempLog(t)
	e := audit.Entry{
		Command: "validate",
		Files:   []string{"prod.env"},
		Issues:  []string{"missing KEY_A"},
	}
	_ = audit.Append(p, e)
	l, _ := audit.LoadLog(p)
	if l.Entries[0].Command != "validate" {
		t.Errorf("expected validate, got %s", l.Entries[0].Command)
	}
	if len(l.Entries[0].Issues) != 1 {
		t.Error("expected 1 issue")
	}
}
