// Package lint implements static analysis rules for .env files.
//
// It detects common authoring mistakes including:
//   - Duplicate keys within the same file
//   - Empty values that may indicate missing configuration
//   - Unquoted values containing whitespace
//
// Usage:
//
//	issues := lint.Check(entries)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package lint
