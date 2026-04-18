// Package runner wires together parsing, diffing/validating, and reporting.
package runner

import (
	"fmt"
	"io"

	"github.com/user/envcmp/internal/config"
	"github.com/user/envcmp/internal/diff"
	"github.com/user/envcmp/internal/envfile"
	"github.com/user/envcmp/internal/mask"
	"github.com/user/envcmp/internal/report"
	"github.com/user/envcmp/internal/validate"
)

// Result holds the exit code after running.
type Result struct {
	Code int
}

// Run executes the appropriate mode and writes output to w.
func Run(cfg *config.Config, w io.Writer) Result {
	switch cfg.Mode {
	case config.ModeDiff:
		return runDiff(cfg, w)
	case config.ModeValidate:
		return runValidate(cfg, w)
	default:
		fmt.Fprintf(w, "unknown mode: %s\n", cfg.Mode)
		return Result{Code: 2}
	}
}

func runDiff(cfg *config.Config, w io.Writer) Result {
	left, err := envfile.Parse(cfg.Files[0])
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", cfg.Files[0], err)
		return Result{Code: 2}
	}
	right, err := envfile.Parse(cfg.Files[1])
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", cfg.Files[1], err)
		return Result{Code: 2}
	}

	if cfg.MaskSecrets {
		left = mask.MaskMap(left)
		right = mask.MaskMap(right)
	}

	result := diff.Compare(left, right)
	report.Render(result, cfg.Color)

	if diff.HasDiff(result) {
		return Result{Code: 1}
	}
	return Result{Code: 0}
}

func runValidate(cfg *config.Config, w io.Writer) Result {
	ref, err := envfile.Parse(cfg.Files[0])
	if err != nil {
		fmt.Fprintf(w, "error reading reference %s: %v\n", cfg.Files[0], err)
		return Result{Code: 2}
	}
	target, err := envfile.Parse(cfg.Files[1])
	if err != nil {
		fmt.Fprintf(w, "error reading target %s: %v\n", cfg.Files[1], err)
		return Result{Code: 2}
	}

	result := validate.Against(ref, target)
	report.RenderValidation(result, cfg.Color)

	if len(result.Missing) > 0 {
		return Result{Code: 1}
	}
	return Result{Code: 0}
}
