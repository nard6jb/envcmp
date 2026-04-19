package lint_test

import (
	"testing"

	"github.com/yourusername/envcmp/internal/lint"
)

func entries(pairs ...string) [][2]string {
	var out [][2]string
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, [2]string{pairs[i], pairs[i+1]})
	}
	return out
}

func TestCheck_NoIssues(t *testing.T) {
	issue := lint.Check(entries("HOST", "localhost", "PORT", "8080"))
	if len(issue) != 0 {
		t.Fatalf("expected no issues, got %v", issue)
	}
}

func TestCheck_DuplicateKey(t *testing.T) {
	issues := lint.Check(entries("KEY", "a", "KEY", "b"))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "KEY" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	issues := lint.Check(entries("EMPTY", ""))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Message != "empty value" {
		t.Errorf("unexpected message %q", issues[0].Message)
	}
}

func TestCheck_UnquotedWhitespace(t *testing.T) {
	issues := lint.Check(entries("VAR", "hello world"))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestCheck_MultipleIssues(t *testing.T) {
	issues := lint.Check(entries("A", "", "A", "val", "B", "x y"))
	// empty A + duplicate A + whitespace B = 3
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d: %v", len(issues), issues)
	}
}

func TestIssue_String(t *testing.T) {
	i := lint.Issue{Line: 3, Key: "FOO", Message: "empty value"}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
