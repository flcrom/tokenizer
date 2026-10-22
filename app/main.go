/**
 * @file app/main.go
 * @description Dispatches validated tokenizer operations with structured errors, learned-model inference, diagnostics, training, and persistence.
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
		_ = json.NewEncoder(os.Stderr).Encode(errorResponse{Error: codeForError(err)})
		os.Exit(1)
	}
}

type errorCode string

const (
	errInvalidInput errorCode = "invalid_input"
	errIO           errorCode = "io_error"
	errInvalidModel errorCode = "invalid_model"
	errUnknownToken errorCode = "unknown_token"
)

type tokenizerError struct {
	code  errorCode
	cause error
}

func (err *tokenizerError) Error() string {
	return err.cause.Error()
}

func (err *tokenizerError) Unwrap() error {
	return err.cause
}

type errorResponse struct {
	Error errorCode `json:"error"`
}

func classifiedError(code errorCode, format string, arguments ...any) error {
	return &tokenizerError{code: code, cause: fmt.Errorf(format, arguments...)}
}

func codeForError(err error) errorCode {
	var classified *tokenizerError
	if errors.As(err, &classified) {
		return classified.code
	}
	return errInvalidInput
}

func run(arguments []string, stdout io.Writer) error {
	if len(arguments) == 0 {
		return classifiedError(errInvalidInput, "usage: tokenizer <encode|decode|pairs|train>")
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
		return classifiedError(errInvalidInput, "usage: tokenizer <encode|decode|pairs|train>")
	}
}

func runEncode(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("encode", flag.ContinueOnError)
	// flags for the new "encode" option belong to flags variable
	flags.SetOutput(io.Discard)
	text := flags.String("text", "", "UTF-8 text to encode")
	modelPath := flags.String("model", "", "path to a tokenizer model")
	if err := flags.Parse(arguments); err != nil {
		return classifiedError(errInvalidInput, "usage: tokenizer encode [--model PATH] --text TEXT")
	}

	modelProvided := false
	textProvided := false
	flags.Visit(func(parsedFlag *flag.Flag) {
		switch parsedFlag.Name {
		case "model":
			modelProvided = true
		case "text":
			textProvided = true
		}
	})
	if flags.NArg() != 0 || !textProvided {
		return classifiedError(errInvalidInput, "usage: tokenizer encode [--model PATH] --text TEXT")
	}
	if modelProvided {
		model, err := loadModel(*modelPath)
		if err != nil {
			return err
		}
		ids, err := encodeWithModel(*text, model)
		if err != nil {
			return err
		}
		if err := json.NewEncoder(stdout).Encode(ids); err != nil {
			return classifiedError(errIO, "write stdout: %w", err)
		}
		return nil
	}

	if err := json.NewEncoder(stdout).Encode(encodeBytes(*text)); err != nil {
		return classifiedError(errIO, "write stdout: %w", err)
	}
	return nil
}

func runDecode(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("decode", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	encodedIDs := flags.String("ids", "", "JSON array of base token IDs")
	modelPath := flags.String("model", "", "path to a tokenizer model")
	if err := flags.Parse(arguments); err != nil {
		return classifiedError(errInvalidInput, "usage: tokenizer decode [--model PATH] --ids JSON")
	}
	idsProvided := false
	modelProvided := false
	flags.Visit(func(parsedFlag *flag.Flag) {
		switch parsedFlag.Name {
		case "ids":
			idsProvided = true
		case "model":
			modelProvided = true
		}
	})
	if flags.NArg() != 0 || !idsProvided {
		return classifiedError(errInvalidInput, "usage: tokenizer decode [--model PATH] --ids JSON")
	}

	ids, err := parseIDs(*encodedIDs)
	if err != nil {
		return err
	}

	if modelProvided {
		model, err := loadModel(*modelPath)
		if err != nil {
			return err
		}
		decoded, err := decodeWithModel(ids, model)
		if err != nil {
			return err
		}
		if _, err = stdout.Write(decoded); err != nil {
			return classifiedError(errIO, "write stdout: %w", err)
		}
		return nil
	}

	decoded, err := decodeBaseIDs(ids)
	if err != nil {
		return err
	}
	if _, err = stdout.Write(decoded); err != nil {
		return classifiedError(errIO, "write stdout: %w", err)
	}
	return nil
}

func runPairs(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("pairs", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	encodedIDs := flags.String("ids", "", "JSON array of token IDs")
	if err := flags.Parse(arguments); err != nil {
		return classifiedError(errInvalidInput, "usage: tokenizer pairs --ids JSON")
	}
	idsProvided := false
	flags.Visit(func(parsedFlag *flag.Flag) {
		if parsedFlag.Name == "ids" {
			idsProvided = true
		}
	})
	if flags.NArg() != 0 || !idsProvided {
		return classifiedError(errInvalidInput, "usage: tokenizer pairs --ids JSON")
	}

	ids, err := parseIDs(*encodedIDs)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(stdout).Encode(countAdjacentPairs(ids)); err != nil {
		return classifiedError(errIO, "write stdout: %w", err)
	}
	return nil
}

func runTrain(arguments []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("train", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	text := flags.String("text", "", "UTF-8 training text")
	inputPath := flags.String("input", "", "path to a training corpus")
	outputPath := flags.String("output", "", "path for the trained model")
	mergeLimit := flags.Int("merges", 0, "maximum number of merge rules to learn")
	if err := flags.Parse(arguments); err != nil {
		return classifiedError(errInvalidInput, "usage: tokenizer train <--text TEXT|--input PATH> --merges N [--output PATH]")
	}
	mergeLimitProvided := false
	textProvided := false
	inputProvided := false
	outputProvided := false
	flags.Visit(func(parsedFlag *flag.Flag) {
		switch parsedFlag.Name {
		case "merges":
			mergeLimitProvided = true
		case "text":
			textProvided = true
		case "input":
			inputProvided = true
		case "output":
			outputProvided = true
		}
	})
	if flags.NArg() != 0 || !mergeLimitProvided || *mergeLimit < 0 || textProvided == inputProvided {
		return classifiedError(errInvalidInput, "usage: tokenizer train <--text TEXT|--input PATH> --merges N [--output PATH]")
	}

	var corpus []byte
	var err error
	if inputProvided {
		corpus, err = os.ReadFile(*inputPath)
		if err != nil {
			return classifiedError(errIO, "read --input: %w", err)
		}
	} else {
		corpus = []byte(*text)
	}

	result, err := trainBytes(corpus, *mergeLimit)
	if err != nil {
		return err
	}
	encodedModel, err := json.Marshal(result.model)
	if err != nil {
		return classifiedError(errInvalidModel, "encode model: %w", err)
	}
	encodedModel = append(encodedModel, '\n')

	if outputProvided {
		if err := os.WriteFile(*outputPath, encodedModel, 0o644); err != nil {
			return classifiedError(errIO, "write --output: %w", err)
		}
	}
	if _, err = stdout.Write(encodedModel); err != nil {
		return classifiedError(errIO, "write stdout: %w", err)
	}
	return nil
}

func parseIDs(encoded string) ([]int, error) {
	var ids []int
	if err := json.Unmarshal([]byte(encoded), &ids); err != nil {
		return nil, classifiedError(errInvalidInput, "parse --ids: %w", err)
	}
	if ids == nil {
		return nil, classifiedError(errInvalidInput, "parse --ids: expected an array")
	}
	for _, id := range ids {
		if id < 0 {
			return nil, classifiedError(errInvalidInput, "parse --ids: token IDs must be nonnegative")
		}
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

func decodeBaseIDs(ids []int) ([]byte, error) {
	decoded := make([]byte, len(ids))
	for index, id := range ids {
		if id < 0 || id >= baseVocabSize {
			return nil, classifiedError(errUnknownToken, "unknown token ID %d", id)
		}
		decoded[index] = byte(id)
	}
	return decoded, nil
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

func loadModel(path string) (tokenizerModel, error) {
	encodedModel, err := os.ReadFile(path)
	if err != nil {
		return tokenizerModel{}, classifiedError(errIO, "read --model: %w", err)
	}
	return parseModel(encodedModel)
}

func parseModel(encodedModel []byte) (tokenizerModel, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(encodedModel, &object); err != nil || object == nil {
		return tokenizerModel{}, classifiedError(errInvalidModel, "parse --model: malformed JSON")
	}

	versionJSON, versionPresent := object["version"]
	baseSizeJSON, baseSizePresent := object["base_vocab_size"]
	mergesJSON, mergesPresent := object["merges"]
	if !versionPresent || !baseSizePresent || !mergesPresent {
		return tokenizerModel{}, classifiedError(errInvalidModel, "parse --model: missing required field")
	}

	var model tokenizerModel
	if err := json.Unmarshal(versionJSON, &model.Version); err != nil {
		return tokenizerModel{}, classifiedError(errInvalidModel, "parse --model version: %w", err)
	}
	if err := json.Unmarshal(baseSizeJSON, &model.BaseVocabSize); err != nil {
		return tokenizerModel{}, classifiedError(errInvalidModel, "parse --model base_vocab_size: %w", err)
	}
	if err := json.Unmarshal(mergesJSON, &model.Merges); err != nil || model.Merges == nil {
		return tokenizerModel{}, classifiedError(errInvalidModel, "parse --model merges: expected an array")
	}
	if err := validateModel(model); err != nil {
		return tokenizerModel{}, err
	}
	return model, nil
}

func validateModel(model tokenizerModel) error {
	if model.Version != 1 || model.BaseVocabSize != baseVocabSize {
		return classifiedError(errInvalidModel, "unsupported model metadata")
	}

	seenPairs := make(map[tokenPair]struct{}, len(model.Merges))
	for mergeIndex, encodedPair := range model.Merges {
		if len(encodedPair) != 2 {
			return classifiedError(errInvalidModel, "model merge %d must contain exactly two token IDs", mergeIndex)
		}
		newID := model.BaseVocabSize + mergeIndex
		left, right := encodedPair[0], encodedPair[1]
		if left < 0 || left >= newID || right < 0 || right >= newID {
			return classifiedError(errInvalidModel, "model merge %d references an unavailable token ID", mergeIndex)
		}
		pair := tokenPair{left: left, right: right}
		if _, exists := seenPairs[pair]; exists {
			return classifiedError(errInvalidModel, "model merge %d duplicates an earlier pair", mergeIndex)
		}
		seenPairs[pair] = struct{}{}
	}
	return nil
}

func encodeWithModel(text string, model tokenizerModel) ([]int, error) {
	ids := encodeBytes(text)
	for mergeIndex, encodedPair := range model.Merges {
		if len(encodedPair) != 2 {
			return nil, classifiedError(errInvalidModel, "model merge %d must contain exactly two token IDs", mergeIndex)
		}
		pair := tokenPair{left: encodedPair[0], right: encodedPair[1]}
		ids = replacePair(ids, pair, model.BaseVocabSize+mergeIndex)
	}
	return ids, nil
}

func buildTokenExpansions(model tokenizerModel) ([][]byte, error) {
	expansions := make([][]byte, model.BaseVocabSize+len(model.Merges))
	for id := 0; id < model.BaseVocabSize; id++ {
		expansions[id] = []byte{byte(id)}
	}

	for mergeIndex, encodedPair := range model.Merges {
		if len(encodedPair) != 2 {
			return nil, classifiedError(errInvalidModel, "model merge %d must contain exactly two token IDs", mergeIndex)
		}
		newID := model.BaseVocabSize + mergeIndex
		left, right := encodedPair[0], encodedPair[1]
		if left < 0 || left >= newID || right < 0 || right >= newID {
			return nil, classifiedError(errInvalidModel, "model merge %d references an unavailable token ID", mergeIndex)
		}

		expansion := make([]byte, 0, len(expansions[left])+len(expansions[right]))
		expansion = append(expansion, expansions[left]...)
		expansion = append(expansion, expansions[right]...)
		expansions[newID] = expansion
	}
	return expansions, nil
}

func decodeWithModel(ids []int, model tokenizerModel) ([]byte, error) {
	expansions, err := buildTokenExpansions(model)
	if err != nil {
		return nil, err
	}

	decoded := make([]byte, 0)
	for _, id := range ids {
		if id < 0 || id >= len(expansions) {
			return nil, classifiedError(errUnknownToken, "unknown token ID %d", id)
		}
		decoded = append(decoded, expansions[id]...)
	}
	return decoded, nil
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
	return trainBytes([]byte(text), mergeLimit)
}

func trainBytes(corpus []byte, mergeLimit int) (trainingResult, error) {
	ids := make([]int, len(corpus))
	for index, value := range corpus {
		ids[index] = int(value)
	}
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
