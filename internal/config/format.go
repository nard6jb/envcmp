package config

import "fmt"

// Format represents an output format for envcmp results.
type Format string

const (
	// FormatText is the default human-readable coloured output.
	FormatText Format = "text"
	// FormatJSON emits machine-readable JSON.
	FormatJSON Format = "json"
)

// ParseFormat converts a raw string to a Format value.
// Returns an error for unrecognised format strings.
func ParseFormat(s string) (Format, error) {
	switch s {
	case string(FormatText), "":
		return FormatText, nil
	case string(FormatJSON):
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of [text, json]", s)
	}
}

// IsJSON returns true when the format is JSON.
func (f Format) IsJSON() bool { return f == FormatJSON }
