/**
 * @file challenge-cli/internal/spec/loader.go
 * @description Loads, strictly parses, validates, sorts, and selects stage manifests.
 * @depends_on challenge-cli/internal/spec/types.go
 * @depended_on_by challenge-cli/internal/cli/cli.go, challenge-cli/internal/spec/loader_test.go
 */

package spec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func LoadDirectory(directory string) ([]Stage, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("find stage manifests: %w", err)
	}
	sort.Strings(paths)

	stages := make([]Stage, 0, len(paths))
	seenNumbers := make(map[int]string, len(paths))
	seenSlugs := make(map[string]string, len(paths))
	for _, path := range paths {
		stage, err := loadFile(path)
		if err != nil {
			return nil, err
		}
		if previous, exists := seenNumbers[stage.Stage]; exists {
			return nil, fmt.Errorf("%s: stage number %d is already declared by %s", path, stage.Stage, previous)
		}
		if previous, exists := seenSlugs[stage.Slug]; exists {
			return nil, fmt.Errorf("%s: slug %q is already declared by %s", path, stage.Slug, previous)
		}
		seenNumbers[stage.Stage] = path
		seenSlugs[stage.Slug] = path
		stages = append(stages, stage)
	}

	sort.Slice(stages, func(i, j int) bool { return stages[i].Stage < stages[j].Stage })
	return stages, nil
}

func loadFile(path string) (Stage, error) {
	file, err := os.Open(path)
	if err != nil {
		return Stage{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var stage Stage
	if err := decoder.Decode(&stage); err != nil {
		return Stage{}, fmt.Errorf("parse %s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Stage{}, fmt.Errorf("parse %s: multiple JSON values", path)
		}
		return Stage{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := stage.Validate(); err != nil {
		return Stage{}, fmt.Errorf("validate %s: %w", path, err)
	}
	return stage, nil
}
