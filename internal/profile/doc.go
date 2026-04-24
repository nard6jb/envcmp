// Package profile manages named environment profiles for envcmp.
//
// A profile is a named set of .env file paths that can be saved and
// loaded from a local JSON store (e.g. ~/.config/envcmp/profiles.json).
// This allows users to define shortcuts for commonly compared environments
// such as "staging", "prod", or "ci" without repeating file paths on
// every invocation.
//
// Example usage:
//
//	envcmp profile add staging .env.staging .env.shared
//	envcmp diff --profile staging
package profile
