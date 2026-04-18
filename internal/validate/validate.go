// Package validate provides functionality for validating .env files
// against a reference template, ensuring all required keys are present.
package validate

import "github.com/your/envcmp/internal/diff"

// Result holds the outcome of a validation check.
type Result struct {
	MissingKeys []string
	ExtraKeys   []string
	Valid       bool
}

// Against validates the target map against a reference (template) map.
// Missing keys are keys present in reference but absent in target.
// Extra keys are keys present in target but absent in reference.
func Against(reference, target map[string]string) Result {
	entries := diff.Compare(reference, target)

	var missing, extra []string
	for _, e := range entries {
		switch e.Status {
		case diff.Missing:
			missing = append(missing, e.Key)
		case diff.Extra:
			extra = append(extra, e.Key)
		}
	}

	return Result{
		MissingKeys: missing,
		ExtraKeys:   extra,
		Valid:       len(missing) == 0,
	}
}
