package main

import (
	"fmt"
	"os"

	"github.com/user/envcmp/internal/config"
	"github.com/user/envcmp/internal/diff"
	"github.com/user/envcmp/internal/envfile"
	"github.com/user/envcmp/internal/mask"
	"github.com/user/envcmp/internal/report"
	"github.com/user/envcmp/internal/validate"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch cfg.Mode {
	case config.ModeDiff:
		runDiff(cfg)
	case config.ModeValidate:
	\t}
}

func runDiff(cfg *config.Config) {
	left, err := envfile.Parse(cfg.Files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", cfg.Files[0], err)
		os.Exit(1)
	}
	right, err := envfile.Parse(cfg.Files[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", cfg.Files[1], err)
		os.Exit(1)
	}

	if cfg.MaskSecrets {
		left = mask.MaskMap(left)
		right = mask.MaskMap(right)
	}

	result := diff.Compare(left, right)
	report.Render(result, cfg.Color)

	if diff.HasDiff(result) {
		os.Exit(1)
	}
}

func runValidate(cfg *config.Config) {
	ref, err := envfile.Parse(cfg.Files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading reference %s: %v\n", cfg.Files[0], err)
		os.Exit(1)
	}
	target, err := envfile.Parse(cfg.Files[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target %s: %v\n", cfg.Files[1], err)
		os.Exit(1)
	}

	result := validate.Against(ref, target)
	report.RenderValidation(result, cfg.Color)

	if len(result.Missing) > 0 {
		os.Exit(1)
	}
}
