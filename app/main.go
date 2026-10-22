/**
 * @file app/main.go
 * @description Parses encode and decode commands and round-trips UTF-8 through base byte token IDs.
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
	if len(arguments) == 0 {
		return errors.New("usage: tokenizer <encode|decode>")
	}

	switch arguments[0] {
	case "encode":
		return runEncode(arguments[1:], stdout)
	case "decode":
		return runDecode(arguments[1:], stdout)
	default:
		return errors.New("usage: tokenizer <encode|decode>")
	}
}

func runEncode(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("encode", flag.ContinueOnError)
	// flags for the new "encode" option belong to flags variable
	flags.SetOutput(io.Discard)
	text := flags.String("text", "", "UTF-8 text to encode")
	if err := flags.Parse(arguments); err != nil {
		return errors.New("usage: tokenizer encode --text TEXT")
	}
	if flags.NArg() != 0 {
		return errors.New("usage: tokenizer encode --text TEXT")
	}

	return json.NewEncoder(stdout).Encode(encodeBytes(*text))
}

func runDecode(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("decode", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	encodedIDs := flags.String("ids", "", "JSON array of base token IDs")
	if err := flags.Parse(arguments); err != nil {
		return errors.New("usage: tokenizer decode --ids JSON")
	}
	if flags.NArg() != 0 {
		return errors.New("usage: tokenizer decode --ids JSON")
	}

	var ids []int
	if err := json.Unmarshal([]byte(*encodedIDs), &ids); err != nil {
		return fmt.Errorf("parse --ids: %w", err)
	}

	_, err := stdout.Write(decodeBaseIDs(ids))
	return err
}

func encodeBytes(text string) []int {
	bytes := []byte(text)
	ids := make([]int, len(bytes))
	for index, value := range bytes {
		ids[index] = int(value)
	}
	return ids
}

func decodeBaseIDs(ids []int) []byte {
	decoded := make([]byte, len(ids))
	for index, id := range ids {
		decoded[index] = byte(id)
	}
	return decoded
}
