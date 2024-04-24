/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-19 17:37:49
 * Description: 登录控制器
 ******************************************************************************/
package controllers

import (
	"encoding/json"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/utils"
)

func init() {
	web.InsertFilter("/*", web.BeforeRouter, loginFilter)
}

const (
	apiPreffix     = "/api"
	versionPreffix = "/v1"
)

var loginFilter = func(ctx *context.Context) {
	// implement login check
	username := ctx.Input.Session("username")
	if username == nil {
		if !strings.Contains(ctx.Request.RequestURI, "/login") && !strings.Contains(ctx.Request.RequestURI, "/remote") {
			if strings.HasPrefix(ctx.Request.RequestURI, apiPreffix) {
				ctx.Redirect(403, "session timeout")
			}
		}
	}
}

type LoginController struct {
	web.Controller
}

func (lc *LoginController) Post() {
	logs.Debug("handle post request in LoginController.")
	result := struct {
		Action bool   `json:"action"`
		Error  string `json:"error,omitempty"`
		Info   string `json:"info,omitempty"`
	}{true, "", ""}

	d := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(lc.Ctx.Input.RequestBody, &d); err != nil {
		result.Action = false
		result.Error = err.Error()
		logs.Error("Login failed:", err)
		goto ret
	}
	if d.Username == "" || d.Password == "" {
		result.Action = false
		result.Error = gettext.Gettext("username or password empty")
		logs.Error("Login failed: username or password is empty")
		goto ret
	}
	if d.Username != "hacluster" {
		result.Action = false
		result.Error = gettext.Gettext("username is not allowed to login")
		logs.Error("Login failed: username is not allowed to login")
		goto ret
	}

	// check password
	if utils.CheckAuth(d.Username, d.Password) {
		// update session
		lc.SetSession("username", d.Username)
		lc.SetSession("password", d.Password)

		result.Action = true
		result.Info = gettext.Gettext("Login success")
	} else {
		result.Action = false
		result.Error = gettext.Gettext("Username or password error")
		logs.Error("Login failed: username or password error")
	}

ret:
	lc.Data["json"] = &result
	lc.ServeJSON()
}

type LogoutController struct {
	web.Controller
}

func (lc *LogoutController) Post() {
	// delete session
	lc.DelSession("username")
	lc.DelSession("password")

	result := struct {
		Action bool   `json:"action"`
		Error  string `json:"error,omitempty"`
		Info   string `json:"info,omitempty"`
	}{true, "", ""}
	result.Info = gettext.Gettext("Logout success")
	logs.Info("Logout success")
	lc.Data["json"] = &result
	lc.ServeJSON()
}

type PasswordChangeController struct {
	web.Controller
}

func (lc *PasswordChangeController) Post() {
	logs.Debug("handle post request in PasswordChangeController.")
	var cmd string = ""
	result := struct {
		Action bool   `json:"action"`
		Error  string `json:"error,omitempty"`
		Info   string `json:"info,omitempty"`
	}{false, "", ""}

	data := struct {
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(lc.Ctx.Input.RequestBody, &data); err != nil {
		result.Action = false
		result.Error = err.Error()
		goto ret
	}

	if data.Password == "" {
		result.Action = false
		result.Error = gettext.Gettext("The new password is empty")
		goto ret
	}
	cmd = "echo 'hacluster:" + data.Password + "'|chpasswd >/dev/null 2>&1"
	if _, err := utils.RunCommand(cmd); err != nil {
		logs.Error("run command error: ", err)
		result.Action = false
		result.Error = err.Error()
		goto ret
	} else {
		result.Action = true
		result.Info = gettext.Gettext("Change password success")
	}

ret:
	lc.Data["json"] = &result
	lc.ServeJSON()
}

// for test

type TestController struct {
	web.Controller
}

func (lc *TestController) Get() {

	// result := struct {
	// 	Action bool   `json:"action"`
	// 	Error  string `json:"info"`
	// }{true, "logout"}
	// lc.Data["json"] = &result
	lc.ServeJSON()
}
