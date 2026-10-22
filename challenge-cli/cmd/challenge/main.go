/**
 * @file challenge-cli/cmd/challenge/main.go
 * @description Starts the local challenge CLI and maps command failures to process exit status.
 * @depends_on challenge-cli/go.mod, challenge-cli/internal/cli/cli.go
 * @depended_on_by challenge, challenge-cli/README.md
 */

package main

import (
	"fmt"
	"os"

	"github.com/flcrom/tokenizer/challenge-cli/internal/cli"
)

func main() {
	projectRoot, err := cli.FindProjectRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	application := cli.Application{
		ProjectRoot: projectRoot,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
	}
	if code := application.Run(os.Args[1:]); code != 0 {
		os.Exit(code)
	}
}
