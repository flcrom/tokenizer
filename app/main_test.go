/**
 * @file app/main_test.go
 * @description Verifies base-token round trips and adjacent-pair counting independently of the process harness.
 * @depends_on app/main.go
 * @depended_on_by none
 */

package main

import (
	"bytes"
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
			if got := decodeBaseIDs(test.ids); !bytes.Equal(got, test.want) {
				t.Fatalf("decodeBaseIDs(%v) = %v, want %v", test.ids, got, test.want)
			}
		})
	}
}

func TestBaseTokenRoundTrip(t *testing.T) {
	fixtures := []string{"", "hello", "café", "🙂", "é", "line one\nline two"}
	for _, fixture := range fixtures {
		decoded := decodeBaseIDs(encodeBytes(fixture))
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
