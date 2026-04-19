// Package audit provides change tracking for env file comparisons.
package audit

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Command   string            `json:"command"`
	Files     []string          `json:"files"`
	Changes   map[string]string `json:"changes,omitempty"`
	Issues    []string          `json:"issues,omitempty"`
}

// Log holds a list of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Append adds a new entry to the audit log file at path.
// If the file does not exist it is created.
func Append(path string, e Entry) error {
	l, err := load(path)
	if err != nil {
		l = &Log{}
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	l.Entries = append(l.Entries, e)
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadLog reads and returns the audit log from path.
func LoadLog(path string) (*Log, error) {
	return load(path)
}

func load(path string) (*Log, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}
