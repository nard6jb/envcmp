// Package config parses and validates CLI arguments for envcmp.
//
// It supports two modes of operation:
//
//	- diff:     compare two .env files and report differences
//	- validate: check a target .env file against a reference for missing keys
//
// Flags such as --no-color and --mask-secrets are also handled here.
package config
