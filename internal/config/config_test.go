package config

import (
	"testing"
)

func TestParse_DiffMode(t *testing.T) {
	cfg, err := Parse([]string{"-mode", "diff", ".env.dev", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mode != ModeDiff {
		t.Errorf("expected diff mode, got %q", cfg.Mode)
	}
	if cfg.FileA != ".env.dev" || cfg.FileB != ".env.prod" {
		t.Errorf("unexpected files: %q %q", cfg.FileA, cfg.FileB)
	}
	if !cfg.MaskKeys {
		t.Error("expected mask to default to true")
	}
}

func TestParse_ValidateMode(t *testing.T) {
	cfg, err := Parse([]string{"-mode", "validate", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mode != ModeValidate {
		t.Errorf("expected validate mode, got %q", cfg.Mode)
	}
}

func TestParse_NoColorFlag(t *testing.T) {
	cfg, err := Parse([]string{"-no-color", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("expected no-color to be true")
	}
}

func TestParse_MissingFiles(t *testing.T) {
	_, err := Parse([]string{"-mode", "diff"})
	if err == nil {
		t.Error("expected error for missing positional args")
	}
}

func TestParse_UnknownMode(t *testing.T) {
	_, err := Parse([]string{"-mode", "export", "a.env", "b.env"})
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}

func TestParse_MaskFalse(t *testing.T) {
	cfg, err := Parse([]string{"-mask=false", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaskKeys {
		t.Error("expected mask to be false")
	}
}
