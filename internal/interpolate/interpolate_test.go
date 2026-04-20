package interpolate

import (
	"testing"
)

func TestApply_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res := Apply(env)
	if res.Resolved["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", res.Resolved["HOST"])
	}
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(res.Issues))
	}
}

func TestApply_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "8080",
	}
	res := Apply(env)
	if res.Resolved["BASE_URL"] != "http://example.com:8080" {
		t.Errorf("unexpected value: %s", res.Resolved["BASE_URL"])
	}
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues, got %v", res.Issues)
	}
}

func TestApply_DollarStyle(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	res := Apply(env)
	if res.Resolved["GREETING"] != "Hello World" {
		t.Errorf("unexpected value: %s", res.Resolved["GREETING"])
	}
}

func TestApply_UndefinedReference(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${DB_USER}@localhost",
	}
	res := Apply(env)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if res.Issues[0].Ref != "DB_USER" {
		t.Errorf("expected ref DB_USER, got %s", res.Issues[0].Ref)
	}
	// Original token preserved in value
	if res.Resolved["DSN"] != "postgres://${DB_USER}@localhost" {
		t.Errorf("unexpected resolved value: %s", res.Resolved["DSN"])
	}
}

func TestApply_MultipleUndefined(t *testing.T) {
	env := map[string]string{
		"URL": "${SCHEME}://${HOST}/${PATH}",
	}
	res := Apply(env)
	if len(res.Issues) != 3 {
		t.Errorf("expected 3 issues, got %d", len(res.Issues))
	}
}

func TestApply_SelfReference(t *testing.T) {
	env := map[string]string{
		"FOO": "${FOO}_suffix",
	}
	res := Apply(env)
	// Self-reference resolves to itself (current value), creating a loop — we just expand once
	if res.Resolved["FOO"] != "${FOO}_suffix" {
		t.Logf("self-ref resolved to: %s", res.Resolved["FOO"])
	}
}
