// Package config parses CLI arguments into a structured Config.
package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Mode represents the operating mode of envcmp.
type Mode string

const (
	ModeDiff     Mode = "diff"
	ModeValidate Mode = "validate"
)

// Config holds all parsed CLI configuration.
type Config struct {
	Mode     Mode
	Files    []string
	NoColor  bool
	Prefix   string
	Exclude  []string
	OnlyKeys []string
}

// Parse reads os.Args and returns a Config or an error.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envcmp", flag.ContinueOnError)

	noColor := fs.Bool("no-color", false, "disable colored output")
	prefix := fs.String("prefix", "", "only include keys with this prefix")
	exclude := fs.String("exclude", "", "comma-separated keys to exclude")
	onlyKeys := fs.String("keys", "", "comma-separated keys to include exclusively")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return nil, errors.New("usage: envcmp <mode> <file1> [file2]")
	}

	mode := Mode(remaining[0])
	if mode != ModeDiff && mode != ModeValidate {
		return nil, fmt.Errorf("unknown mode %q: use 'diff' or 'validate'", mode)
	}

	files := remaining[1:]
	if mode == ModeDiff && len(files) < 2 {
		return nil, errors.New("diff mode requires two files")
	}

	cfg := &Config{
		Mode:    mode,
		Files:   files,
		NoColor: *noColor,
		Prefix:  *prefix,
	}
	if *exclude != "" {
		cfg.Exclude = splitCSV(*exclude)
	}
	if *onlyKeys != "" {
		cfg.OnlyKeys = splitCSV(*onlyKeys)
	}
	return cfg, nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
