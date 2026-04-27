package env

import (
	"strings"

	"github.com/jabuta/envcmp/internal/envfile"
)

// TrimOptions controls how Trim behaves.
type TrimOptions struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// TrimPrefix removes a specific prefix from all keys.
	TrimPrefix string
	// TrimSuffix removes a specific suffix from all keys.
	TrimSuffix string
}

// TrimResult holds a single trimmed entry alongside what changed.
type TrimResult struct {
	OriginalKey   string
	OriginalValue string
	Key           string
	Value         string
	KeyChanged    bool
	ValueChanged  bool
}

// Trim applies whitespace and affix trimming to a map of env entries.
// It returns the cleaned map and a slice of results describing each change.
func Trim(entries map[string]string, opts TrimOptions) (map[string]string, []TrimResult) {
	out := make(map[string]string, len(entries))
	results := make([]TrimResult, 0, len(entries))

	for k, v := range entries {
		nk := k
		nv := v

		if opts.TrimKeys {
			nk = strings.TrimSpace(nk)
		}
		if opts.TrimPrefix != "" {
			nk = strings.TrimPrefix(nk, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			nk = strings.TrimSuffix(nk, opts.TrimSuffix)
		}
		if opts.TrimValues {
			nv = strings.TrimSpace(nv)
		}

		out[nk] = nv
		results = append(results, TrimResult{
			OriginalKey:   k,
			OriginalValue: v,
			Key:           nk,
			Value:         nv,
			KeyChanged:    nk != k,
			ValueChanged:  nv != v,
		})
	}

	_ = envfile.Entry{} // ensure import is referenced indirectly via shared types
	return out, results
}
