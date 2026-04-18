package mask_test

import (
	"testing"

	"github.com/yourorg/envcmp/internal/mask"
)

func TestIsSensitive_DetectsSecretKeys(t *testing.T) {
	sensitive := []string{
		"DB_PASSWORD", "API_KEY", "AUTH_TOKEN",
		"PRIVATE_KEY", "AWS_SECRET_ACCESS_KEY", "USER_CREDENTIAL",
	}
	for _, key := range sensitive {
		if !mask.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_AllowsNonSecretKeys(t *testing.T) {
	plain := []string{"PORT", "HOST", "APP_ENV", "LOG_LEVEL", "TIMEOUT"}
	for _, key := range plain {
		if mask.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMaskValue_MasksSensitive(t *testing.T) {
	got := mask.MaskValue("DB_PASSWORD", "supersecret")
	if got == "supersecret" {
		t.Error("expected value to be masked")
	}
	if got == "" {
		t.Error("expected non-empty mask placeholder")
	}
}

func TestMaskValue_PreservesPlain(t *testing.T) {
	got := mask.MaskValue("PORT", "8080")
	if got != "8080" {
		t.Errorf("expected '8080', got %q", got)
	}
}

func TestMaskMap(t *testing.T) {
	input := map[string]string{
		"PORT":        "8080",
		"DB_PASSWORD": "secret123",
		"API_KEY":     "key-abc",
		"APP_ENV":     "production",
	}
	out := mask.MaskMap(input)

	if out["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", out["PORT"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", out["APP_ENV"])
	}
	if out["DB_PASSWORD"] == "secret123" {
		t.Error("DB_PASSWORD should be masked")
	}
	if out["API_KEY"] == "key-abc" {
		t.Error("API_KEY should be masked")
	}
	// original map must not be mutated
	if input["DB_PASSWORD"] != "secret123" {
		t.Error("original map should not be mutated")
	}
}
