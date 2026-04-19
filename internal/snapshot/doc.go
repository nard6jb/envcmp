// Package snapshot provides functionality for saving and loading .env file
// snapshots to disk, enabling drift detection between a known-good baseline
// and the current environment state.
//
// Snapshots are stored as JSON files with restricted file permissions (0600)
// to prevent accidental exposure of sensitive values.
//
// Basic usage:
//
//	// Save a snapshot
//	err := snapshot.Save("/path/to/snapshot.json", envMap)
//
//	// Load and compare
//	diffs, err := snapshot.CompareEnvAgainstSnapshot("/path/to/snapshot.json", envMap)
package snapshot
