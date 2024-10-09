/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-03-21 17:02:57
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-10-09 17:33:24
 * @FilePath: /ha-api/models/cluster_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2024-08-20 10:30:17
 * Description: 集群测试用例
 ******************************************************************************/
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
