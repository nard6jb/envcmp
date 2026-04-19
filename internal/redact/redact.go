// Package redact provides utilities for redacting sensitive env entries
// before output or export, combining mask and filter logic.
package redact

import (
	"github.com/yourusername/envcmp/internal/mask"
)

// Entry represents a single env key/value pair with redaction metadata.
type Entry struct {
	Key      string
	Value    string
	Redacted bool
}

// ApplyToMap takes a map of env vars and returns a slice of Entry values,
// masking any sensitive keys according to mask.IsSensitive.
func ApplyToMap(env map[string]string) []Entry {
	entries := make([]Entry, 0, len(env))
	for k, v := range env {
		if mask.IsSensitive(k) {
			entries = append(entries, Entry{Key: k, Value: mask.MaskValue(k, v), Redacted: true})
		} else {
			entries = append(entries, Entry{Key: k, Value: v, Redacted: false})
		}
	}
	return entries
}

// ToMap converts a slice of Entry back to a plain string map (masked values preserved).
func ToMap(entries []Entry) map[string]string {
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		out[e.Key] = e.Value
	}
	return out
}

// RedactedKeys returns the list of keys that were redacted.
func RedactedKeys(entries []Entry) []string {
	var keys []string
	for _, e := range entries {
		if e.Redacted {
			keys = append(keys, e.Key)
		}
	}
	return keys
}
