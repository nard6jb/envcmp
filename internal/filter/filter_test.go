package filter_test

import (
	"testing"

	"github.com/yourorg/envcmp/internal/filter"
)

var base = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
}

func TestApply_NoOptions(t *testing.T) {
	out := filter.Apply(base, filter.Options{})
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out := filter.Apply(base, filter.Options{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApply_Keys(t *testing.T) {
	out := filter.Apply(base, filter.Options{Keys: []string{"LOG_LEVEL", "DB_HOST"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_Exclude(t *testing.T) {
	out := filter.Apply(base, filter.Options{Exclude: []string{"DB_PASSWORD"}})
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if len(out) != len(base)-1 {
		t.Fatalf("expected %d keys, got %d", len(base)-1, len(out))
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	out := filter.Apply(base, filter.Options{Prefix: "DB_", Exclude: []string{"DB_PASSWORD"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	out := filter.Apply(map[string]string{}, filter.Options{Prefix: "APP_"})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}
