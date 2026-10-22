/**
 * @file app/main_test.go
 * @description Verifies Stage 1 byte-to-token conversion independently of the process harness.
 * @depends_on app/main.go
 * @depended_on_by none
 */

package main

import (
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
