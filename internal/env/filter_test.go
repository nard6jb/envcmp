package env

import (
	"testing"
)

var filterBase = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
}

func TestFilter_NoOptions(t *testing.T) {
	out := Filter(filterBase, FilterOptions{})
	if len(out) != len(filterBase) {
		t.Fatalf("expected %d entries, got %d", len(filterBase), len(out))
	}
}

func TestFilter_Prefix(t *testing.T) {
	out := Filter(filterBase, FilterOptions{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestFilter_Keys(t *testing.T) {
	out := Filter(filterBase, FilterOptions{Keys: []string{"DB_HOST", "LOG_LEVEL"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out["DB_HOST"] != "db.local" {
		t.Errorf("unexpected value for DB_HOST: %s", out["DB_HOST"])
	}
}

func TestFilter_Exclude(t *testing.T) {
	out := Filter(filterBase, FilterOptions{Exclude: []string{"DB_PASSWORD", "LOG_LEVEL"}})
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := out["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should be excluded")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestFilter_PrefixAndExclude(t *testing.T) {
	out := Filter(filterBase, FilterOptions{
		Prefix:  "DB_",
		Exclude: []string{"DB_PASSWORD"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilter_EmptySource(t *testing.T) {
	out := Filter(map[string]string{}, FilterOptions{Prefix: "APP_"})
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(out))
	}
}
