package redact_test

import (
	"sort"
	"testing"

	"github.com/yourusername/envcmp/internal/redact"
)

func TestApplyToMap_MasksSensitive(t *testing.T) {
	env := map[string]string{
		"APP_SECRET": "supersecret",
		"APP_HOST":   "localhost",
	}
	entries := redact.ApplyToMap(env)
	for _, e := range entries {
		if e.Key == "APP_SECRET" {
			if !e.Redacted {
				t.Error("expected APP_SECRET to be redacted")
			}
			if e.Value == "supersecret" {
				t.Error("expected APP_SECRET value to be masked")
			}
		}
		if e.Key == "APP_HOST" {
			if e.Redacted {
				t.Error("expected APP_HOST not to be redacted")
			}
			if e.Value != "localhost" {
				t.Errorf("expected APP_HOST=localhost, got %s", e.Value)
			}
		}
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"KEY": "val", "OTHER": "data"}
	entries := redact.ApplyToMap(env)
	out := redact.ToMap(entries)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestRedactedKeys_ReturnsSensitive(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "pass",
		"DB_HOST":     "host",
		"API_TOKEN":   "tok",
	}
	entries := redact.ApplyToMap(env)
	keys := redact.RedactedKeys(entries)
	sort.Strings(keys)
	if len(keys) < 1 {
		t.Error("expected at least one redacted key")
	}
	for _, k := range keys {
		if k == "DB_HOST" {
			t.Error("DB_HOST should not be redacted")
		}
	}
}

func TestApplyToMap_EmptyMap(t *testing.T) {
	entries := redact.ApplyToMap(map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
