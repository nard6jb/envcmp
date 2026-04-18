// Package config handles CLI configuration and flag parsing for envcmp.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Mode represents the operating mode of the CLI.
type Mode string

const (
	ModeDiff     Mode = "diff"
	ModeValidate Mode = "validate"
)

// Config holds the resolved CLI configuration.
type Config struct {
	Mode      Mode
	FileA     string
	FileB     string
	MaskKeys  bool
	NoColor   bool
}

// Parse parses os.Args and returns a Config or an error.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envcmp", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	mode := fs.String("mode", "diff", "operating mode: diff | validate")
	mask := fs.Bool("mask", true, "mask sensitive values in output")
	noColor := fs.Bool("no-color", false, "disable colored output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return nil, errors.New("usage: envcmp [flags] <fileA> <fileB>")
	}

	m := Mode(*mode)
	if m != ModeDiff && m != ModeValidate {
		return nil, fmt.Errorf("unknown mode %q: must be diff or validate", *mode)
	}

	return &Config{
		Mode:     m,
		FileA:    positional[0],
		FileB:    positional[1],
		MaskKeys: *mask,
		NoColor:  *noColor,
	}, nil
}
