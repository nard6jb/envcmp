// Package redact combines masking and filtering to produce redacted
// representations of environment variable maps.
//
// It wraps the mask package to annotate each key/value pair with
// redaction metadata, making it easy to track which values were
// hidden before display or export.
//
// Typical usage:
//
//	entries := redact.ApplyToMap(envMap)
//	maskedMap := redact.ToMap(entries)
//	redactedKeys := redact.RedactedKeys(entries)
package redact
