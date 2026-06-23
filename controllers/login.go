/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */
package controllers

import (
	"encoding/pem"
	"log/slog"

	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/filters"
	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
)

func init() {
	web.InsertFilter("/*", web.BeforeRouter, filters.LoginFilter)
	web.InsertFilter("/api/v1/*", web.BeforeExec, filters.RemoteClusterRedirectFilter)
	web.InsertFilter("*", web.BeforeExec, filters.LogRequestControllerFilter)
	web.InsertFilter("/api/v1/*", web.FinishRouter, filters.LogResponseControllerFilter, web.WithReturnOnOutput(false))
}

type LoginController struct {
	web.Controller
}

func (lc *LoginController) Post() {
	slog.Debug("handle post request in LoginController.")
	resp := map[string]interface{}{}
	requestInput := new(validations.UserS)
	authChecker := new(utils.PamAuthChecker)
	err := validations.UnmarshalAndValidation(lc.Ctx.Input.RequestBody, requestInput)
	if err != nil {
		resp["action"] = false
		resp["error"] = err.Error()
		goto ret
	}
	resp = models.Login(*requestInput, authChecker)
	if resp["action"].(bool) {
		// update session
		lc.SetSession("username", requestInput.UserName)
		lc.SetSession("password", requestInput.Password)
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
	result := new(utils.GeneralResponse)

	result.Info = gettext.Gettext("Logout success")
	result.Action = true
	slog.Info("Logout success")
	lc.Data["json"] = &result
	lc.ServeJSON()
}

type PasswordChangeController struct {
	web.Controller
}

func (lc *PasswordChangeController) Post() {
	slog.Debug("handle post request in PasswordChangeController.")
	var result utils.GeneralResponse
	requestInput := new(validations.PasswordS)
	authChecker := new(utils.PamAuthChecker)
	err := validations.UnmarshalAndValidation(lc.Ctx.Input.RequestBody, requestInput)
	if err != nil {
		result.Action = false
		result.Error = err.Error()
		goto ret
	}
	result = models.PasswordChange(*requestInput, authChecker)
ret:
	lc.Data["json"] = &result
	lc.ServeJSON()
}

type PublicKey struct {
	Key string `json:"public_key"`
}

type PublicKeyResponse struct {
	Action bool      `json:"action"`
	Data   PublicKey `json:"data"`
}

type KeyController struct {
	web.Controller
}

func (kc *KeyController) Get() {
	var resp PublicKeyResponse
	var block *pem.Block
	publicKey, err := utils.ReadPublicKey(settings.RSA_PUBLIC_KEY)
	if err != nil {
		resp = PublicKeyResponse{
			Action: false,
			Data:   PublicKey{Key: ""},
		}
		goto ret
	}
	if block, _ = pem.Decode(publicKey); block == nil {
		resp = PublicKeyResponse{
			Action: false,
			Data:   PublicKey{Key: ""},
		}
		goto ret
	}
	resp.Action = true
	resp.Data.Key = string(publicKey)

ret:
	kc.Data["json"] = &resp
	kc.ServeJSON()
}
