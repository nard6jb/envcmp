// Package snapshot provides functionality for saving and loading
// environment file snapshots to disk in JSON format.
//
// Snapshots allow users to capture the state of an environment at a
// point in time and later compare it against a current environment
// to detect drift or unintended changes.
//
// Usage:
//
//	err := snapshot.Save("/path/to/snapshot.json", envMap)
//	envMap, err := snapshot.Load("/path/to/snapshot.json")
package snapshot
