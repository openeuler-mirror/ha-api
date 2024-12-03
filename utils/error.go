/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Wed Mar 13 11:02:27 2024 +0800
 */

package utils

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/pkg/errors"
)

type ErrorInfo struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
}

func HandleCmdError(errMsg string, action bool) GeneralResponse {
	return GeneralResponse{
		Action: action,
		Error:  errMsg,
	}
}

func HandleXmlError(errMsg string, action bool) GeneralResponse {
	return GeneralResponse{
		Action: action,
		Error:  errMsg,
	}
}

// handleJsonError 函数处理JSON错误。
func HandleJsonError(errMsg string, action bool) ErrorInfo {
	result := ErrorInfo{
		Action: action,
		Error:  errMsg,
	}
	return result
}

func LogTraceWithMsg(err error, msg string) {
	logs.Error(msg)
	LogTrace(err)
}

func LogTrace(err error) {
	logs.Error("original error: %T %v", errors.Cause(err), errors.Cause(err))
	logs.Error("stack trace: %+v", err)
}
