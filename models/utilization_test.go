/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-11-22 14:30:59
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-11-25 17:08:25
 * @FilePath: /ha-api/models/utilization_test.go
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
 * Author: bixiaoyan
 * Date: 2024-11-22 14:25:51
 * LastEditTime: 2024-11-22 14:37:48
 * Description: 集群利用率功能单元测试
 ******************************************************************************/

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
