package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Issue describes a problem found during interpolation.
type Issue struct {
	Key      string
	Ref      string
	Message  string
}

// Result holds the resolved map and any issues encountered.
type Result struct {
	Resolved map[string]string
	Issues   []Issue
}

// Apply resolves variable references within env values using the same map.
// References to undefined variables are left unexpanded and recorded as issues.
func Apply(env map[string]string) Result {
	resolved := make(map[string]string, len(env))
	var issues []Issue

	for k, v := range env {
		expanded, errs := expand(k, v, env)
		resolved[k] = expanded
		issues = append(issues, errs...)
	}

	return Result{Resolved: resolved, Issues: issues}
}

// expand replaces variable references in val using the lookup map.
func expand(key, val string, lookup map[string]string) (string, []Issue) {
	var issues []Issue

	result := varPattern.ReplaceAllStringFunc(val, func(match string) string {
		ref := extractName(match)
		if replacement, ok := lookup[ref]; ok {
			return replacement
		}
		issues = append(issues, Issue{
			Key:     key,
			Ref:     ref,
			Message: fmt.Sprintf("undefined variable reference: %s", ref),
		})
		return match
	})

	return result, issues
}

// extractName strips ${ } or leading $ from a matched token.
func extractName(match string) string {
	match = strings.TrimPrefix(match, "${")
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
