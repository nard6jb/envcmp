// Package config provides CLI flag parsing and configuration resolution
// for the envcmp tool.
//
// It exposes a single Parse function that accepts a slice of string arguments
// (typically os.Args[1:]) and returns a populated Config struct.
//
// Supported modes:
//   - diff:     compare two .env files and report differences
//   - validate: check a .env file against a reference for missing/extra keys
//
// Example:
//
//	cfg, err := config.Parse(os.Args[1:])
package config
