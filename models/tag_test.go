/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Dec 05 15:26:24 2024 +0800
 */

package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTag(t *testing.T) {
	result := GetTag()
	if result.Action != true {
		t.Fatal("Get tag failed")
	}
}

func TestSetTag(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected map[string]interface{}
	}{
		{
			[]byte(`{"id":"tag1","tag_resource": ["res1"]}`),
			map[string]interface{}{"action": true, "info": "Add tag success"},
		},
	}

	for _, testCase := range testCases {
		result := SetTag(testCase.input)
		resultJson, err := json.Marshal(result)
		require.NoError(t, err, "Marshal not return an error")

		expectedJson, err := json.Marshal(testCase.expected)
		require.NoError(t, err, "Marshal expected map not return an error")

		assert.JSONEq(t, string(expectedJson), string(resultJson), "Add tag success")
	}
}
