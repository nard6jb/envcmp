// Package lint provides checks for common .env file issues such as
// duplicate keys, empty values, and suspicious patterns.
package lint

import "fmt"

// Issue represents a single lint finding.
type Issue struct {
	Line    int
	Key     string
	Message string
}

func (i Issue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("line %d [%s]: %s", i.Line, i.Key, i.Message)
	}
	return fmt.Sprintf("[%s]: %s", i.Key, i.Message)
}

// Check runs all lint rules against the provided ordered entries.
// entries is a slice of [2]string{key, value} in parse order.
func Check(entries [][2]string) []Issue {
	var issues []Issue
	seen := make(map[string]int)

	for idx, kv := range entries {
		key, val := kv[0], kv[1]
		lineNum := idx + 1

		if prev, ok := seen[key]; ok {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("duplicate key (first seen at line %d)", prev),
			})
		} else {
			seen[key] = lineNum
		}

		if val == "" {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: "empty value",
			})
		}

		if containsWhitespace(val) {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: "value contains unquoted whitespace",
			})
		}
	}
	return issues
}

func containsWhitespace(s string) bool {
	for _, r := range s {
		if r == ' ' || r == '\t' {
			return true
		}
	}
	return false
}
