// Package validate checks a target .env file against a reference template.
//
// It reports which keys are missing (present in reference but not in target)
// and which keys are extra (present in target but not in reference).
//
// A validation result is considered Valid only when no required keys are
// missing. Extra keys are reported as informational warnings.
package validate
