/**
 * @file challenge-cli/internal/spec/loader_test.go
 * @description Verifies strict stage loading, ordering, validation, and duplicate detection.
 * @depends_on challenge-cli/internal/spec/loader.go
 * @depended_on_by none
 */

package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDirectorySortsStages(t *testing.T) {
	directory := t.TempDir()
	writeManifest(t, directory, "02-second.json", `{
  "version": 1,
  "stage": 2,
  "slug": "second",
  "title": "Second",
  "cases": []
}`)
	writeManifest(t, directory, "01-first.json", `{
  "version": 1,
  "stage": 1,
  "slug": "first",
  "title": "First",
  "cases": []
}`)

	stages, err := LoadDirectory(directory)
	if err != nil {
		t.Fatalf("LoadDirectory() error = %v", err)
	}
	if len(stages) != 2 || stages[0].Stage != 1 || stages[1].Stage != 2 {
		t.Fatalf("LoadDirectory() stages = %#v", stages)
	}
}

func TestLoadDirectoryRejectsUnknownFields(t *testing.T) {
	directory := t.TempDir()
	writeManifest(t, directory, "01-first.json", `{
  "version": 1,
  "stage": 1,
  "slug": "first",
  "title": "First",
  "surprise": true,
  "cases": []
}`)

	_, err := LoadDirectory(directory)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("LoadDirectory() error = %v, want unknown field", err)
	}
}

func TestLoadDirectoryRejectsDuplicateNumbers(t *testing.T) {
	directory := t.TempDir()
	writeManifest(t, directory, "01-first.json", `{"version":1,"stage":1,"slug":"first","title":"First","cases":[]}`)
	writeManifest(t, directory, "01-other.json", `{"version":1,"stage":1,"slug":"other","title":"Other","cases":[]}`)

	_, err := LoadDirectory(directory)
	if err == nil || !strings.Contains(err.Error(), "already declared") {
		t.Fatalf("LoadDirectory() error = %v, want duplicate error", err)
	}
}

func writeManifest(t *testing.T, directory string, name string, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
