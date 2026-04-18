// Package export provides functionality to export diff and validation
// results to structured formats such as JSON and plain text.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envcmp/internal/diff"
	"github.com/user/envcmp/internal/validate"
)

// DiffExport is the JSON-serializable representation of a diff result.
type DiffExport struct {
	Entries []DiffEntry `json:"entries"`
}

// DiffEntry represents a single diff line.
type DiffEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	Left   string `json:"left,omitempty"`
	Right  string `json:"right,omitempty"`
}

// ValidationExport is the JSON-serializable representation of a validation result.
type ValidationExport struct {
	Valid       bool     `json:"valid"`
	MissingKeys []string `json:"missing_keys,omitempty"`
	ExtraKeys   []string `json:"extra_keys,omitempty"`
}

// WriteDiffJSON writes the diff result as JSON to w.
func WriteDiffJSON(w io.Writer, result map[string]diff.Entry) error {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]DiffEntry, 0, len(keys))
	for _, k := range keys {
		e := result[k]
		entries = append(entries, DiffEntry{
			Key:    k,
			Status: string(e.Status),
			Left:   e.Left,
			Right:  e.Right,
		})
	}

	return json.NewEncoder(w).Encode(DiffExport{Entries: entries})
}

// WriteValidationJSON writes the validation result as JSON to w.
func WriteValidationJSON(w io.Writer, result validate.Result) error {
	out := ValidationExport{
		Valid:       result.Valid,
		MissingKeys: result.MissingKeys,
		ExtraKeys:   result.ExtraKeys,
	}
	return json.NewEncoder(w).Encode(out)
}

// WriteDiffText writes a plain-text summary of the diff to w.
func WriteDiffText(w io.Writer, result map[string]diff.Entry) error {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		e := result[k]
		fmt.Fprintf(w, "%s\t%s\n", k, e.Status)
	}
	return nil
}
