// Package mask provides utilities for masking sensitive values
// in environment variable output.
package mask

import "strings"

const defaultMask = "********"

// SecretKeyPatterns holds substrings that indicate a key is sensitive.
var SecretKeyPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// IsSensitive reports whether the given key name looks like a secret.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range SecretKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked placeholder if the key is sensitive,
// otherwise it returns the original value unchanged.
func MaskValue(key, value string) string {
	if IsSensitive(key) {
		return defaultMask
	}
	return value
}

// MaskMap returns a copy of the provided map with sensitive values masked.
func MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = MaskValue(k, v)
	}
	return out
}
