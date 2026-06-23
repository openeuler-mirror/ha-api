/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/msteinert/pam"
)

// auth.go
type AuthChecker interface {
	CheckAuth(username, password string) bool
}

// 真实实现
type PamAuthChecker struct{}

func (p *PamAuthChecker) CheckAuth(username, password string) bool {
	t, err := pam.StartFunc("login", username, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			return "", nil
		case pam.ErrorMsg:
			return "", nil
		case pam.TextInfo:
			return "", nil
		}
		return "", errors.New("unrecognized message style")
	})
	defer func() {
		runtime.GC()
	}()

	if err != nil {
		slog.Error(fmt.Sprintf("pam start error: %s", err))
		return false
	}

	err = t.Authenticate(0)
	if err != nil {
		slog.Error(fmt.Sprintf("pam auth error: %s", err))
		return false
	}

	return true
}
