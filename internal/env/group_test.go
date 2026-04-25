package env

import (
	"testing"
)

func TestGroup_NoOptions(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "myapp",
	}
	results := Group(env, GroupOptions{})
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	if results[0].Prefix != "APP" || results[1].Prefix != "DB" {
		t.Errorf("unexpected prefixes: %v, %v", results[0].Prefix, results[1].Prefix)
	}
}

func TestGroup_MinGroupSize(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SOLO_KEY": "val",
	}
	results := Group(env, GroupOptions{MinGroupSize: 2})
	// SOLO group has only 1 key → merged into ungrouped ("")
	var ungrouped *GroupResult
	for i := range results {
		if results[i].Prefix == "" {
			ungrouped = &results[i]
		}
	}
	if ungrouped == nil {
		t.Fatal("expected ungrouped bucket")
	}
	if _, ok := ungrouped.Entries["SOLO_KEY"]; !ok {
		t.Error("expected SOLO_KEY in ungrouped bucket")
	}
}

func TestGroup_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
		"app.name": "myapp",
	}
	results := Group(env, GroupOptions{PrefixSep: "."})
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
}

func TestGroup_NoPrefix(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	results := Group(env, GroupOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 group, got %d", len(results))
	}
	if results[0].Prefix != "" {
		t.Errorf("expected ungrouped prefix, got %q", results[0].Prefix)
	}
}

func TestGroup_EmptyMap(t *testing.T) {
	results := Group(map[string]string{}, GroupOptions{})
	if len(results) != 0 {
		t.Errorf("expected empty result, got %d groups", len(results))
	}
}

func TestGroup_SortedPrefixes(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "1",
		"A_KEY": "2",
		"M_KEY": "3",
	}
	results := Group(env, GroupOptions{})
	if results[0].Prefix != "A" || results[1].Prefix != "M" || results[2].Prefix != "Z" {
		t.Errorf("unexpected order: %v %v %v", results[0].Prefix, results[1].Prefix, results[2].Prefix)
	}
}
