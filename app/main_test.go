/**
 * @file app/main_test.go
 * @description Verifies model validation, error classification, learned-model inference, pair operations, training, and raw byte preservation.
 * @depends_on app/main.go
 * @depended_on_by none
 */

package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []int
	}{
		{name: "single byte", text: "A", want: []int{65}},
		{name: "ordered bytes", text: "hello", want: []int{104, 101, 108, 108, 111}},
		{name: "space", text: "A A", want: []int{65, 32, 65}},
		{name: "multibyte text", text: "café", want: []int{99, 97, 102, 195, 169}},
		{name: "emoji", text: "🙂", want: []int{240, 159, 153, 130}},
		{name: "empty", text: "", want: []int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := encodeBytes(test.text); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("encodeBytes(%q) = %v, want %v", test.text, got, test.want)
			}
		})
	}
}

func TestDecodeBaseIDs(t *testing.T) {
	tests := []struct {
		name string
		ids  []int
		want []byte
	}{
		{name: "UTF-8 bytes", ids: []int{99, 97, 102, 195, 169}, want: []byte("café")},
		{name: "embedded newline", ids: []int{65, 10, 66}, want: []byte{'A', '\n', 'B'}},
		{name: "arbitrary bytes", ids: []int{0, 10, 255}, want: []byte{0, 10, 255}},
		{name: "empty", ids: []int{}, want: []byte{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := decodeBaseIDs(test.ids)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, test.want) {
				t.Fatalf("decodeBaseIDs(%v) = %v, want %v", test.ids, got, test.want)
			}
		})
	}
}

func TestBaseTokenRoundTrip(t *testing.T) {
	fixtures := []string{"", "hello", "café", "🙂", "é", "line one\nline two"}
	for _, fixture := range fixtures {
		decoded, err := decodeBaseIDs(encodeBytes(fixture))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(decoded, []byte(fixture)) {
			t.Fatalf("round trip for %q produced %v", fixture, decoded)
		}
	}
}

func TestCountAdjacentPairs(t *testing.T) {
	tests := []struct {
		name string
		ids  []int
		want [][]int
	}{
		{name: "repeated pairs", ids: []int{97, 98, 97, 98}, want: [][]int{{97, 98, 2}, {98, 97, 1}}},
		{name: "overlapping windows", ids: []int{7, 7, 7, 7}, want: [][]int{{7, 7, 3}}},
		{name: "numeric pair order", ids: []int{256, 97, 256, 97}, want: [][]int{{97, 256, 1}, {256, 97, 2}}},
		{name: "empty", ids: []int{}, want: [][]int{}},
		{name: "one ID", ids: []int{42}, want: [][]int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := countAdjacentPairs(test.ids); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("countAdjacentPairs(%v) = %v, want %v", test.ids, got, test.want)
			}
		})
	}
}

func TestReplacePair(t *testing.T) {
	tests := []struct {
		name  string
		ids   []int
		pair  tokenPair
		newID int
		want  []int
	}{
		{name: "overlap is consumed once", ids: []int{97, 97, 97}, pair: tokenPair{97, 97}, newID: 256, want: []int{256, 97}},
		{name: "two disjoint matches", ids: []int{97, 97, 97, 97}, pair: tokenPair{97, 97}, newID: 256, want: []int{256, 256}},
		{name: "preserves IDs between matches", ids: []int{97, 98, 99, 97, 98}, pair: tokenPair{97, 98}, newID: 256, want: []int{256, 99, 256}},
		{name: "absent pair", ids: []int{1, 2, 3}, pair: tokenPair{8, 9}, newID: 256, want: []int{1, 2, 3}},
		{name: "empty", ids: []int{}, pair: tokenPair{1, 2}, newID: 256, want: []int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := replacePair(test.ids, test.pair, test.newID); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("replacePair(%v, %v, %d) = %v, want %v", test.ids, test.pair, test.newID, got, test.want)
			}
		})
	}
}

func TestSelectTopPair(t *testing.T) {
	tests := []struct {
		name      string
		counts    map[tokenPair]int
		wantPair  tokenPair
		wantCount int
	}{
		{
			name:      "highest frequency",
			counts:    map[tokenPair]int{{97, 98}: 3, {98, 97}: 2},
			wantPair:  tokenPair{97, 98},
			wantCount: 3,
		},
		{
			name:      "smaller left ID breaks tie",
			counts:    map[tokenPair]int{{98, 97}: 2, {97, 122}: 2},
			wantPair:  tokenPair{97, 122},
			wantCount: 2,
		},
		{
			name:      "smaller right ID breaks tie after left",
			counts:    map[tokenPair]int{{97, 99}: 2, {97, 98}: 2},
			wantPair:  tokenPair{97, 98},
			wantCount: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotPair, gotCount, err := selectTopPair(test.counts)
			if err != nil {
				t.Fatal(err)
			}
			if gotPair != test.wantPair || gotCount != test.wantCount {
				t.Fatalf("selectTopPair(%v) = (%v, %d), want (%v, %d)", test.counts, gotPair, gotCount, test.wantPair, test.wantCount)
			}
		})
	}
}

func TestFirstMergeStep(t *testing.T) {
	tests := []struct {
		text string
		want mergeStep
	}{
		{text: "ababab", want: mergeStep{pair: tokenPair{97, 98}, count: 3, newID: 256, ids: []int{256, 256, 256}}},
		{text: "zzzzz", want: mergeStep{pair: tokenPair{122, 122}, count: 4, newID: 256, ids: []int{256, 256, 122}}},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			got, err := firstMergeStep(test.text)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("firstMergeStep(%q) = %#v, want %#v", test.text, got, test.want)
			}
		})
	}
}

func TestTrainFirstMerge(t *testing.T) {
	want := tokenizerModel{
		Version:       1,
		BaseVocabSize: 256,
		Merges:        [][]int{{97, 98}},
	}

	got, err := trainFirstMerge("ababab")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trainFirstMerge(%q) = %#v, want %#v", "ababab", got, want)
	}
}

func TestTrainMerges(t *testing.T) {
	want := trainingResult{
		model: tokenizerModel{
			Version:       1,
			BaseVocabSize: 256,
			Merges:        [][]int{{97, 110}, {98, 256}, {256, 97}},
		},
		ids: []int{257, 258, 32, 257, 258},
	}

	got, err := trainMerges("banana banana", 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trainMerges(%q, %d) = %#v, want %#v", "banana banana", 3, got, want)
	}
}

func TestTrainMergesStopsSafely(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		mergeLimit int
		want       trainingResult
	}{
		{
			name:       "zero budget",
			text:       "abab",
			mergeLimit: 0,
			want: trainingResult{
				model: tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{}},
				ids:   []int{97, 98, 97, 98},
			},
		},
		{
			name:       "no repeated pair",
			text:       "abc",
			mergeLimit: 10,
			want: trainingResult{
				model: tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{}},
				ids:   []int{97, 98, 99},
			},
		},
		{
			name:       "stops below budget",
			text:       "aaaaa",
			mergeLimit: 10,
			want: trainingResult{
				model: tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{97, 97}}},
				ids:   []int{256, 256, 97},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := trainMerges(test.text, test.mergeLimit)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("trainMerges(%q, %d) = %#v, want %#v", test.text, test.mergeLimit, got, test.want)
			}
		})
	}
}

func TestTrainBytesPreservesArbitraryCorpusBytes(t *testing.T) {
	want := trainingResult{
		model: tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{255, 0}}},
		ids:   []int{256, 256},
	}

	got, err := trainBytes([]byte{255, 0, 255, 0}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trainBytes() = %#v, want %#v", got, want)
	}
}

func TestLoadModel(t *testing.T) {
	path := t.TempDir() + "/model.json"
	if err := os.WriteFile(path, []byte("{\n  \"merges\": [], \"base_vocab_size\": 256, \"version\": 1\n}"), 0o644); err != nil {
		t.Fatal(err)
	}

	want := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{}}
	got, err := loadModel(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loadModel(%q) = %#v, want %#v", path, got, want)
	}
}

func TestEncodeWithOneLearnedPair(t *testing.T) {
	model := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{97, 98}}}
	tests := []struct {
		name string
		text string
		want []int
	}{
		{name: "pair throughout new text", text: "zab!ab", want: []int{122, 256, 33, 256}},
		{name: "pair occurs once", text: "ab", want: []int{256}},
		{name: "uncovered bytes", text: "z!", want: []int{122, 33}},
		{name: "empty", text: "", want: []int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := encodeWithModel(test.text, model)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("encodeWithModel(%q) = %v, want %v", test.text, got, test.want)
			}
		})
	}
}

func TestEncodeWithOneLearnedPairDoesNotOverlap(t *testing.T) {
	model := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{97, 97}}}
	got, err := encodeWithModel("aaa", model)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{256, 97}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encodeWithModel(%q) = %v, want %v", "aaa", got, want)
	}
}

func TestEncodeWithMultipleBasePairsUsesModelOrder(t *testing.T) {
	model := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{97, 98}, {98, 99}}}
	got, err := encodeWithModel("abc bc bc", model)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{256, 99, 32, 257, 32, 257}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encodeWithModel() = %v, want %v", got, want)
	}

	reversed := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{98, 99}, {97, 98}}}
	got, err = encodeWithModel("abc bc bc", reversed)
	if err != nil {
		t.Fatal(err)
	}
	want = []int{97, 256, 32, 256, 32, 256}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("encodeWithModel() with reversed rules = %v, want %v", got, want)
	}
}

func TestEncodeWithModelAppliesHierarchicalPairs(t *testing.T) {
	model := tokenizerModel{
		Version:       1,
		BaseVocabSize: 256,
		Merges:        [][]int{{97, 110}, {98, 256}, {256, 97}},
	}
	tests := []struct {
		text string
		want []int
	}{
		{text: "bana ana ana", want: []int{257, 97, 32, 258, 32, 258}},
		{text: "banana!", want: []int{257, 258, 33}},
		{text: "banana banana", want: []int{257, 258, 32, 257, 258}},
	}

	for _, test := range tests {
		got, err := encodeWithModel(test.text, model)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("encodeWithModel(%q) = %v, want %v", test.text, got, test.want)
		}
	}
}

func TestDecodeWithModelExpandsLearnedTokens(t *testing.T) {
	model := tokenizerModel{
		Version:       1,
		BaseVocabSize: 256,
		Merges:        [][]int{{97, 110}, {98, 256}, {256, 97}},
	}
	tests := []struct {
		name string
		ids  []int
		want []byte
	}{
		{name: "hierarchical tokens", ids: []int{257, 258, 33}, want: []byte("banana!")},
		{name: "base and learned tokens", ids: []int{257, 97, 32, 258}, want: []byte("bana ana")},
		{name: "empty", ids: []int{}, want: []byte{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := decodeWithModel(test.ids, model)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, test.want) {
				t.Fatalf("decodeWithModel(%v) = %v, want %v", test.ids, got, test.want)
			}
		})
	}
}

func TestDecodeWithModelPreservesRawBytes(t *testing.T) {
	model := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{169, 195}}}
	got, err := decodeWithModel([]int{256}, model)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{169, 195}
	if !bytes.Equal(got, want) {
		t.Fatalf("decodeWithModel([256]) = %v, want %v", got, want)
	}
}

func TestParseModelRejectsInvalidTokenGraphs(t *testing.T) {
	tests := []struct {
		name  string
		model string
	}{
		{name: "malformed JSON", model: "{not json"},
		{name: "missing field", model: `{"version":1,"base_vocab_size":256}`},
		{name: "wrong version", model: `{"version":2,"base_vocab_size":256,"merges":[]}`},
		{name: "wrong base size", model: `{"version":1,"base_vocab_size":128,"merges":[]}`},
		{name: "wrong field type", model: `{"version":"1","base_vocab_size":256,"merges":[]}`},
		{name: "wrong merge width", model: `{"version":1,"base_vocab_size":256,"merges":[[97]]}`},
		{name: "non-integer operand", model: `{"version":1,"base_vocab_size":256,"merges":[[97.5,98]]}`},
		{name: "negative operand", model: `{"version":1,"base_vocab_size":256,"merges":[[-1,97]]}`},
		{name: "self reference", model: `{"version":1,"base_vocab_size":256,"merges":[[256,97]]}`},
		{name: "future reference", model: `{"version":1,"base_vocab_size":256,"merges":[[257,97],[97,98]]}`},
		{name: "duplicate pair", model: `{"version":1,"base_vocab_size":256,"merges":[[97,98],[97,98]]}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseModel([]byte(test.model)); codeForError(err) != errInvalidModel {
				t.Fatalf("parseModel() error = %v, code = %q; want %q", err, codeForError(err), errInvalidModel)
			}
		})
	}
}

func TestParseModelAcceptsUnknownFields(t *testing.T) {
	encoded := []byte(`{"version":1,"base_vocab_size":256,"merges":[],"note":"ignored"}`)
	model, err := parseModel(encoded)
	if err != nil {
		t.Fatal(err)
	}
	want := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{}}
	if !reflect.DeepEqual(model, want) {
		t.Fatalf("parseModel() = %#v, want %#v", model, want)
	}
}

func TestDecodeRejectsUnknownTokens(t *testing.T) {
	if _, err := decodeBaseIDs([]int{256}); codeForError(err) != errUnknownToken {
		t.Fatalf("decodeBaseIDs() error = %v, code = %q; want %q", err, codeForError(err), errUnknownToken)
	}
	model := tokenizerModel{Version: 1, BaseVocabSize: 256, Merges: [][]int{{97, 110}}}
	if _, err := decodeWithModel([]int{257}, model); codeForError(err) != errUnknownToken {
		t.Fatalf("decodeWithModel() error = %v, code = %q; want %q", err, codeForError(err), errUnknownToken)
	}
}
