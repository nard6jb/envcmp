package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envcmp/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
		"API_KEY": "secret123",
	}

	if err := snapshot.Save(path, ".env.prod", entries); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if snap.Source != ".env.prod" {
		t.Errorf("expected source .env.prod, got %s", snap.Source)
	}

	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}

	for k, v := range entries {
		if snap.Entries[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, snap.Entries[k])
		}
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte("not json{"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSave_CreatesFileWithRestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(path, "test", map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected perm 0600, got %o", perm)
	}
}

func TestSave_EmptyEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	if err := snapshot.Save(path, ".env.empty", map[string]string{}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(snap.Entries) != 0 {
		t.Errorf("expected empty entries, got %d entries", len(snap.Entries))
	}

	if snap.Source != ".env.empty" {
		t.Errorf("expected source .env.empty, got %s", snap.Source)
	}
}
