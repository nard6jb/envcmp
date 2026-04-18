// Package envfile provides functionality for parsing .env files into
// structured representations suitable for comparison and validation.
//
// Supported syntax:
//   - KEY=VALUE
//   - KEY="VALUE WITH SPACES"
//   - Lines starting with '#' are treated as comments and ignored.
//   - Blank lines are ignored.
//
// Example usage:
//
//	env, err := envfile.Parse(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for key, entry := range env.Entries {
//		fmt.Printf("%s = %s (line %d)\n", key, entry.Value, entry.Line)
//	}
package envfile
