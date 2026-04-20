package schema_test

import (
	"testing"

	"github.com/user/envcmp/internal/schema"
)

func TestValidate_NoIssues(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "PORT", Type: schema.TypeInt, Required: true},
			{Key: "DEBUG", Type: schema.TypeBool, Required: false},
			{Key: "API_URL", Type: schema.TypeURL, Required: true},
		},
	}
	env := map[string]string{
		"PORT":    "8080",
		"DEBUG":   "true",
		"API_URL": "https://example.com",
	}
	issues := s.Validate(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "PORT", Type: schema.TypeInt, Required: true},
		},
	}
	issues := s.Validate(map[string]string{})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected key PORT, got %q", issues[0].Key)
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "DEBUG", Type: schema.TypeBool, Required: false},
		},
	}
	issues := s.Validate(map[string]string{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues for optional missing key, got %v", issues)
	}
}

func TestValidate_WrongType(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "PORT", Type: schema.TypeInt, Required: true},
			{Key: "ENABLED", Type: schema.TypeBool, Required: true},
			{Key: "API_URL", Type: schema.TypeURL, Required: true},
		},
	}
	env := map[string]string{
		"PORT":    "not-a-number",
		"ENABLED": "yes",
		"API_URL": "ftp://bad",
	}
	issues := s.Validate(env)
	if len(issues) != 3 {
		t.Fatalf("expected 3 type issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "ENV", Type: schema.TypeString, Required: true, Pattern: `^(production|staging|development)$`},
		},
	}
	env := map[string]string{"ENV": "unknown"}
	issues := s.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 pattern issue, got %d", len(issues))
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "ENV", Type: schema.TypeString, Required: true, Pattern: `^(production|staging|development)$`},
		},
	}
	env := map[string]string{"ENV": "staging"}
	issues := s.Validate(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}
