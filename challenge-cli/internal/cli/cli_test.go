/**
 * @file challenge-cli/internal/cli/cli_test.go
 * @description Verifies stage scaffolding, listing, selection, and CLI argument failures without touching Git.
 * @depends_on challenge-cli/internal/cli/cli.go
 * @depended_on_by none
 */

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewStageAndList(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	application := Application{ProjectRoot: root, Stdout: &stdout, Stderr: &stderr}

	code := application.Run([]string{"new-stage", "--number", "1", "--slug", "encode-ascii", "--title", "Encode ASCII"})
	if code != 0 {
		t.Fatalf("new-stage code = %d, stderr = %q", code, stderr.String())
	}
	path := filepath.Join(root, "challenge-cli", "stages", "01-encode-ascii.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("generated stage: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	code = application.Run([]string{"list"})
	if code != 0 {
		t.Fatalf("list code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "encode-ascii") || !strings.Contains(stdout.String(), "draft") {
		t.Fatalf("list output = %q", stdout.String())
	}
}

func TestNewStageDoesNotOverwrite(t *testing.T) {
	root := t.TempDir()
	application := Application{ProjectRoot: root, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}}
	arguments := []string{"new-stage", "--number", "1", "--slug", "encode-ascii"}
	if code := application.Run(arguments); code != 0 {
		t.Fatalf("first new-stage code = %d", code)
	}
	if code := application.Run(arguments); code == 0 {
		t.Fatal("second new-stage unexpectedly overwrote the manifest")
	}
}

func TestTestRequiresDefinedStage(t *testing.T) {
	root := t.TempDir()
	var stderr bytes.Buffer
	application := Application{ProjectRoot: root, Stdout: &bytes.Buffer{}, Stderr: &stderr}
	if code := application.Run([]string{"test", "--stage", "1"}); code != 2 {
		t.Fatalf("test code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "not defined") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
