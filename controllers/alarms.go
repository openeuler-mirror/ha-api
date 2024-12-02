/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: liqiuyu <liqiuyu@kylinos.cn>
 * Date: Mon Jan 18 11:44:18 2021 +0800
 */
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

	"encoding/json"

	"gitee.com/openeuler/ha-api/models"
)

type AlarmConfig struct {
	web.Controller
}

func (ac *AlarmConfig) Get() {
	ac.Data["json"] = models.AlarmsGet()
	ac.ServeJSON()
}

func (ac *AlarmConfig) Post() {
	var result map[string]interface{}

	reqData := make(map[string]string)
	if err := json.Unmarshal(ac.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.AlarmsSet(reqData)
	}
	ac.Data["json"] = &result
	ac.ServeJSON()
}
