// Package profile provides named environment profiles, allowing users to
// define and switch between sets of env file paths (e.g. "staging", "prod").
package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named collection of env file paths.
type Profile struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// Store holds a set of named profiles persisted to disk.
type Store struct {
	Profiles []Profile `json:"profiles"`
}

// Add adds or replaces a profile in the store.
func (s *Store) Add(p Profile) {
	for i, existing := range s.Profiles {
		if existing.Name == p.Name {
			s.Profiles[i] = p
			return
		}
	}
	s.Profiles = append(s.Profiles, p)
}

// Get returns the profile with the given name, or an error if not found.
func (s *Store) Get(name string) (Profile, error) {
	for _, p := range s.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("profile %q not found", name)
}

// Remove deletes a profile by name. Returns false if it did not exist.
func (s *Store) Remove(name string) bool {
	for i, p := range s.Profiles {
		if p.Name == name {
			s.Profiles = append(s.Profiles[:i], s.Profiles[i+1:]...)
			return true
		}
	}
	return false
}

// SaveStore writes the store to the given file path as JSON.
func SaveStore(path string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("profile: mkdir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("profile: open: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadStore reads a store from the given file path.
// If the file does not exist, an empty store is returned.
func LoadStore(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("profile: read: %w", err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("profile: parse: %w", err)
	}
	return &s, nil
}
