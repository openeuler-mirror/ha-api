package utils

import (
	"reflect"
	"sort"
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
	sort.Strings(keys)
	sort.Strings(expectedKeys)
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Keys() = %v, want %v", keys, expectedKeys)
	}

	// 创建一个map来记录期望值出现的次数
	expectedCount := make(map[string]int)
	for _, v := range expectedKeys {
		expectedCount[v]++
	}

	// 检查返回的slice中的每个值是否都在期望的slice中
	for _, v := range keys {
		if count, ok := expectedCount[v]; !ok || count == 0 {
			t.Errorf("Keys() contains unexpected value %v", v)
		}
		expectedCount[v]--
	}

	// 检查期望的slice中的每个值是否都被找到了
	for _, count := range expectedCount {
		if count != 0 {
			t.Errorf("Keys() is missing value")
		}
	}

}

func TestValues(t *testing.T) {
	testMap := map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 1,
	}

	values := Values2[string, int](testMap)

	expectedKeys := []int{1, 2, 1}

	// 创建一个map来记录期望值出现的次数
	expectedCount := make(map[int]int)
	for _, v := range expectedKeys {
		expectedCount[v]++
	}

	// 检查返回的slice中的每个值是否都在期望的slice中
	for _, v := range values {
		if count, ok := expectedCount[v]; !ok || count == 0 {
			t.Errorf("Values() contains unexpected value %v", v)
		}
		expectedCount[v]--
	}

	// 检查期望的slice中的每个值是否都被找到了
	for _, count := range expectedCount {
		if count != 0 {
			t.Errorf("Values() is missing value")
		}
	}
}
