package env

import (
	"testing"
)

func TestNormalize_NoOptions(t *testing.T) {
	input := map[string]string{"KEY": "  value  ", "OTHER": ""}
	out := Normalize(input, NormalizeOptions{})
	if out["KEY"] != "  value  " {
		t.Errorf("expected untrimmed value, got %q", out["KEY"])
	}
	if _, ok := out["OTHER"]; !ok {
		t.Error("expected empty key to be preserved")
	}
}

func TestNormalize_TrimSpace(t *testing.T) {
	input := map[string]string{"KEY": "  hello  ", "B": "\t world\t"}
	out := Normalize(input, NormalizeOptions{TrimSpace: true})
	if out["KEY"] != "hello" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
	if out["B"] != "world" {
		t.Errorf("expected trimmed value, got %q", out["B"])
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	input := map[string]string{"KEY": "value", "EMPTY": "", "SPACES": "   "}
	out := Normalize(input, NormalizeOptions{TrimSpace: true, RemoveEmpty: true})
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := out["SPACES"]; ok {
		t.Error("expected SPACES to be removed after trim")
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY to be preserved, got %q", out["KEY"])
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	input := map[string]string{"key": "val", "mixed_Key": "v2"}
	out := Normalize(input, NormalizeOptions{UppercaseKeys: true})
	if out["KEY"] != "val" {
		t.Errorf("expected KEY, got map %v", out)
	}
	if out["MIXED_KEY"] != "v2" {
		t.Errorf("expected MIXED_KEY, got map %v", out)
	}
}

func TestNormalize_LowercaseKeys(t *testing.T) {
	input := map[string]string{"KEY": "val", "Mixed": "v2"}
	out := Normalize(input, NormalizeOptions{LowercaseKeys: true})
	if out["key"] != "val" {
		t.Errorf("expected key, got map %v", out)
	}
	if out["mixed"] != "v2" {
		t.Errorf("expected mixed, got map %v", out)
	}
}

func TestNormalize_LowercaseTakesPrecedence(t *testing.T) {
	input := map[string]string{"Key": "val"}
	out := Normalize(input, NormalizeOptions{UppercaseKeys: true, LowercaseKeys: true})
	if _, ok := out["key"]; !ok {
		t.Errorf("expected lowercase to take precedence, got map %v", out)
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	_ = Normalize(input, NormalizeOptions{TrimSpace: true})
	if input["KEY"] != "  value  " {
		t.Error("original map was mutated")
	}
}
