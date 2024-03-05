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
		if !strings.Contains(ctx.Request.RequestURI, "/login") {
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
		Error  string `json:"error"`
	}{false, ""}

	d := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(lc.Ctx.Input.RequestBody, &d); err != nil {
		result.Action = false
		result.Error = "Login failed"
		goto ret
	}
	if d.Username == "" || d.Password == "" {
		result.Action = false
		result.Error = "Username or password error"
		goto ret
	}
	if d.Username != "hacluster" {
		result.Action = false
		result.Error = "Username is not allowed to login"
		goto ret
	}

	// check password
	if utils.CheckAuth(d.Username, d.Password) {
		// update session
		lc.SetSession("username", d.Username)
		lc.SetSession("password", d.Password)

		result.Action = true
		result.Error = ""
	} else {
		result.Action = false
		result.Error = "Username or password error"
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
		Error  string `json:"info"`
	}{true, "logout"}
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
