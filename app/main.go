/**
 * @file app/main.go
 * @description Dispatches byte encoding, base decoding, pair counting, and safely terminating deterministic tokenizer training.
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
		return errors.New("usage: tokenizer <encode|decode|pairs|train>")
	}

	switch arguments[0] {
	case "encode":
		return runEncode(arguments[1:], stdout)
	case "decode":
		return runDecode(arguments[1:], stdout)
	case "pairs":
		return runPairs(arguments[1:], stdout)
	case "train":
		return runTrain(arguments[1:], stdout)
	default:
		return errors.New("usage: tokenizer <encode|decode|pairs|train>")
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

func runTrain(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("train", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	text := flags.String("text", "", "UTF-8 training text")
	mergeLimit := flags.Int("merges", 0, "maximum number of merge rules to learn")
	if err := flags.Parse(arguments); err != nil {
		return errors.New("usage: tokenizer train --text TEXT --merges N")
	}
	mergeLimitProvided := false
	flags.Visit(func(parsedFlag *flag.Flag) {
		if parsedFlag.Name == "merges" {
			mergeLimitProvided = true
		}
	})
	if flags.NArg() != 0 || !mergeLimitProvided || *mergeLimit < 0 {
		return errors.New("usage: tokenizer train --text TEXT --merges N")
	}

	result, err := trainMerges(*text, *mergeLimit)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(result.model)
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

func pairFrequencies(ids []int) map[tokenPair]int {
	counts := make(map[tokenPair]int)
	for index := 0; index+1 < len(ids); index++ {
		pair := tokenPair{left: ids[index], right: ids[index+1]}
		counts[pair]++
	}
	return counts
}

func countAdjacentPairs(ids []int) [][]int {
	counts := pairFrequencies(ids)

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

const baseVocabSize = 256

var errNoRepeatedPair = errors.New("no pair occurs at least twice")

type tokenizerModel struct {
	Version       int     `json:"version"`
	BaseVocabSize int     `json:"base_vocab_size"`
	Merges        [][]int `json:"merges"`
}

type trainingResult struct {
	model tokenizerModel
	ids   []int
}

type mergeStep struct {
	pair  tokenPair
	count int
	newID int
	ids   []int
}

func trainFirstMerge(text string) (tokenizerModel, error) {
	result, err := trainMerges(text, 1)
	return result.model, err
}

func firstMergeStep(text string) (mergeStep, error) {
	return nextMergeStep(encodeBytes(text), baseVocabSize)
}

func trainMerges(text string, mergeLimit int) (trainingResult, error) {
	ids := encodeBytes(text)
	merges := make([][]int, 0, mergeLimit)
	for mergeIndex := 0; mergeIndex < mergeLimit; mergeIndex++ {
		step, err := nextMergeStep(ids, baseVocabSize+mergeIndex)
		if errors.Is(err, errNoRepeatedPair) {
			break
		}
		if err != nil {
			return trainingResult{}, err
		}
		merges = append(merges, []int{step.pair.left, step.pair.right})
		ids = step.ids
	}

	return trainingResult{
		model: tokenizerModel{
			Version:       1,
			BaseVocabSize: baseVocabSize,
			Merges:        merges,
		},
		ids: ids,
	}, nil
}

func nextMergeStep(ids []int, newID int) (mergeStep, error) {
	pair, count, err := selectTopPair(pairFrequencies(ids))
	if err != nil {
		return mergeStep{}, err
	}

	return mergeStep{
		pair:  pair,
		count: count,
		newID: newID,
		ids:   replacePair(ids, pair, newID),
	}, nil
}

func selectTopPair(counts map[tokenPair]int) (tokenPair, int, error) {
	pairs := make([]tokenPair, 0, len(counts))
	for pair := range counts {
		pairs = append(pairs, pair)
	}
	sort.Slice(pairs, func(first, second int) bool {
		firstCount := counts[pairs[first]]
		secondCount := counts[pairs[second]]
		if firstCount != secondCount {
			return firstCount > secondCount
		}
		if pairs[first].left != pairs[second].left {
			return pairs[first].left < pairs[second].left
		}
		return pairs[first].right < pairs[second].right
	})

	if len(pairs) == 0 || counts[pairs[0]] < 2 {
		return tokenPair{}, 0, errNoRepeatedPair
	}
	selected := pairs[0]
	return selected, counts[selected], nil
}

func replacePair(ids []int, pair tokenPair, newID int) []int {
	result := make([]int, 0, len(ids))
	for index := 0; index < len(ids); {
		if index+1 < len(ids) && ids[index] == pair.left && ids[index+1] == pair.right {
			result = append(result, newID)
			index += 2
			continue
		}
		result = append(result, ids[index])
		index++
	}
	return result
}
