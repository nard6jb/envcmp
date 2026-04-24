package env

import (
	"github.com/subtlepseudonym/envcmp/internal/diff"
	"github.com/subtlepseudonym/envcmp/internal/mask"
)

// maskEntry replaces sensitive values inside a diff.Entry with "***".
func maskEntry(e diff.Entry) diff.Entry {
	if !mask.IsSensitive(e.Key) {
		return e
	}
	if e.Left != "" {
		e.Left = mask.MaskValue(e.Key, e.Left)
	}
	if e.Right != "" {
		e.Right = mask.MaskValue(e.Key, e.Right)
	}
	return e
}
