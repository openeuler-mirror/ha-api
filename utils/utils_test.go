package utils

import (
	"testing"
)

func TestKeys(t *testing.T) {
	m := map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 1,
	}
	keys := Keys[string, int](m)
	expectedKeys := []string{"key1", "key2", "key3"}
	assertSliceElements[string](t, keys, expectedKeys)
}

func TestValues(t *testing.T) {
	testMap := map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 1,
	}

	values := Values[string, int](testMap)

	expectedValues := []int{1, 2, 1}
	assertSliceElements[int](t, values, expectedValues)
}

func TestDifferenceSlice(t *testing.T) {
	mainSlice := []int{1, 2, 3, 4, 5}
	subtractSlice := []int{3, 4, 6}
	expectedSlice := []int{1, 2, 5}
	diffSlice := DifferenceSlice(mainSlice, subtractSlice)

	assertSliceElements[int](t, diffSlice, expectedSlice)
}

func assertSliceElements[K comparable](t *testing.T, actual, expected []K) {
	t.Helper()

	expectedCount := make(map[K]int)
	for _, v := range expected {
		expectedCount[v]++
	}

	for _, v := range actual {
		if count, ok := expectedCount[v]; !ok || count == 0 {
			t.Errorf("Slice contains unexpected value %v", v)
		}
		expectedCount[v]--
	}

	for _, count := range expectedCount {
		if count != 0 {
			t.Errorf("Slice is missing value")
		}
	}
}
