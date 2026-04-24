package env

import (
	"testing"

	"github.com/subtlepseudonym/envcmp/internal/merge"
)

func TestEnvMerge_NoConflicts(t *testing.T) {
	left := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	right := map[string]string{"APP_DEBUG": "true"}

	res, err := Merge(left, right, MergeOptions{Strategy: merge.StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Merged["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %s", res.Merged["APP_HOST"])
	}
	if res.Merged["APP_DEBUG"] != "true" {
		t.Errorf("expected APP_DEBUG=true, got %s", res.Merged["APP_DEBUG"])
	}
}

func TestEnvMerge_ConflictStrategyFirst(t *testing.T) {
	left := map[string]string{"APP_HOST": "localhost"}
	right := map[string]string{"APP_HOST": "remotehost"}

	res, err := Merge(left, right, MergeOptions{Strategy: merge.StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Chosen != "localhost" {
		t.Errorf("expected chosen=localhost, got %s", res.Conflicts[0].Chosen)
	}
}

func TestEnvMerge_ConflictStrategyLast(t *testing.T) {
	left := map[string]string{"DB_PASS": "secret1"}
	right := map[string]string{"DB_PASS": "secret2"}

	res, err := Merge(left, right, MergeOptions{Strategy: merge.StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Conflicts[0].Chosen != "secret2" {
		t.Errorf("expected chosen=secret2, got %s", res.Conflicts[0].Chosen)
	}
}

func TestEnvMerge_MaskSecrets(t *testing.T) {
	left := map[string]string{"API_SECRET": "abc123"}
	right := map[string]string{"API_SECRET": "xyz789"}

	res, err := Merge(left, right, MergeOptions{Strategy: merge.StrategyFirst, MaskSecrets: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	c := res.Conflicts[0]
	if c.Left == "abc123" || c.Right == "xyz789" {
		t.Errorf("expected secret values to be masked, got left=%s right=%s", c.Left, c.Right)
	}
}

func TestEnvMerge_NilLeft(t *testing.T) {
	_, err := Merge(nil, map[string]string{}, MergeOptions{})
	if err == nil {
		t.Error("expected error for nil left map")
	}
}

func TestEnvMerge_NilRight(t *testing.T) {
	_, err := Merge(map[string]string{}, nil, MergeOptions{})
	if err == nil {
		t.Error("expected error for nil right map")
	}
}
