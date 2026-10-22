/**
 * @file challenge-cli/internal/runner/runner_test.go
 * @description Verifies black-box process execution, output modes, files, failures, stdin, repetitions, and timeouts.
 * @depends_on challenge-cli/internal/runner/runner.go, challenge-cli/internal/spec/types.go, challenge-cli/testdata/programs/fake_program.sh
 * @depended_on_by none
 */

package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/flcrom/tokenizer/challenge-cli/internal/spec"
)

func TestRunCaseComparesTextAndStderr(t *testing.T) {
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name: "echo",
		Args: []string{"echo", "hello"},
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "exact_text", Text: "hello"},
			Stderr:   spec.OutputExpectation{Mode: "contains_text", Text: "diagnostic:hello"},
			ExitCode: 0,
		},
	})
	assertPassed(t, result)
}

func TestRunCaseComparesJSONStructurally(t *testing.T) {
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name: "json",
		Args: []string{"json"},
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "json", JSON: json.RawMessage(`{"a":1,"b":2}`)},
			Stderr:   spec.OutputExpectation{Mode: "empty"},
			ExitCode: 0,
		},
	})
	assertPassed(t, result)
}

func TestRunCasePreservesStdinBytes(t *testing.T) {
	raw := []byte{0, 10, 255}
	encoded := base64.StdEncoding.EncodeToString(raw)
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name:  "stdin",
		Args:  []string{"stdin"},
		Stdin: &spec.Payload{Base64: &encoded},
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "exact_bytes", Base64: encoded},
			Stderr:   spec.OutputExpectation{Mode: "empty"},
			ExitCode: 0,
		},
	})
	assertPassed(t, result)
}

func TestRunCaseChecksExitCode(t *testing.T) {
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name: "failure",
		Args: []string{"fail"},
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "empty"},
			Stderr:   spec.OutputExpectation{Mode: "json", JSON: json.RawMessage(`{"error":"bad"}`)},
			ExitCode: 7,
		},
	})
	assertPassed(t, result)
}

func TestRunCaseChecksCreatedFiles(t *testing.T) {
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name: "file",
		Args: []string{"write", "model", "output/model.json"},
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "empty"},
			Stderr:   spec.OutputExpectation{Mode: "empty"},
			ExitCode: 0,
		},
		Files: []spec.InputFile{{
			Path:    "output/.keep",
			Content: spec.Payload{Text: stringPointer("")},
		}},
		ExpectedFiles: []spec.ExpectedFile{{
			Path:    "output/model.json",
			Content: spec.OutputExpectation{Mode: "exact_text", Text: "model"},
		}},
	})
	assertPassed(t, result)
}

func TestRunCaseEnforcesTimeout(t *testing.T) {
	result := RunCase(context.Background(), fixtureConfig(t), spec.Case{
		Name:      "timeout",
		Args:      []string{"timeout"},
		TimeoutMS: 25,
		Expect: spec.Expectation{
			Stdout:   spec.OutputExpectation{Mode: "empty"},
			Stderr:   spec.OutputExpectation{Mode: "empty"},
			ExitCode: 0,
		},
	})
	if result.Passed || result.Failure == "" {
		t.Fatalf("RunCase() result = %#v, want timeout failure", result)
	}
}

func fixtureConfig(t *testing.T) Config {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "programs", "fake_program.sh"))
	if err != nil {
		t.Fatal(err)
	}
	return Config{ProjectRoot: t.TempDir(), ProgramPath: path}
}

func assertPassed(t *testing.T, result Result) {
	t.Helper()
	if !result.Passed {
		t.Fatalf("RunCase() failure = %s, stdout=%q stderr=%q exit=%d", result.Failure, result.Stdout, result.Stderr, result.ExitCode)
	}
}

func stringPointer(value string) *string {
	return &value
}
