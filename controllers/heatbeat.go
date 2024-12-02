/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Thu Jan 14 13:33:38 2021 +0800
 */
package controllers

import (
	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"
)

type HeartBeatController struct {
	web.Controller
}

func (hbc *HeartBeatController) Get() {
	result := map[string]interface{}{}

	data, err := models.GetHeartBeatConfig()
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = data
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

func (hbc *HeartBeatController) Post() {
	result := map[string]interface{}{}

	data := hbc.Ctx.Input.RequestBody
	if err := models.EditHeartbeatInfo(data); err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = gettext.Gettext("Change cluster success")
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

type HeartBeatStatusController struct {
	web.Controller
}

func (hbsc *HeartBeatStatusController) Get() {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = 0

	hbsc.Data["json"] = &result
	hbsc.ServeJSON()
}

type DiskHeartBeatController struct {
	web.Controller
}

func (dhbc *DiskHeartBeatController) Get() {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = []string{}

	dhbc.Data["json"] = &result
	dhbc.ServeJSON()
}
