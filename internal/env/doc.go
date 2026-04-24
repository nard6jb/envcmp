// Package env provides utilities for working with environment variable maps,
// including variable expansion, normalization, and cloning with transformations.
//
// # Expand
//
// Expand resolves ${VAR} and $VAR references within values, optionally falling
// back to OS environment variables and supporting strict mode.
//
// # Normalize
//
// Normalize applies bulk transformations such as trimming whitespace, removing
// empty values, and converting key casing.
//
// # Clone
//
// Clone produces a deep copy of an env map with optional prefix renaming,
// key overrides, and key removal — useful when projecting one environment's
// configuration into another namespace.
package env
