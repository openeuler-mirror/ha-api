/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 11:28:40
 * Description: PAM管理
 ******************************************************************************/
package utils

import (
	"errors"
	"runtime"

	"github.com/beego/beego/v2/core/logs"
	"github.com/msteinert/pam"
)

func CheckAuth(username, password string) bool {
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
		return "", errors.New("Unrecognized message style")
	})
	defer func() {
		runtime.GC()
	}()

	if err != nil {
		logs.Error("pam start error:", err)
		return false
	}

	err = t.Authenticate(0)
	if err != nil {
		logs.Error("pam auth error:", err)
		return false
	}

	return true
}
