package env

import (
	"testing"
)

func TestTrim_NoOptions(t *testing.T) {
	in := map[string]string{"KEY": "value", "OTHER": "data"}
	out, results := Trim(in, TrimOptions{})

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, r := range results {
		if r.KeyChanged || r.ValueChanged {
			t.Errorf("expected no changes, got key=%q value=%q", r.Key, r.Value)
		}
	}
}

func TestTrim_TrimValues(t *testing.T) {
	in := map[string]string{"KEY": "  hello  ", "OTHER": "\tworld\t"}
	out, results := Trim(in, TrimOptions{TrimValues: true})

	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected 'world', got %q", out["OTHER"])
	}
	for _, r := range results {
		if r.ValueChanged && r.OriginalValue == r.Value {
			t.Errorf("ValueChanged set but values identical for key %q", r.Key)
		}
	}
}

func TestTrim_TrimKeys(t *testing.T) {
	in := map[string]string{" KEY ": "value"}
	out, _ := Trim(in, TrimOptions{TrimKeys: true})

	if _, ok := out["KEY"]; !ok {
		t.Error("expected trimmed key 'KEY' to exist")
	}
	if _, ok := out[" KEY "]; ok {
		t.Error("expected original key ' KEY ' to be removed")
	}
}

func TestTrim_TrimPrefix(t *testing.T) {
	in := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "val"}
	out, results := Trim(in, TrimOptions{TrimPrefix: "APP_"})

	if _, ok := out["HOST"]; !ok {
		t.Error("expected key 'HOST'")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected key 'PORT'")
	}
	if _, ok := out["OTHER"]; !ok {
		t.Error("expected key 'OTHER' unchanged")
	}

	changed := 0
	for _, r := range results {
		if r.KeyChanged {
			changed++
		}
	}
	if changed != 2 {
		t.Errorf("expected 2 key changes, got %d", changed)
	}
}

func TestTrim_TrimSuffix(t *testing.T) {
	in := map[string]string{"HOST_ENV": "prod", "PORT_ENV": "443"}
	out, _ := Trim(in, TrimOptions{TrimSuffix: "_ENV"})

	if _, ok := out["HOST"]; !ok {
		t.Error("expected key 'HOST'")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected key 'PORT'")
	}
}

func TestTrim_CombinedOptions(t *testing.T) {
	in := map[string]string{"APP_ KEY ": "  value  "}
	out, results := Trim(in, TrimOptions{
		TrimKeys:   true,
		TrimValues: true,
		TrimPrefix: "APP_",
	})

	if out["KEY"] != "value" {
		t.Errorf("expected out[KEY]=value, got %q", out["KEY"])
	}
	if len(results) != 1 || !results[0].KeyChanged || !results[0].ValueChanged {
		t.Error("expected both key and value to be marked changed")
	}
}
