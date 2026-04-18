package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/your/envcmp/internal/validate"
)

// RenderValidation writes a human-readable validation report to w.
func RenderValidation(w io.Writer, result validate.Result, filename string) {
	if result.Valid && len(result.ExtraKeys) == 0 {
		fmt.Fprintf(w, "✓ %s is valid — all required keys present\n", filename)
		return
	}

	if !result.Valid {
		sort.Strings(result.MissingKeys)
		for _, k := range result.MissingKeys {
			printColored(w, colorRed, fmt.Sprintf("MISSING  %s", k))
		}
	}

	if len(result.ExtraKeys) > 0 {
		sort.Strings(result.ExtraKeys)
		for _, k := range result.ExtraKeys {
			printColored(w, colorYellow, fmt.Sprintf("EXTRA    %s", k))
		}
	}

	if !result.Valid {
		fmt.Fprintf(w, "✗ %s is invalid — %d key(s) missing\n", filename, len(result.MissingKeys))
	} else {
		fmt.Fprintf(w, "⚠ %s has %d extra key(s)\n", filename, len(result.ExtraKeys))
	}
}
