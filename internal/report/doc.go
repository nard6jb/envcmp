// Package report provides formatting and rendering of environment diff results.
//
// It supports text output with optional ANSI color coding and secret masking.
// Results are sorted alphabetically by key for consistent, readable output.
//
// Usage:
//
//	opts := report.Options{Format: report.FormatText, NoColor: false}
//	report.Render(os.Stdout, results, opts)
package report
