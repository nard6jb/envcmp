package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from an .env file.
type Entry struct {
	Key   string
	Value string
	Line  int
}

// EnvFile holds all parsed entries from an .env file.
type EnvFile struct {
	Path    string
	Entries map[string]Entry
}

// Parse reads and parses an .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:    path,
		Entries: make(map[string]Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("%s:%d: invalid line (missing '='): %q", path, lineNum, line)
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"`)

		if key == "" {
			return nil, fmt.Errorf("%s:%d: empty key", path, lineNum)
		}

		env.Entries[key] = Entry{Key: key, Value: value, Line: lineNum}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}

	return env, nil
}
