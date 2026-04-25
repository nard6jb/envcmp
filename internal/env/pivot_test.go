package env

import (
	"testing"

	"github.com/jasonlovesdoggo/envcmp/internal/envfile"
)

func mkMap(pairs ...string) map[string]envfile.Entry {
	m := make(map[string]envfile.Entry, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = envfile.Entry{Key: pairs[i], Value: pairs[i+1]}
	}
	return m
}

func TestPivot_NoMaps(t *testing.T) {
	labels, rows := Pivot(nil, PivotOptions{})
	if len(labels) != 0 {
		t.Fatalf("expected 0 labels, got %d", len(labels))
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(rows))
	}
}

func TestPivot_DefaultLabels(t *testing.T) {
	m1 := mkMap("HOST", "localhost")
	m2 := mkMap("HOST", "prod.example.com")
	labels, _ := Pivot([]map[string]envfile.Entry{m1, m2}, PivotOptions{})
	if labels[0] != "env0" || labels[1] != "env1" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestPivot_CustomLabels(t *testing.T) {
	m1 := mkMap("HOST", "localhost")
	m2 := mkMap("HOST", "prod.example.com")
	labels, _ := Pivot([]map[string]envfile.Entry{m1, m2}, PivotOptions{Labels: []string{"dev", "prod"}})
	if labels[0] != "dev" || labels[1] != "prod" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestPivot_SameValues(t *testing.T) {
	m1 := mkMap("PORT", "8080")
	m2 := mkMap("PORT", "8080")
	_, rows := Pivot([]map[string]envfile.Entry{m1, m2}, PivotOptions{})
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if !rows[0].Same {
		t.Error("expected Same=true for identical values")
	}
}

func TestPivot_DifferentValues(t *testing.T) {
	m1 := mkMap("PORT", "8080")
	m2 := mkMap("PORT", "9090")
	_, rows := Pivot([]map[string]envfile.Entry{m1, m2}, PivotOptions{})
	if rows[0].Same {
		t.Error("expected Same=false for different values")
	}
	if rows[0].Values[0] != "8080" || rows[0].Values[1] != "9090" {
		t.Errorf("unexpected values: %v", rows[0].Values)
	}
}

func TestPivot_MissingKey(t *testing.T) {
	m1 := mkMap("HOST", "localhost", "PORT", "8080")
	m2 := mkMap("HOST", "prod.example.com")
	_, rows := Pivot([]map[string]envfile.Entry{m1, m2}, PivotOptions{})
	// rows should be sorted: HOST, PORT
	if rows[1].Key != "PORT" {
		t.Fatalf("expected PORT row, got %s", rows[1].Key)
	}
	if rows[1].Values[1] != "" {
		t.Error("expected empty string for missing key in second map")
	}
	if rows[1].Same {
		t.Error("expected Same=false when key is absent in one map")
	}
}

func TestPivot_SortedRows(t *testing.T) {
	m := mkMap("ZEBRA", "z", "ALPHA", "a", "MIDDLE", "m")
	_, rows := Pivot([]map[string]envfile.Entry{m}, PivotOptions{})
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, r := range rows {
		if r.Key != expected[i] {
			t.Errorf("row %d: expected %s, got %s", i, expected[i], r.Key)
		}
	}
}
