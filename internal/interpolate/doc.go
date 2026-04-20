// Package interpolate resolves variable references within .env file values.
//
// It supports both ${VAR_NAME} and $VAR_NAME syntax. References are resolved
// using the same env map (intra-file interpolation). Unresolved references are
// preserved verbatim and reported as Issues so callers can decide how to handle
// missing variables — e.g. warn the user or abort.
//
// Usage:
//
//	res := interpolate.Apply(envMap)
//	for _, issue := range res.Issues {
//		fmt.Println(issue.Message)
//	}
package interpolate
