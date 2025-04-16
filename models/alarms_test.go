/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Mon Sep 2 15:58:49 2024 +0800
 */
package models

import (
	"fmt"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
	"github.com/stretchr/testify/assert"
)

func TestAlarmsGet_Success(t *testing.T) {
	// 备份并替换 RunCommand
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 模拟返回有效 XML
	utils.RunCommand = func(cmd string) ([]byte, error) {
		xml := `
        <configuration>
            <alerts>
                <alert>
                    <instance_attributes>
                        <nvpair name="email_sender" value="test@example.com"/>
                        <nvpair name="email_server" value="smtp.example.com"/>
                        <nvpair name="password" value="encrypted_pwd"/>
                        <nvpair name="port" value="587"/>
                        <nvpair name="switCh" value="on"/>
                    </instance_attributes>
                    <recipient value="user1@example.com"/>
                    <recipient value="user2@example.com"/>
                </alert>
            </alerts>
        </configuration>
        `
		return []byte(xml), nil
	}

	response := AlarmsGet()

	// 验证结果
	assert.True(t, response.Action)
	assert.True(t, response.Data.Flag) // switCh=on 应映射为 true
	assert.Equal(t, "test@example.com", response.Data.Sender)
	assert.Equal(t, "smtp.example.com", response.Data.Smtp)
	assert.Equal(t, 587.0, response.Data.Port)
	assert.Equal(t, []string{"user1@example.com", "user2@example.com"}, response.Data.Receiver)
}

func TestAlarmsGet_CommandError(t *testing.T) {
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 模拟命令返回错误
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return nil, fmt.Errorf("command failed")
	}

	response := AlarmsGet()

	// 即使出错，Action 仍为 true（根据代码逻辑）
	assert.True(t, response.Action)
	assert.Empty(t, response.Data.Sender) // 数据应为默认零值
}

func TestAlarmsGet_PasswordDecryptError(t *testing.T) {
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 模拟 XML 返回加密密码，但解密失败
	utils.RunCommand = func(cmd string) ([]byte, error) {
		if cmd == utils.CmdCibQueryConfig {
			xml := `
            <configuration>
                <alerts>
                    <alert>
                        <instance_attributes>
                            <nvpair name="password" value="encrypted_pwd"/>
                        </instance_attributes>
                    </alert>
                </alerts>
            </configuration>
            `
			return []byte(xml), nil
		} else if cmd == "/usr/bin/pwd_decode encrypted_pwd" {
			return []byte("the parameter is less\n"), nil // 解密失败
		}
		return nil, nil
	}

	response := AlarmsGet()

	assert.True(t, response.Action)
	assert.Equal(t, "", response.Data.Password) // 解密失败时 password 应为空
}

func TestAlarmsSet_Success(t *testing.T) {
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 记录实际执行的命令
	var executedCommands []string
	utils.RunCommand = func(cmd string) ([]byte, error) {
		executedCommands = append(executedCommands, cmd)
		return []byte("success"), nil
	}

	data := AlarmData{
		Flag:     true,
		Sender:   "admin@example.com",
		Smtp:     "smtp.example.com",
		Password: "secret",
		Port:     465,
		Receiver: []string{"user1@example.com", "user2@example.com"},
	}

	result := AlarmsSet(data)

	// 验证返回结果
	assert.True(t, result["action"].(bool))
}

func TestAlarmsSet_CommandFailure(t *testing.T) {
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	utils.RunCommand = func(cmd string) ([]byte, error) {
		return nil, fmt.Errorf("permission denied")
	}

	data := AlarmData{Sender: "test@example.com"}
	result := AlarmsSet(data)

	assert.False(t, result["action"].(bool))
	assert.Equal(t, gettext.Gettext("Set alarm failed"), result["error"])
}

func TestIsDataEmpty(t *testing.T) {
	// 空数据
	emptyData := AlarmData{}
	assert.True(t, isDataEmpty(emptyData))

	// 非空数据
	nonEmptyData := AlarmData{Sender: "test@example.com"}
	assert.False(t, isDataEmpty(nonEmptyData))
}
