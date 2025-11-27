package shim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandElem(t *testing.T) {

	// 示例：从字符串切片中随机选择一个元素
	strings := []string{"apple", "banana", "cherry", "date"}
	randomString := RandElem(strings)
	t.Logf("随机选择的字符串: %s", randomString)

	// 示例：从整数切片中随机选择一个元素
	integers := []int{1, 2, 3, 4, 5}
	randomInt := RandElem(integers)
	t.Logf("随机选择的字符串: %v", randomInt)
}

func containsSameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[T]int)
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		if counts[v] == 0 {
			return false
		}
		counts[v]--
	}
	return true
}

func TestShuffle(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{"empty slice", args[int]{input: []int{}}},
		{"single element", args[int]{input: []int{42}}},
		{"multiple elements", args[int]{input: []int{1, 2, 3, 4, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Shuffle(tt.args.input)
			t.Logf("%v", got)
			assert.Equal(t, len(tt.args.input), len(got), "length should be equal")
			assert.True(t, containsSameElements(tt.args.input, got), "shuffled slice should contain same elements")
		})
	}
}
