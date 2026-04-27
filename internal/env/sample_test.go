package env

import (
	"testing"

	"github.com/subtlepseudonym/envcmp/internal/envfile"
)

func sampleEntries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "GAMMA", Value: "3"},
		{Key: "DELTA", Value: "4"},
		{Key: "EPSILON", Value: "5"},
	}
}

func TestSample_NoOptions(t *testing.T) {
	entries := sampleEntries()
	result := Sample(entries, SampleOptions{})
	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestSample_NLessThanTotal(t *testing.T) {
	entries := sampleEntries()
	result := Sample(entries, SampleOptions{N: 3, Seed: 42})
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestSample_NGreaterThanTotal(t *testing.T) {
	entries := sampleEntries()
	result := Sample(entries, SampleOptions{N: 100, Seed: 1})
	if len(result) != len(entries) {
		t.Errorf("expected all %d entries, got %d", len(entries), len(result))
	}
}

func TestSample_Sorted(t *testing.T) {
	entries := sampleEntries()
	result := Sample(entries, SampleOptions{N: 4, Seed: 7, Sorted: true})
	for i := 1; i < len(result); i++ {
		if result[i].Key < result[i-1].Key {
			t.Errorf("result not sorted: %s before %s", result[i-1].Key, result[i].Key)
		}
	}
}

func TestSample_Reproducible(t *testing.T) {
	entries := sampleEntries()
	opts := SampleOptions{N: 2, Seed: 99}
	a := Sample(entries, opts)
	b := Sample(entries, opts)
	if len(a) != len(b) {
		t.Fatalf("lengths differ: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Key != b[i].Key {
			t.Errorf("position %d: got %q and %q with same seed", i, a[i].Key, b[i].Key)
		}
	}
}

func TestSample_EmptyInput(t *testing.T) {
	result := Sample([]envfile.Entry{}, SampleOptions{N: 3})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestSample_DoesNotMutateInput(t *testing.T) {
	entries := sampleEntries()
	original := make([]envfile.Entry, len(entries))
	copy(original, entries)
	Sample(entries, SampleOptions{N: 2, Seed: 5})
	for i := range entries {
		if entries[i].Key != original[i].Key {
			t.Errorf("input mutated at index %d: got %q, want %q", i, entries[i].Key, original[i].Key)
		}
	}
}
