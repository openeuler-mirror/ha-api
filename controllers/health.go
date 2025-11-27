/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Fri Jul 04 15:54:28 2025 +0800
 */

package controllers

import (
	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/server/web"
)

type HealthConfig struct {
	web.Controller
}

func (hc *HealthConfig) Get() {
	hc.Data["json"] = models.HealthGet()
	hc.ServeJSON()
}

func (hc *HealthConfig) Post() {
	jsonStr := hc.Ctx.Input.RequestBody
	hc.Data["json"] = models.HealthSet(jsonStr)
	hc.ServeJSON()
}

func (hc *HealthConfig) Put() {
	var healthTestData models.HealthTestData
	healthTestData = models.HealthTest()
	hc.Data["json"] = &healthTestData
	hc.ServeJSON()
}
