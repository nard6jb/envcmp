package validate_test

import (
	"testing"

	"github.com/your/envcmp/internal/validate"
)

func TestAgainst_Valid(t *testing.T) {
	ref := map[string]string{"HOST": "localhost", "PORT": "8080"}
	tgt := map[string]string{"HOST": "prod.example.com", "PORT": "443"}

	r := validate.Against(ref, tgt)
	if !r.Valid {
		t.Errorf("expected valid, got invalid with missing: %v", r.MissingKeys)
	}
}

func TestAgainst_MissingKeys(t *testing.T) {
	ref := map[string]string{"HOST": "localhost", "PORT": "8080", "SECRET": "x"}
	tgt := map[string]string{"HOST": "prod.example.com"}

	r := validate.Against(ref, tgt)
	if r.Valid {
		t.Error("expected invalid result")
	}
	if len(r.MissingKeys) != 2 {
		t.Errorf("expected 2 missing keys, got %d: %v", len(r.MissingKeys), r.MissingKeys)
	}
}

func TestAgainst_ExtraKeys(t *testing.T) {
	ref := map[string]string{"HOST": "localhost"}
	tgt := map[string]string{"HOST": "prod", "EXTRA_VAR": "surprise"}

	r := validate.Against(ref, tgt)
	if !r.Valid {
		t.Error("expected valid (extra keys do not invalidate)")
	}
	if len(r.ExtraKeys) != 1 {
		t.Errorf("expected 1 extra key, got %d", len(r.ExtraKeys))
	}
}

func TestAgainst_EmptyReference(t *testing.T) {
	ref := map[string]string{}
	tgt := map[string]string{"FOO": "bar"}

	r := validate.Against(ref, tgt)
	if !r.Valid {
		t.Error("expected valid for empty reference")
	}
}
