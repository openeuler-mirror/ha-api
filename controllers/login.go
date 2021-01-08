package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	"openkylin.com/ha-api/utils"
)

func init() {
	web.InsertFilter("/*", web.BeforeRouter, loginFilter)
}

var loginFilter = func(ctx *context.Context) {
	// TODO: implement login check
}

type LoginController struct {
	web.Controller
}

func (lc *LoginController) Post() {
	logs.Debug("handle psot request in LoginController.")
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
		// TODO: update session
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
	// TODO: delete session

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
