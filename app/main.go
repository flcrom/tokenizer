/**
 * @file app/main.go
 * @description Dispatches tokenizer commands for byte encoding, base decoding, and adjacent-pair counting.
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
	"sort"
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
		return errors.New("usage: tokenizer <encode|decode|pairs>")
	}

	switch arguments[0] {
	case "encode":
		return runEncode(arguments[1:], stdout)
	case "decode":
		return runDecode(arguments[1:], stdout)
	case "pairs":
		return runPairs(arguments[1:], stdout)
	default:
		return errors.New("usage: tokenizer <encode|decode|pairs>")
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

	ids, err := parseIDs(*encodedIDs)
	if err != nil {
		return fmt.Errorf("parse --ids: %w", err)
	}

	_, err = stdout.Write(decodeBaseIDs(ids))
	return err
}

func runPairs(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("pairs", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	encodedIDs := flags.String("ids", "", "JSON array of token IDs")
	if err := flags.Parse(arguments); err != nil {
		return errors.New("usage: tokenizer pairs --ids JSON")
	}
	if flags.NArg() != 0 {
		return errors.New("usage: tokenizer pairs --ids JSON")
	}

	ids, err := parseIDs(*encodedIDs)
	if err != nil {
		return fmt.Errorf("parse --ids: %w", err)
	}

	return json.NewEncoder(stdout).Encode(countAdjacentPairs(ids))
}

func parseIDs(encoded string) ([]int, error) {
	var ids []int
	if err := json.Unmarshal([]byte(encoded), &ids); err != nil {
		return nil, err
	}
	return ids, nil
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

type tokenPair struct {
	left  int
	right int
}

func countAdjacentPairs(ids []int) [][]int {
	counts := make(map[tokenPair]int)
	for index := 0; index+1 < len(ids); index++ {
		pair := tokenPair{left: ids[index], right: ids[index+1]}
		counts[pair]++
	}

	pairs := make([]tokenPair, 0, len(counts))
	for pair := range counts {
		pairs = append(pairs, pair)
	}
	sort.Slice(pairs, func(first, second int) bool {
		if pairs[first].left != pairs[second].left {
			return pairs[first].left < pairs[second].left
		}
		return pairs[first].right < pairs[second].right
	})

	result := make([][]int, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, []int{pair.left, pair.right, counts[pair]})
	}
	return result
}
