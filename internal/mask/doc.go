// Package mask provides secret-detection and value-masking helpers for
// environment variable keys and maps.
//
// Keys are considered sensitive when their names contain well-known
// substrings such as SECRET, PASSWORD, TOKEN, API_KEY, AUTH, PRIVATE,
// or CREDENTIAL (case-insensitive).
//
// Typical usage:
//
//	masked := mask.MaskMap(parsedEnv)
//	fmt.Println(masked["DB_PASSWORD"]) // prints ********
//
// The package is intentionally stateless and safe for concurrent use.
package mask
