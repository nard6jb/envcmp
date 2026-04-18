package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestMain_DiffNoDiff(t *testing.T) {
	a := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	cmd := buildCmd(t, "diff", a, b)
	if err := cmd.Run(); err != nil {
		t.Errorf("expected exit 0, got: %v", err)
	}
}

func TestMain_DiffHasDiff(t *testing.T) {
	a := writeTempFile(t, "FOO=bar\n")
	b := writeTempFile(t, "FOO=changed\n")
	cmd := buildCmd(t, "diff", a, b)
	err := cmd.Run()
	if err == nil {
		t.Error("expected non-zero exit when diff found")
	}
}

func TestMain_ValidateValid(t *testing.T) {
	ref := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	target := writeTempFile(t, "FOO=x\nBAZ=y\n")
	cmd := buildCmd(t, "validate", ref, target)
	if err := cmd.Run(); err != nil {
		t.Errorf("expected exit 0, got: %v", err)
	}
}

func TestMain_ValidateMissing(t *testing.T) {
	ref := writeTempFile(t, "FOO=bar\nREQUIRED=x\n")
	target := writeTempFile(t, "FOO=bar\n")
	cmd := buildCmd(t, "validate", ref, target)
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit for missing keys")
	}
}

func buildCmd(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "envcmp")
	build := exec.Command("go", "build", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	cmd := exec.Command(binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
