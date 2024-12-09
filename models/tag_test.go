/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-12-05 15:28:13
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-12-09 11:28:15
 * @FilePath: /ha-api/models/tag_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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

func TestUpdateTag(t *testing.T) {
	testCases := []struct {
		inputTagName string
		inputData    []byte
		expected     map[string]interface{}
	}{
		{
			`tag1`,
			[]byte(`{"id":"tag1","tag_resource": ["res1"]}`),
			map[string]interface{}{"action": true, "info": "Update tag success"},
		},
	}

	for _, testCase := range testCases {
		result := UpdateTag(testCase.inputTagName, testCase.inputData)
		resultJson, err := json.Marshal(result)
		require.NoError(t, err, "Marshal not return an error")

		expectedJson, err := json.Marshal(testCase.expected)
		require.NoError(t, err, "Marshal expected map not return an error")

		assert.JSONEq(t, string(expectedJson), string(resultJson), "Update tag success")
	}
}
