// Package export provides serialization of diff and validation results
// to structured output formats.
//
// Supported formats:
//   - JSON: machine-readable output suitable for CI pipelines
//   - Text: tab-separated plain text for scripting
//
// Usage:
//
//	export.WriteDiffJSON(os.Stdout, diffResult)
//	export.WriteValidationJSON(os.Stdout, validationResult)
package export
