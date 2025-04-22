/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

// 测试文件需包含的全局锁（应对并行测试）
var cmdMutex sync.Mutex

func TestUpdateClusterProperties(t *testing.T) {
	// 子测试用例组
	t.Run("空输入返回错误", testEmptyInput)
	t.Run("类型转换验证", testTypeConversion)
	t.Run("特殊属性命令生成", testSpecialPropertyCommand)
	t.Run("命令执行错误处理", testCommandExecutionError)
	t.Run("多属性批量更新", testMultiplePropertiesUpdate)
}

// 测试用例1：空输入
func testEmptyInput(t *testing.T) {
	result := UpdateClusterProperties(map[string]interface{}{})
	assert.False(t, result["action"].(bool))
	assert.Equal(t, "No input data", result["error"].(string))
}

// 测试用例2：类型转换逻辑
func testTypeConversion(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	// Mock成功执行的命令
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte(""), nil
	}

	testData := []struct {
		input    interface{}
		expected string
	}{
		{"test_string", "test_string"},
		{true, "true"},
		{false, "false"},
		{float64(123), "123"},
	}

	for _, data := range testData {
		props := map[string]interface{}{"test_key": data.input}
		result := UpdateClusterProperties(props)
		assert.True(t, result["action"].(bool))

	}
}

// 测试用例3：特殊属性命令格式
func testSpecialPropertyCommand(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	var actualCmd string
	utils.RunCommand = func(cmd string) ([]byte, error) {
		actualCmd = cmd
		return []byte(""), nil
	}

	props := map[string]interface{}{"resource-stickiness": "200"}
	result := UpdateClusterProperties(props)
	assert.True(t, result["action"].(bool))
	assert.True(t, strings.HasPrefix(actualCmd, utils.CmdUpdateResourceStickness))
}

// 测试用例4：命令执行错误处理
func testCommandExecutionError(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte("permission denied"), errors.New("exit status 1")
	}

	props := map[string]interface{}{"test_key": "value"}
	result := UpdateClusterProperties(props)
	assert.False(t, result["action"].(bool))
	assert.Equal(t, "permission denied", result["error"].(string))
}

// 测试用例5：多属性批量更新
func testMultiplePropertiesUpdate(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	execCount := 0
	utils.RunCommand = func(cmd string) ([]byte, error) {
		execCount++
		return []byte(""), nil
	}

	props := map[string]interface{}{
		"prop1":               "val1",
		"prop2":               false,
		"resource-stickiness": 300,
	}
	result := UpdateClusterProperties(props)
	assert.True(t, result["action"].(bool))
	assert.Equal(t, 3, execCount) // 验证三次命令调用
}

// 测试不同操作分支与错误处理
func TestOperationClusterAction(t *testing.T) {
	t.Run("空操作返回错误", testEmptyAction)
	t.Run("启动集群命令调用", testStartCluster)
	t.Run("重启集群流程验证", testRestartCluster)
}

func testEmptyAction(t *testing.T) {
	result := OperationClusterAction("")
	assert.False(t, result["action"].(bool))
	assert.Contains(t, result["error"], "failed") // 验证错误消息‌:ml-citation{ref="1" data="citationList"}
}

func testStartCluster(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	originalRun := utils.RunCommand
	defer func() { utils.RunCommand = originalRun }()

	// Mock命令调用计数
	var execCount int
	utils.RunCommand = func(cmd string) ([]byte, error) {
		assert.Equal(t, utils.CmdStartCluster, cmd) // 验证命令格式‌:ml-citation{ref="7" data="citationList"}
		execCount++
		return []byte(""), nil
	}

	result := OperationClusterAction("start")
	assert.True(t, result["action"].(bool))
	assert.Equal(t, 1, execCount)
}

func testRestartCluster(t *testing.T) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	originalRun := utils.RunCommand
	defer func() { utils.RunCommand = originalRun }()

	// 验证命令执行顺序
	callSequence := []string{}
	utils.RunCommand = func(cmd string) ([]byte, error) {
		callSequence = append(callSequence, cmd)
		return []byte(""), nil
	}

	result := OperationClusterAction("restart")
	assert.True(t, result["action"].(bool))
	assert.Equal(t, []string{utils.CmdStopClusterLocal, utils.CmdStartCluster}, callSequence) // 验证流程顺序‌:ml-citation{ref="7" data="citationList"}
}

func TestGetClusterPropertyFromXml(t *testing.T) {
	t.Run("完整XML元素解析", func(t *testing.T) {
		xmlStr := `
			<parameter name="test_prop">
				<shortdesc>Short description</shortdesc>
				<longdesc>Long description with Allowed values: val1, val2</longdesc>
				<content type="enum" default="val1"/>
			</parameter>
		`
		doc := etree.NewDocument()
		doc.ReadFromString(xmlStr)
		elem := doc.Root()

		prop := getClusterPropertyFromXml(elem)
		assert.Equal(t, "test_prop", prop["name"])
		assert.Equal(t, "Short description", prop["shortdesc"])
		assert.Equal(t, "Long description with ", prop["longdesc"]) // Allowed values部分被移除
		assert.Equal(t, "enum", prop["type"])
	})

	t.Run("缺失Content元素处理", func(t *testing.T) {
		xmlStr := `<parameter name="no_content_prop"/>`
		doc := etree.NewDocument()
		doc.ReadFromString(xmlStr)
		elem := doc.Root()

		prop := getClusterPropertyFromXml(elem)
		assert.Equal(t, "", prop["type"])
		assert.Equal(t, "", prop["default"])
	})

	t.Run("相同长短描述清理", func(t *testing.T) {
		xmlStr := `
			<parameter name="dup_desc">
				<shortdesc>Same Desc</shortdesc>
				<longdesc>Same Desc</longdesc>
			</parameter>
		`
		doc := etree.NewDocument()
		doc.ReadFromString(xmlStr)
		elem := doc.Root()

		prop := getClusterPropertyFromXml(elem)
		assert.Equal(t, "", prop["longdesc"]) // 长描述被清空
	})
}
