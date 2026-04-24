package mask

import "strings"

// sensitiveKeywords is the list of substrings that indicate a key holds
// a sensitive value. Comparisons are performed case-insensitively.
var sensitiveKeywords = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"_KEY",
	"APIKEY",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"ACCESS_KEY",
	"SIGNING",
}

// Masked is the placeholder string used in place of sensitive values.
const Masked = "***"

// IsSensitive reports whether the given key name contains any of the
// sensitiveKeywords. The comparison is case-insensitive.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}
