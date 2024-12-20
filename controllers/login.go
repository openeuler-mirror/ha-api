/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package controllers

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
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
	resp := new(utils.Response)
	requestInput := new(validations.UserS)

	err := validations.UnmarshalAndValidation(lc.Ctx.Input.RequestBody, requestInput)
	if err != nil {
		resp.Action = false
		resp.Error = err.Error()
		goto ret
	}

	// check password
	if utils.CheckAuth(requestInput.UserName, requestInput.Password) {
		// update session
		lc.SetSession("username", requestInput.UserName)
		lc.SetSession("password", requestInput.Password)

		resp.Action = true
		resp.Info = gettext.Gettext("Login success")
	} else {
		resp.Action = false
		resp.Error = gettext.Gettext("Username or password error")
		logs.Error("Login failed: username or password error")
	}

ret:
	lc.Data["json"] = &resp
	lc.ServeJSON()
}

type LogoutController struct {
	web.Controller
}

func (lc *LogoutController) Post() {
	// delete session
	lc.DelSession("username")
	lc.DelSession("password")
	result := new(utils.Response)

	result.Info = gettext.Gettext("Logout success")
	result.Action = true
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

	result := new(utils.Response)
	requestInput := new(validations.PasswordS)

	err := validations.UnmarshalAndValidation(lc.Ctx.Input.RequestBody, requestInput)
	if err != nil {
		result.Action = false
		result.Error = err.Error()
		goto ret
	}

	cmd = "echo 'hacluster:" + requestInput.Password + "'|chpasswd >/dev/null 2>&1"
	if _, err := utils.RunCommand(cmd); err != nil {
		logs.Error("run command error: ", err)
		result.Action = false
		result.Error = gettext.Gettext("Change password failed")
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
