package diff_test

import (
	"testing"

	"github.com/yourorg/envcmp/internal/diff"
)

func TestCompare_NoDiff(t *testing.T) {
	left := map[string]string{"KEY": "value", "PORT": "8080"}
	right := map[string]string{"KEY": "value", "PORT": "8080"}
	res := diff.Compare(left, right)
	if res.HasDiff() {
		t.Errorf("expected no diff, got %+v", res)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"KEY": "value", "EXTRA": "only-left"}
	right := map[string]string{"KEY": "value"}
	res := diff.Compare(left, right)
	if len(res.MissingInRight) != 1 || res.MissingInRight[0] != "EXTRA" {
		t.Errorf("expected EXTRA missing in right, got %v", res.MissingInRight)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"KEY": "value"}
	right := map[string]string{"KEY": "value", "NEW": "only-right"}
	res := diff.Compare(left, right)
	if len(res.MissingInLeft) != 1 || res.MissingInLeft[0] != "NEW" {
		t.Errorf("expected NEW missing in left, got %v", res.MissingInLeft)
	}
}

func TestCompare_ChangedValues(t *testing.T) {
	left := map[string]string{"DB_URL": "localhost", "PORT": "8080"}
	right := map[string]string{"DB_URL": "prod-host", "PORT": "8080"}
	res := diff.Compare(left, right)
	pair, ok := res.Changed["DB_URL"]
	if !ok {
		t.Fatal("expected DB_URL to be in Changed")
	}
	if pair[0] != "localhost" || pair[1] != "prod-host" {
		t.Errorf("unexpected changed values: %v", pair)
	}
	if _, ok := res.Changed["PORT"]; ok {
		t.Error("PORT should not appear in Changed")
	}
}

func TestCompare_HasDiff(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "2"}
	res := diff.Compare(left, right)
	if !res.HasDiff() {
		t.Error("expected HasDiff to return true")
	}
}
