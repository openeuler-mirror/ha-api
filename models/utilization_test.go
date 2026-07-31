/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

// 测试GetOneTypeUtilization基础解析功能
func TestGetOneTypeUtilization_Normal(t *testing.T) {
	// Mock正常命令输出
	mockOutput := `Utilization:
    node1: cpu=4 ram=16
    node2: cpu=8 ram=32 special=value-with=equal
    `

	utils.RunCommand = func(cmd string) ([]byte, error) {
		expectedCmd := "pcs node utilization"
		assert.Equal(t, expectedCmd, cmd)
		return []byte(mockOutput), nil
	}

	result := GetOneTypeUtilization("node")

	// 验证解析结果
	assert.Len(t, result, 2, "应解析出2个节点")

}

// 测试GetUtilization聚合功能
func TestGetUtilization(t *testing.T) {
	// 设置双mock
	utils.RunCommand = func(cmd string) ([]byte, error) {
		switch {
		case strings.Contains(cmd, "node"):
			return []byte("Utilization:\nnodeA: cpu=2"), nil
		case strings.Contains(cmd, "resource"):
			return []byte("Utilization:\nres1: priority=1"), nil
		}
		return nil, nil
	}

	result := GetUtilization()

	assert.True(t, result.Action)
}

// 测试设置利用率参数构建
func TestSetUtilization_ParameterBuild(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]string
		expected string
	}{
		{
			"NormalCase",
			map[string]string{
				"type": "node",
				"name": "node1",
				"cpu":  "4",
			},
			"pcs node utilization node1 cpu=4 ",
		},
		{
			"SpecialChars",
			map[string]string{
				"type": "resource",
				"name": "res1",
				"key":  "value=with=equals",
			},
			"pcs resource utilization res1 key=value=with=equals ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actualCmd string
			utils.RunCommand = func(cmd string) ([]byte, error) {
				actualCmd = cmd
				return []byte("success"), nil
			}

			data, _ := json.Marshal(tc.input)
			SetUtilization(data)

			assert.Equal(t, tc.expected, actualCmd)
		})
	}
}

// 测试设置利用率异常场景
func TestSetUtilization_ErrorCases(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		resp := SetUtilization(nil)
		assert.False(t, resp.Action)
		assert.Contains(t, resp.Error, "No input data")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		resp := SetUtilization([]byte("{invalid}"))
		assert.False(t, resp.Action)
		assert.Contains(t, resp.Error, "Cannot convert data")
	})

	t.Run("CommandFailure", func(t *testing.T) {
		utils.RunCommand = func(cmd string) ([]byte, error) {
			return []byte("Error: Not found"), errors.New("command failed")
		}

		data, _ := json.Marshal(map[string]string{
			"type": "node",
			"name": "invalid-node",
		})
		resp := SetUtilization(data)

		assert.False(t, resp.Action)
	})
}

// 测试删除利用率功能
func TestDelUtilization(t *testing.T) {
	t.Run("NormalDeletion", func(t *testing.T) {
		var actualCmd string
		utils.RunCommand = func(cmd string) ([]byte, error) {
			actualCmd = cmd
			return []byte(""), nil
		}

		data, _ := json.Marshal(map[string]string{
			"type": "resource",
			"name": "res1",
			"cpu":  "", // 删除cpu属性
		})
		resp := DelUtilization(data)

		assert.True(t, resp.Action)
		assert.Equal(t, "pcs resource utilization res1 cpu=", actualCmd)
	})

	t.Run("MultiAttribute", func(t *testing.T) {
		data, _ := json.Marshal(map[string]string{
			"type":  "node",
			"name":  "node1",
			"attr1": "",
			"attr2": "",
		})

		var actualCmd string
		utils.RunCommand = func(cmd string) ([]byte, error) {
			actualCmd = cmd
			return []byte(""), nil
		}

		DelUtilization(data)
		assert.Contains(t, actualCmd, "attr1=")
		assert.Contains(t, actualCmd, "attr2=")
	})
}
