/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClusterPropertiesInfo(t *testing.T) {
	result := GetClusterPropertiesInfo()
	if result["action"] != true {
		t.Fatal("Get cluster properties failed")
	}
}

func TestUpdateClusterProperties(t *testing.T) {
	clusterPropJson := map[string]interface{}{}
	clusterPropJson["no-quorum-policy"] = "ignore"
	res := UpdateClusterProperties(clusterPropJson)
	if res["action"] != true {
		t.Fatal("Update cluster properties failed")
	}
}

func TestGetClusterStatus(t *testing.T) {
	result := GetClusterStatus()
	if result == -1 {
		t.Fatal("Get cluster status failed")
	}
}

func TestGetClusterProperties(t *testing.T) {
	result, _ := getClusterProperties()
	if result == nil {
		t.Fatal("Get cluster properties failed")
	}
}
func TestGetResourceStickiness(t *testing.T) {
	result := getResourceStickiness()
	if result == 0 {
		t.Fatal("Get resource stickiness failed")
	}
}

// ["start","stop","restart",""]
func TestOperationClusterAction(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]interface{}
	}{
		{"start", map[string]interface{}{"action": true, "info": "Action on node success"}},
		{"stop", map[string]interface{}{"action": true, "info": "Action on node success"}},
		{"restart", map[string]interface{}{"action": true, "info": "Action on node success"}},
	}
	for _, testCase := range testCases {
		result := OperationClusterAction(testCase.input)
		resultJson, err := json.Marshal(result)
		require.NoError(t, err, "Marshal not return an error")

		expectedJson, err := json.Marshal(testCase.expected)
		require.NoError(t, err, "Marshal expected map not return an error")

		assert.JSONEq(t, string(expectedJson), string(resultJson), "case success")
	}
}
