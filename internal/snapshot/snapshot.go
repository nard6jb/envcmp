// Package snapshot provides functionality to save and load env file snapshots
// for later comparison against current state.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a saved state of an env file.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// Save writes the given env entries to a snapshot file at the given path.
func Save(path, source string, entries map[string]string) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	if snap.Entries == nil {
		snap.Entries = make(map[string]string)
	}

	return &snap, nil
}
