package env

import (
	"testing"
)

func TestIntersect_OnlyCommonKeys(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"B": "20", "C": "3", "D": "4"}

	result := Intersect(left, right, IntersectOptions{})

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "B" || result[1].Key != "C" {
		t.Errorf("unexpected keys: %v", result)
	}
}

func TestIntersect_ValuesPreserved(t *testing.T) {
	left := map[string]string{"HOST": "localhost"}
	right := map[string]string{"HOST": "prod.example.com"}

	result := Intersect(left, right, IntersectOptions{})

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Left != "localhost" || result[0].Right != "prod.example.com" {
		t.Errorf("unexpected values: %+v", result[0])
	}
}

func TestIntersect_EmptyMaps(t *testing.T) {
	result := Intersect(map[string]string{}, map[string]string{}, IntersectOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestIntersect_NoCommonKeys(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"B": "2"}

	result := Intersect(left, right, IntersectOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestIntersect_MaskSecrets(t *testing.T) {
	left := map[string]string{"SECRET_KEY": "abc123", "HOST": "localhost"}
	right := map[string]string{"SECRET_KEY": "xyz789", "HOST": "prod"}

	result := Intersect(left, right, IntersectOptions{MaskSecrets: true})

	for _, e := range result {
		if e.Key == "SECRET_KEY" {
			if !e.Masked {
				t.Errorf("expected SECRET_KEY to be masked")
			}
			if e.Left != "***" || e.Right != "***" {
				t.Errorf("expected masked values, got left=%q right=%q", e.Left, e.Right)
			}
		}
		if e.Key == "HOST" && e.Masked {
			t.Errorf("HOST should not be masked")
		}
	}
}

func TestIntersect_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "2", "M": "3"}
	right := map[string]string{"Z": "1", "A": "2", "M": "3"}

	result := Intersect(left, right, IntersectOptions{})

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Key != "A" || result[1].Key != "M" || result[2].Key != "Z" {
		t.Errorf("expected sorted output, got %v", result)
	}
}
