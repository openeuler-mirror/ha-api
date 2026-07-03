/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Wed Aug 21 09:42:19 2024 +0800
 */
package models

import (
	"fmt"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateLog_Success(t *testing.T) {
	// 备份原函数，测试结束后恢复
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// Mock 成功场景
	utils.RunCommand = func(cmd string) ([]byte, error) {
		// 验证传入的命令是否正确
		if cmd != utils.CmdGenLog {
			t.Fatalf("Expected command: %s, got: %s", utils.CmdGenLog, cmd)
		}
		// 返回模拟的成功输出（包含换行符）
		return []byte("kylinha-log-test.tar.gz\n"), nil
	}

	// 执行被测函数
	result, err := GenerateLog()

	// 验证结果
	assert.Nil(t, err)
	assert.Contains(t, result, ".tar.gz")
}
