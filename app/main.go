/**
 * @file app/main.go
 * @description Parses the tokenizer command line and writes Stage 1 ASCII byte IDs as JSON.
 * @depends_on none
 * @depended_on_by your_program.sh, app/main_test.go
 */

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// take all the arguments other than the executable itself and pass it to the run function.
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string, stdout io.Writer) error {
	if len(arguments) == 0 || arguments[0] != "encode" {
		return errors.New("usage: tokenizer encode --text TEXT")
	}

	flags := flag.NewFlagSet("encode", flag.ContinueOnError)
	// flags for the new "encode" option belong to flags variable
	flags.SetOutput(io.Discard)
	text := flags.String("text", "", "ASCII text to encode")
	if err := flags.Parse(arguments[1:]); err != nil {
		return errors.New("usage: tokenizer encode --text TEXT")
	}
	if flags.NArg() != 0 {
		return errors.New("usage: tokenizer encode --text TEXT")
	}

	return json.NewEncoder(stdout).Encode(encodeBytes(*text))
}

func encodeBytes(text string) []int {
	bytes := []byte(text)
	ids := make([]int, len(bytes))
	for index, value := range bytes {
		ids[index] = int(value)
	}
	return ids
}
