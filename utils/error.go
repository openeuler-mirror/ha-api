/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

import (
	"fmt"
	"log/slog"

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
	slog.Error(msg)
	LogTrace(err)
}

func LogTrace(err error) {
	slog.Error(fmt.Sprintf("original error: %T %v", errors.Cause(err), errors.Cause(err)))
	slog.Error(fmt.Sprintf("stack trace: %+v", err))
}
