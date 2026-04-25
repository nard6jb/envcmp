package env

import (
	"testing"
)

func TestFlatten_NoOptions(t *testing.T) {
	groups := map[string]map[string]string{
		"db": {"HOST": "localhost", "PORT": "5432"},
	}
	result, entries := Flatten(groups, FlattenOptions{})

	if result["db_HOST"] != "localhost" {
		t.Errorf("expected db_HOST=localhost, got %q", result["db_HOST"])
	}
	if result["db_PORT"] != "5432" {
		t.Errorf("expected db_PORT=5432, got %q", result["db_PORT"])
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	groups := map[string]map[string]string{
		"app": {"NAME": "envcmp"},
	}
	result, _ := Flatten(groups, FlattenOptions{Separator: "."})
	if result["app.NAME"] != "envcmp" {
		t.Errorf("expected app.NAME=envcmp, got %q", result["app.NAME"])
	}
}

func TestFlatten_UppercaseKeys(t *testing.T) {
	groups := map[string]map[string]string{
		"cache": {"ttl": "300"},
	}
	result, _ := Flatten(groups, FlattenOptions{UppercaseKeys: true})
	if result["CACHE_TTL"] != "300" {
		t.Errorf("expected CACHE_TTL=300, got %q", result["CACHE_TTL"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	groups := map[string]map[string]string{
		"db": {"USER": "admin"},
	}
	result, _ := Flatten(groups, FlattenOptions{Prefix: "APP"})
	if result["APP_db_USER"] != "admin" {
		t.Errorf("expected APP_db_USER=admin, got %q", result["APP_db_USER"])
	}
}

func TestFlatten_MultipleGroups_DeterministicOrder(t *testing.T) {
	groups := map[string]map[string]string{
		"z": {"KEY": "z-val"},
		"a": {"KEY": "a-val"},
	}
	_, entries := Flatten(groups, FlattenOptions{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Groups sorted: "a" comes before "z"
	if entries[0].Key != "a_KEY" {
		t.Errorf("expected first entry a_KEY, got %q", entries[0].Key)
	}
	if entries[1].Key != "z_KEY" {
		t.Errorf("expected second entry z_KEY, got %q", entries[1].Key)
	}
}

func TestFlatten_EmptyGroups(t *testing.T) {
	result, entries := Flatten(map[string]map[string]string{}, FlattenOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestFlatten_OriginalKeyPreserved(t *testing.T) {
	groups := map[string]map[string]string{
		"svc": {"timeout": "30s"},
	}
	_, entries := Flatten(groups, FlattenOptions{UppercaseKeys: true})
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Original != "timeout" {
		t.Errorf("expected original=timeout, got %q", entries[0].Original)
	}
	if entries[0].Key != "SVC_TIMEOUT" {
		t.Errorf("expected key=SVC_TIMEOUT, got %q", entries[0].Key)
	}
}
