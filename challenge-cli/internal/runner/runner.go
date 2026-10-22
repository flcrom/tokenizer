/**
 * @file challenge-cli/internal/runner/runner.go
 * @description Executes stage cases against your_program.sh in isolated workspaces and compares observable results.
 * @depends_on your_program.sh, challenge-cli/internal/spec/types.go
 * @depended_on_by challenge-cli/internal/cli/cli.go, challenge-cli/internal/runner/runner_test.go
 */

package runner

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/flcrom/tokenizer/challenge-cli/internal/spec"
)

const defaultTimeout = 10 * time.Second

type Config struct {
	ProjectRoot string
	ProgramPath string
}

type Result struct {
	Passed   bool
	Failure  string
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
}

func RunCase(parent context.Context, config Config, testCase spec.Case) Result {
	started := time.Now()
	result := Result{ExitCode: -1}

	workspace, err := os.MkdirTemp("", "tokenizer-stage-")
	if err != nil {
		result.Failure = fmt.Sprintf("create case workspace: %v", err)
		return result
	}
	defer os.RemoveAll(workspace)

	if err := writeInputFiles(workspace, testCase.Files); err != nil {
		result.Failure = err.Error()
		return result
	}

	timeout := defaultTimeout
	if testCase.TimeoutMS > 0 {
		timeout = time.Duration(testCase.TimeoutMS) * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	programPath := config.ProgramPath
	if programPath == "" {
		programPath = filepath.Join(config.ProjectRoot, "your_program.sh")
	}
	command := exec.CommandContext(ctx, programPath, testCase.Args...)
	command.Dir = workspace
	command.Env = append(os.Environ(), "TOKENIZER_PROJECT_ROOT="+config.ProjectRoot)
	if testCase.Stdin != nil {
		stdin, err := testCase.Stdin.Bytes()
		if err != nil {
			result.Failure = fmt.Sprintf("decode stdin: %v", err)
			return result
		}
		command.Stdin = bytes.NewReader(stdin)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	result.Duration = time.Since(started)
	result.Stdout = append([]byte(nil), stdout.Bytes()...)
	result.Stderr = append([]byte(nil), stderr.Bytes()...)

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		result.Failure = fmt.Sprintf("program exceeded %s timeout", timeout)
		return result
	case err == nil:
		result.ExitCode = 0
	default:
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.Failure = fmt.Sprintf("start program: %v", err)
			return result
		}
	}

	if result.ExitCode != testCase.Expect.ExitCode {
		result.Failure = fmt.Sprintf("exit code: got %d, want %d", result.ExitCode, testCase.Expect.ExitCode)
		return result
	}
	if failure := compareOutput("stdout", result.Stdout, testCase.Expect.Stdout); failure != "" {
		result.Failure = failure
		return result
	}
	if failure := compareOutput("stderr", result.Stderr, testCase.Expect.Stderr); failure != "" {
		result.Failure = failure
		return result
	}
	for _, expectedFile := range testCase.ExpectedFiles {
		contents, err := os.ReadFile(filepath.Join(workspace, filepath.Clean(expectedFile.Path)))
		if err != nil {
			result.Failure = fmt.Sprintf("expected file %s: %v", expectedFile.Path, err)
			return result
		}
		if failure := compareOutput("file "+expectedFile.Path, contents, expectedFile.Content); failure != "" {
			result.Failure = failure
			return result
		}
	}

	result.Passed = true
	return result
}

func writeInputFiles(workspace string, files []spec.InputFile) error {
	for _, file := range files {
		contents, err := file.Content.Bytes()
		if err != nil {
			return fmt.Errorf("decode input file %s: %w", file.Path, err)
		}
		path := filepath.Join(workspace, filepath.Clean(file.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create directory for %s: %w", file.Path, err)
		}
		if err := os.WriteFile(path, contents, 0o644); err != nil {
			return fmt.Errorf("write input file %s: %w", file.Path, err)
		}
	}
	return nil
}

func compareOutput(label string, actual []byte, expectation spec.OutputExpectation) string {
	switch expectation.Mode {
	case "empty":
		if len(actual) != 0 {
			return fmt.Sprintf("%s: got %q, want empty", label, actual)
		}
	case "exact_text":
		if !bytes.Equal(actual, []byte(expectation.Text)) {
			return fmt.Sprintf("%s: got %q, want %q", label, actual, expectation.Text)
		}
	case "contains_text":
		if !bytes.Contains(actual, []byte(expectation.Text)) {
			return fmt.Sprintf("%s: %q does not contain %q", label, actual, expectation.Text)
		}
	case "exact_bytes":
		expected, err := base64.StdEncoding.DecodeString(expectation.Base64)
		if err != nil {
			return fmt.Sprintf("%s expectation: %v", label, err)
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Sprintf("%s bytes: got %v, want %v", label, actual, expected)
		}
	case "json":
		actualValue, err := decodeJSON(actual)
		if err != nil {
			return fmt.Sprintf("%s: invalid JSON: %v", label, err)
		}
		expectedValue, err := decodeJSON(expectation.JSON)
		if err != nil {
			return fmt.Sprintf("%s expectation: invalid JSON: %v", label, err)
		}
		if !reflect.DeepEqual(actualValue, expectedValue) {
			return fmt.Sprintf("%s JSON: got %s, want %s", label, strings.TrimSpace(string(actual)), string(expectation.JSON))
		}
	default:
		return fmt.Sprintf("%s: unsupported expectation mode %q", label, expectation.Mode)
	}
	return ""
}

func decodeJSON(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		return nil, fmt.Errorf("multiple JSON values")
	} else if !errors.Is(err, io.EOF) {
		return nil, err
	}
	return value, nil
}
