/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Fri Nov 22 14:35:35 2024 +0800
 */

package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUtilization(t *testing.T) {
	result := GetUtilization()
	if result.Action != true {
		t.Fatal("Get utilization failed")
	}
}

func TestSetUtilization(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected map[string]interface{}
	}{
		{
			[]byte(`{"type":"node","name": "ha1","cpu":"1"}`),
			map[string]interface{}{"action": true, "info": "Utilization set success"},
		},
		{
			[]byte(`{"type":"resource","name": "d","cpu":"1"}`),
			map[string]interface{}{"action": true, "info": "Utilization set success"},
		},
	}

	for _, testCase := range testCases {
		result := SetUtilization(testCase.input)
		resultJson, err := json.Marshal(result)
		require.NoError(t, err, "Marshal not return an error")

		expectedJson, err := json.Marshal(testCase.expected)
		require.NoError(t, err, "Marshal expected map not return an error")

		assert.JSONEq(t, string(expectedJson), string(resultJson), "Set utilization case success")
	}
}

func TestDelUtilization(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected map[string]interface{}
	}{
		{
			[]byte(`{"type":"node","name": "ha1","cpu":"1"}`),
			map[string]interface{}{"action": true, "info": "Delete Utilization success"},
		},
		{
			[]byte(`{"type":"resource","name": "d","cpu":"1"}`),
			map[string]interface{}{"action": true, "info": "Delete Utilization success"},
		},
	}

	for _, testCase := range testCases {
		result := DelUtilization(testCase.input)
		resultJson, err := json.Marshal(result)
		require.NoError(t, err, "Marshal not return an error")

		expectedJson, err := json.Marshal(testCase.expected)
		require.NoError(t, err, "Marshal expected map not return an error")

		assert.JSONEq(t, string(expectedJson), string(resultJson), "Delete Utilization success")
	}
}
