/**
 * @file challenge-cli/internal/cli/cli.go
 * @description Implements stage creation, discovery, selection, execution, and terminal reporting for the local challenge CLI.
 * @depends_on challenge-cli/internal/spec/loader.go, challenge-cli/internal/spec/types.go, challenge-cli/internal/runner/runner.go
 * @depended_on_by challenge-cli/cmd/challenge/main.go, challenge-cli/internal/cli/cli_test.go, challenge-cli/stages/README.md
 */

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/flcrom/tokenizer/challenge-cli/internal/runner"
	"github.com/flcrom/tokenizer/challenge-cli/internal/spec"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Application struct {
	ProjectRoot string
	Stdout      io.Writer
	Stderr      io.Writer
}

func FindProjectRoot() (string, error) {
	if configured := os.Getenv("TOKENIZER_PROJECT_ROOT"); configured != "" {
		return filepath.Abs(configured)
	}
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("determine working directory: %w", err)
	}
	for {
		if fileExists(filepath.Join(current, "your_program.sh")) && directoryExists(filepath.Join(current, "challenge-cli")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("could not find tokenizer project root")
		}
		current = parent
	}
}

func (application Application) Run(arguments []string) int {
	if len(arguments) == 0 {
		application.printUsage()
		return 2
	}
	switch arguments[0] {
	case "help", "-h", "--help":
		application.printUsage()
		return 0
	case "list":
		return application.list(arguments[1:])
	case "new-stage":
		return application.newStage(arguments[1:])
	case "test":
		return application.test(arguments[1:])
	default:
		fmt.Fprintf(application.Stderr, "unknown command %q\n", arguments[0])
		application.printUsage()
		return 2
	}
}

func (application Application) list(arguments []string) int {
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.SetOutput(application.Stderr)
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(application.Stderr, "list does not accept positional arguments")
		return 2
	}

	stages, err := spec.LoadDirectory(application.stageDirectory())
	if err != nil {
		fmt.Fprintln(application.Stderr, err)
		return 2
	}
	if len(stages) == 0 {
		fmt.Fprintln(application.Stdout, "No stages defined.")
		return 0
	}
	for _, stage := range stages {
		status := "ready"
		if len(stage.Cases) == 0 {
			status = "draft"
		}
		fmt.Fprintf(application.Stdout, "%02d  %-24s  %s  (%d cases)\n", stage.Stage, stage.Slug, status, len(stage.Cases))
	}
	return 0
}

func (application Application) newStage(arguments []string) int {
	flags := flag.NewFlagSet("new-stage", flag.ContinueOnError)
	flags.SetOutput(application.Stderr)
	number := flags.Int("number", 0, "positive stage number")
	slug := flags.String("slug", "", "lowercase hyphenated stage slug")
	title := flags.String("title", "", "human-readable stage title")
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 || *number <= 0 || !slugPattern.MatchString(*slug) {
		fmt.Fprintln(application.Stderr, "usage: ./challenge new-stage --number N --slug lowercase-slug [--title TITLE]")
		return 2
	}
	if *title == "" {
		words := strings.Split(*slug, "-")
		for index, word := range words {
			words[index] = strings.ToUpper(word[:1]) + word[1:]
		}
		*title = strings.Join(words, " ")
	}
	stage := spec.Stage{
		Version: spec.CurrentVersion,
		Stage:   *number,
		Slug:    *slug,
		Title:   *title,
		Cases:   []spec.Case{},
	}
	if err := stage.Validate(); err != nil {
		fmt.Fprintln(application.Stderr, err)
		return 2
	}
	contents, err := json.MarshalIndent(stage, "", "  ")
	if err != nil {
		fmt.Fprintln(application.Stderr, err)
		return 2
	}
	contents = append(contents, '\n')
	if err := os.MkdirAll(application.stageDirectory(), 0o755); err != nil {
		fmt.Fprintln(application.Stderr, err)
		return 2
	}
	path := filepath.Join(application.stageDirectory(), fmt.Sprintf("%02d-%s.json", *number, *slug))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		fmt.Fprintf(application.Stderr, "create stage: %v\n", err)
		return 2
	}
	if _, err := file.Write(contents); err != nil {
		file.Close()
		fmt.Fprintf(application.Stderr, "write stage: %v\n", err)
		return 2
	}
	if err := file.Close(); err != nil {
		fmt.Fprintf(application.Stderr, "close stage: %v\n", err)
		return 2
	}
	fmt.Fprintf(application.Stdout, "Created %s\n", path)
	return 0
}

func (application Application) test(arguments []string) int {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.SetOutput(application.Stderr)
	stageNumber := flags.Int("stage", 0, "stage number to test")
	previous := flags.Bool("previous", false, "run all stages through the selected stage")
	caseName := flags.String("case", "", "run one named case from the selected stage")
	debug := flags.Bool("debug", false, "show captured program output for every case")
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 || *stageNumber <= 0 {
		fmt.Fprintln(application.Stderr, "usage: ./challenge test --stage N [--previous] [--case NAME] [--debug]")
		return 2
	}
	if *caseName != "" && *previous {
		fmt.Fprintln(application.Stderr, "--case cannot be combined with --previous")
		return 2
	}

	stages, err := spec.LoadDirectory(application.stageDirectory())
	if err != nil {
		fmt.Fprintln(application.Stderr, err)
		return 2
	}
	selected := selectStages(stages, *stageNumber, *previous)
	if len(selected) == 0 || selected[len(selected)-1].Stage != *stageNumber {
		fmt.Fprintf(application.Stderr, "stage %d is not defined\n", *stageNumber)
		return 2
	}

	total := 0
	failed := 0
	for _, stage := range selected {
		cases := stage.Cases
		if stage.Stage == *stageNumber && *caseName != "" {
			cases = filterCases(cases, *caseName)
			if len(cases) == 0 {
				fmt.Fprintf(application.Stderr, "stage %d has no case named %q\n", stage.Stage, *caseName)
				return 2
			}
		}
		if len(cases) == 0 {
			fmt.Fprintf(application.Stderr, "stage %d is a draft with no cases\n", stage.Stage)
			return 2
		}
		fmt.Fprintf(application.Stdout, "[stage-%02d] %s\n", stage.Stage, stage.Title)
		for _, testCase := range cases {
			repetitions := testCase.Repeat
			if repetitions == 0 {
				repetitions = 1
			}
			for attempt := 1; attempt <= repetitions; attempt++ {
				total++
				result := runner.RunCase(context.Background(), runner.Config{ProjectRoot: application.ProjectRoot}, testCase)
				label := testCase.Name
				if repetitions > 1 {
					label = fmt.Sprintf("%s (run %d/%d)", label, attempt, repetitions)
				}
				if result.Passed {
					fmt.Fprintf(application.Stdout, "  PASS  %s (%s)\n", label, result.Duration.Round(time.Millisecond))
				} else {
					failed++
					fmt.Fprintf(application.Stdout, "  FAIL  %s: %s\n", label, result.Failure)
				}
				if *debug || !result.Passed {
					printCaptured(application.Stdout, result.Stdout, result.Stderr)
				}
			}
		}
	}
	if failed > 0 {
		fmt.Fprintf(application.Stdout, "Failed: %d of %d cases.\n", failed, total)
		return 1
	}
	fmt.Fprintf(application.Stdout, "Passed: %d cases.\n", total)
	return 0
}

func selectStages(stages []spec.Stage, number int, previous bool) []spec.Stage {
	selected := make([]spec.Stage, 0)
	for _, stage := range stages {
		if stage.Stage == number || (previous && stage.Stage < number) {
			selected = append(selected, stage)
		}
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].Stage < selected[j].Stage })
	return selected
}

func filterCases(cases []spec.Case, name string) []spec.Case {
	for _, testCase := range cases {
		if testCase.Name == name {
			return []spec.Case{testCase}
		}
	}
	return nil
}

func printCaptured(writer io.Writer, stdout []byte, stderr []byte) {
	if len(stdout) > 0 {
		fmt.Fprintf(writer, "    stdout: %q\n", stdout)
	}
	if len(stderr) > 0 {
		fmt.Fprintf(writer, "    stderr: %q\n", stderr)
	}
}

func (application Application) stageDirectory() string {
	return filepath.Join(application.ProjectRoot, "challenge-cli", "stages")
}

func (application Application) printUsage() {
	fmt.Fprintln(application.Stderr, "Tokenizer challenge CLI")
	fmt.Fprintln(application.Stderr, "  ./challenge list")
	fmt.Fprintln(application.Stderr, "  ./challenge new-stage --number N --slug SLUG [--title TITLE]")
	fmt.Fprintln(application.Stderr, "  ./challenge test --stage N [--previous] [--case NAME] [--debug]")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
