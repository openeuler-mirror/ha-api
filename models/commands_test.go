/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Tue Aug 20 16:27:26 2024 +0800
 */
package models

import (
	"errors"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCommandsList(t *testing.T) {
	result := GetCommandsList()
	if result["action"] != true {
		t.Fatal("Get commands list failed")
	}
}

func TestRunBuiltinCommand_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		cmdID       int
		mockOutput  string
		mockError   error
		expectError string
	}{
		{"ValidCommand", 3, "config data", nil, ""},
		{"CommandError", 4, "", errors.New("timeout"), "timeout"},
		{"InvalidID", 5, "", nil, "invalid command index"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock 函数设置
			originalRunCommand := utils.RunCommand
			defer func() { utils.RunCommand = originalRunCommand }()

			utils.RunCommand = func(cmd string) ([]byte, error) {
				return []byte(tt.mockOutput), tt.mockError
			}

			// 执行测试
			output, err := RunBuiltinCommand(tt.cmdID)

			// 结果验证
			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockOutput, output)
			}
		})
	}
}
