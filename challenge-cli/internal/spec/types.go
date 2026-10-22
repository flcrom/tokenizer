/**
 * @file challenge-cli/internal/spec/types.go
 * @description Defines and validates the declarative stage, fixture, and expectation contracts used by the challenge CLI.
 * @depends_on none
 * @depended_on_by challenge-cli/internal/cli/cli.go, challenge-cli/internal/spec/loader.go, challenge-cli/internal/runner/runner.go, challenge-cli/internal/runner/runner_test.go, challenge-cli/stages/README.md
 */

package spec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

const CurrentVersion = 1

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Stage struct {
	Version int    `json:"version"`
	Stage   int    `json:"stage"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Cases   []Case `json:"cases"`
}

type Case struct {
	Name          string         `json:"name"`
	Args          []string       `json:"args,omitempty"`
	Stdin         *Payload       `json:"stdin,omitempty"`
	Files         []InputFile    `json:"files,omitempty"`
	Expect        Expectation    `json:"expect"`
	ExpectedFiles []ExpectedFile `json:"expected_files,omitempty"`
	TimeoutMS     int            `json:"timeout_ms,omitempty"`
	Repeat        int            `json:"repeat,omitempty"`
}

type Payload struct {
	Text   *string `json:"text,omitempty"`
	Base64 *string `json:"base64,omitempty"`
}

type InputFile struct {
	Path    string  `json:"path"`
	Content Payload `json:"content"`
}

type ExpectedFile struct {
	Path    string            `json:"path"`
	Content OutputExpectation `json:"content"`
}

type Expectation struct {
	Stdout   OutputExpectation `json:"stdout"`
	Stderr   OutputExpectation `json:"stderr"`
	ExitCode int               `json:"exit_code"`
}

type OutputExpectation struct {
	Mode   string          `json:"mode"`
	Text   string          `json:"text,omitempty"`
	Base64 string          `json:"base64,omitempty"`
	JSON   json.RawMessage `json:"json,omitempty"`
}

func (stage Stage) Validate() error {
	if stage.Version != CurrentVersion {
		return fmt.Errorf("unsupported stage version %d", stage.Version)
	}
	if stage.Stage <= 0 {
		return fmt.Errorf("stage number must be positive")
	}
	if !slugPattern.MatchString(stage.Slug) {
		return fmt.Errorf("slug %q must use lowercase words separated by hyphens", stage.Slug)
	}
	if strings.TrimSpace(stage.Title) == "" {
		return fmt.Errorf("stage title must not be empty")
	}

	names := make(map[string]struct{}, len(stage.Cases))
	for index, testCase := range stage.Cases {
		if err := testCase.Validate(); err != nil {
			return fmt.Errorf("case %d: %w", index+1, err)
		}
		if _, exists := names[testCase.Name]; exists {
			return fmt.Errorf("duplicate case name %q", testCase.Name)
		}
		names[testCase.Name] = struct{}{}
	}
	return nil
}

func (testCase Case) Validate() error {
	if strings.TrimSpace(testCase.Name) == "" {
		return fmt.Errorf("name must not be empty")
	}
	if testCase.TimeoutMS < 0 {
		return fmt.Errorf("timeout_ms must not be negative")
	}
	if testCase.Repeat < 0 {
		return fmt.Errorf("repeat must not be negative")
	}
	if testCase.Stdin != nil {
		if err := testCase.Stdin.Validate(); err != nil {
			return fmt.Errorf("stdin: %w", err)
		}
	}
	for _, file := range testCase.Files {
		if err := validateRelativePath(file.Path); err != nil {
			return fmt.Errorf("input file %q: %w", file.Path, err)
		}
		if err := file.Content.Validate(); err != nil {
			return fmt.Errorf("input file %q: %w", file.Path, err)
		}
	}
	if err := testCase.Expect.Stdout.Validate(); err != nil {
		return fmt.Errorf("stdout: %w", err)
	}
	if err := testCase.Expect.Stderr.Validate(); err != nil {
		return fmt.Errorf("stderr: %w", err)
	}
	for _, file := range testCase.ExpectedFiles {
		if err := validateRelativePath(file.Path); err != nil {
			return fmt.Errorf("expected file %q: %w", file.Path, err)
		}
		if err := file.Content.Validate(); err != nil {
			return fmt.Errorf("expected file %q: %w", file.Path, err)
		}
	}
	return nil
}

func (payload Payload) Validate() error {
	if payload.Text != nil && payload.Base64 != nil {
		return fmt.Errorf("text and base64 are mutually exclusive")
	}
	if payload.Text == nil && payload.Base64 == nil {
		return fmt.Errorf("one of text or base64 is required")
	}
	if payload.Base64 != nil {
		if _, err := base64.StdEncoding.DecodeString(*payload.Base64); err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}
	}
	return nil
}

func (payload Payload) Bytes() ([]byte, error) {
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	if payload.Text != nil {
		return []byte(*payload.Text), nil
	}
	return base64.StdEncoding.DecodeString(*payload.Base64)
}

func (expectation OutputExpectation) Validate() error {
	switch expectation.Mode {
	case "empty":
		if expectation.Text != "" || expectation.Base64 != "" || len(expectation.JSON) != 0 {
			return fmt.Errorf("empty mode cannot include a value")
		}
	case "exact_text", "contains_text":
		if expectation.Base64 != "" || len(expectation.JSON) != 0 {
			return fmt.Errorf("%s mode only accepts text", expectation.Mode)
		}
	case "exact_bytes":
		if expectation.Text != "" || len(expectation.JSON) != 0 {
			return fmt.Errorf("exact_bytes mode only accepts base64")
		}
		if _, err := base64.StdEncoding.DecodeString(expectation.Base64); err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}
	case "json":
		if expectation.Text != "" || expectation.Base64 != "" || len(expectation.JSON) == 0 {
			return fmt.Errorf("json mode requires only a json value")
		}
		var value any
		if err := json.Unmarshal(expectation.JSON, &value); err != nil {
			return fmt.Errorf("invalid json expectation: %w", err)
		}
	default:
		return fmt.Errorf("unsupported mode %q", expectation.Mode)
	}
	return nil
}

func validateRelativePath(path string) error {
	if path == "" || filepath.IsAbs(path) {
		return fmt.Errorf("path must be non-empty and relative")
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path must stay inside the case workspace")
	}
	return nil
}
