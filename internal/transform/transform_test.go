package transform_test

import (
	"testing"

	"github.com/user/envcmp/internal/transform"
)

func base() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"db_pass":  "secret",
	}
}

func TestApply_NoOptions(t *testing.T) {
	result := transform.Apply(base(), transform.Options{})
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", result["APP_HOST"])
	}
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestApply_UppercaseKeys(t *testing.T) {
	result := transform.Apply(base(), transform.Options{UppercaseKeys: true})
	if _, ok := result["DB_PASS"]; !ok {
		t.Error("expected DB_PASS after uppercase, not found")
	}
	if _, ok := result["db_pass"]; ok {
		t.Error("expected original db_pass to be gone")
	}
}

func TestApply_LowercaseKeys(t *testing.T) {
	result := transform.Apply(base(), transform.Options{LowercaseKeys: true})
	if _, ok := result["app_host"]; !ok {
		t.Error("expected app_host after lowercase, not found")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	result := transform.Apply(base(), transform.Options{StripPrefix: "APP_"})
	if _, ok := result["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := result["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_ prefix")
	}
	// db_pass has no APP_ prefix, should be unchanged
	if _, ok := result["db_pass"]; !ok {
		t.Error("expected db_pass unchanged (no prefix match)")
	}
}

func TestApply_AddPrefix(t *testing.T) {
	result := transform.Apply(base(), transform.Options{AddPrefix: "PROD_"})
	if _, ok := result["PROD_APP_HOST"]; !ok {
		t.Error("expected PROD_APP_HOST after adding prefix")
	}
}

func TestApply_RenameMap(t *testing.T) {
	renames := map[string]string{"APP_HOST": "SERVICE_HOST"}
	result := transform.Apply(base(), transform.Options{RenameMap: renames})
	if result["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q", result["SERVICE_HOST"])
	}
	if _, ok := result["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be renamed away")
	}
}

func TestApply_RenameOverridesOtherOpts(t *testing.T) {
	renames := map[string]string{"APP_HOST": "host"}
	opts := transform.Options{
		UppercaseKeys: true,
		RenameMap:     renames,
	}
	result := transform.Apply(map[string]string{"APP_HOST": "localhost"}, opts)
	if _, ok := result["host"]; !ok {
		t.Error("expected rename to take priority over uppercase")
	}
	if _, ok := result["HOST"]; ok {
		t.Error("expected uppercase not applied when rename matches")
	}
}
