package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcmp/internal/profile"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "profiles.json")
}

func TestStore_AddAndGet(t *testing.T) {
	s := &profile.Store{}
	s.Add(profile.Profile{Name: "staging", Files: []string{".env.staging"}})
	p, err := s.Get("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Files) != 1 || p.Files[0] != ".env.staging" {
		t.Errorf("unexpected files: %v", p.Files)
	}
}

func TestStore_AddReplaces(t *testing.T) {
	s := &profile.Store{}
	s.Add(profile.Profile{Name: "prod", Files: []string{".env.prod"}})
	s.Add(profile.Profile{Name: "prod", Files: []string{".env.prod", ".env.prod.local"}})
	if len(s.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(s.Profiles))
	}
	if len(s.Profiles[0].Files) != 2 {
		t.Errorf("expected 2 files after replace, got %d", len(s.Profiles[0].Files))
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := &profile.Store{}
	_, err := s.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestStore_Remove(t *testing.T) {
	s := &profile.Store{}
	s.Add(profile.Profile{Name: "dev", Files: []string{".env"}})
	removed := s.Remove("dev")
	if !removed {
		t.Error("expected Remove to return true")
	}
	if len(s.Profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(s.Profiles))
	}
}

func TestStore_RemoveMissing(t *testing.T) {
	s := &profile.Store{}
	if s.Remove("ghost") {
		t.Error("expected Remove to return false for missing profile")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempStorePath(t)
	s := &profile.Store{}
	s.Add(profile.Profile{Name: "ci", Files: []string{".env.ci", ".env.shared"}})

	if err := profile.SaveStore(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := profile.LoadStore(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	p, err := loaded.Get("ci")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(p.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(p.Files))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	s, err := profile.LoadStore("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("expected empty store, got error: %v", err)
	}
	if len(s.Profiles) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestSave_RestrictedPerms(t *testing.T) {
	path := tempStorePath(t)
	s := &profile.Store{}
	if err := profile.SaveStore(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
