// Package mask provides utilities for detecting and masking sensitive
// environment variable values based on key name heuristics.
//
// Keys containing substrings such as "SECRET", "PASSWORD", "TOKEN", "KEY",
// or "PRIVATE" (case-insensitive) are considered sensitive and their values
// will be replaced with a redacted placeholder when masking is enabled.
//
// Example usage:
//
//	if mask.IsSensitive("DB_PASSWORD") {
//		fmt.Println(mask.MaskValue("DB_PASSWORD", "s3cr3t"))
//	}
package mask
