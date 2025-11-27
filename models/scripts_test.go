/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock HTTP Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) SendRequest(url, method string, body []byte) (*http.Response, error) {
	args := m.Called(url, method, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Mock Command Runner
func mockRunCommand(cmd string) ([]byte, error) {
	switch cmd {
	case utils.CmdNodeStatus:
		return []byte("Online  node1\nOnline  node2\nOffline node3"), nil
	case utils.CmdHostName:
		return []byte("node1"), nil
	}
	return nil, nil
}


func TestGenerateScript_GetNodeStatusFailed(t *testing.T) {
	// Mock command failure
	utils.RunCommand = func(cmd string) ([]byte, error) {
		if cmd == utils.CmdNodeStatus {
			return nil, errors.New("connection error")
		}
		return nil, nil
	}

	result := GenerateScript(map[string]string{})

	assert.False(t, result.Action)
	assert.Contains(t, result.Error, "get node status failed")
}

func TestGenerateScript_LocalGenerationFailed(t *testing.T) {
	utils.RunCommand = mockRunCommand

	// Mock local generation failure
	GenerateLocalScript = func(data map[string]string) error {
		return errors.New("file write error")
	}

	result := GenerateScript(map[string]string{})

	assert.True(t, result.Action)
	// assert.Equal(t, "Generate script failed", result.Data["node1"])
}

func TestGenerateScript_RemoteServerError(t *testing.T) {
	utils.RunCommand = mockRunCommand
	GenerateLocalScript = func(data map[string]string) error { return nil }

	mockHTTP := new(MockHTTPClient)
	utils.SendRequest = func(url string, method string, data interface{}) (*http.Response, error) {
		// 类型断言转换 interface{} -> []byte
		body, ok := data.([]byte)
		if !ok {
			return nil, errors.New("invalid data type")
		}

		return mockHTTP.SendRequest(url, method, body)
	}

	// Mock 500 response
	mockResp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	mockHTTP.On("SendRequest", mock.Anything, mock.Anything, mock.Anything).
		Return(mockResp, nil)

	result := GenerateScript(map[string]string{})

	assert.True(t, result.Action)
	// assert.Equal(t, "Generate script failed", result.Data["node2"])
}

func TestGenerateScript_InvalidRemoteResponse(t *testing.T) {
	utils.RunCommand = mockRunCommand
	GenerateLocalScript = func(data map[string]string) error { return nil }

	mockHTTP := new(MockHTTPClient)
	utils.SendRequest = func(url string, method string, data interface{}) (*http.Response, error) {
		// 类型断言转换 interface{} -> []byte
		body, ok := data.([]byte)
		if !ok {
			return nil, errors.New("invalid data type")
		}

		return mockHTTP.SendRequest(url, method, body)
	}

	// Return invalid JSON
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("{invalid}")),
	}
	mockHTTP.On("SendRequest", mock.Anything, mock.Anything, mock.Anything).
		Return(mockResp, nil)

	result := GenerateScript(map[string]string{})

	assert.True(t, result.Action)
	// assert.Equal(t, "Generate script failed", result.Data["node2"])
}


// Mock 工具函数
func mockRunCommand2param(output string, err error) (func(), error) {
	origRunCommand := utils.RunCommand
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte(output), err
	}
	return func() { utils.RunCommand = origRunCommand }, nil
}

func TestIsScriptExist(t *testing.T) {
	// 测试用例
	tests := []struct {
		name            string
		mockOutput      string
		mockError       error
		scriptName      string
		wantAction      bool
		wantError       string
		wantInfo        string
		wantLogContains string
	}{
		{
			name:            "脚本存在",
			mockOutput:      "script1\ntarget_script\nscript3",
			scriptName:      "target_script",
			wantAction:      false,
			wantError:       gettext.Gettext("The script already exists in the pacemaker directory"),
			wantLogContains: "already exists",
		},
		{
			name:       "脚本不存在",
			mockOutput: "script1\nscript2",
			scriptName: "target_script",
			wantAction: true,
			wantInfo:   gettext.Gettext("Script not exists"),
		},
		{
			name:       "命令执行失败",
			mockOutput: "",
			mockError:  errors.New("command error"),
			scriptName: "any_script",
			wantAction: false,
			wantError:  gettext.Gettext("Execute query script command failed"),
		},
		{
			name:       "空输出处理",
			mockOutput: "",
			scriptName: "test_script",
			wantAction: true,
			wantInfo:   gettext.Gettext("Script not exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock 命令执行
			cleanup, _ := mockRunCommand2param(tt.mockOutput, tt.mockError)
			defer cleanup()

			// 执行测试
			resp := IsScriptExist(tt.scriptName)

			// 验证响应
			assert.Equal(t, tt.wantAction, resp.Action)
			if tt.wantError != "" {
				assert.Contains(t, resp.Error, tt.wantError)
			} else {
				assert.Empty(t, resp.Error)
			}
			if tt.wantInfo != "" {
				assert.Contains(t, resp.Info, tt.wantInfo)
			}
		})
	}
}
