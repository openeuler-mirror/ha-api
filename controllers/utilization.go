/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Tue Mar 12 09:17:37 2024 +0800
 */
package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
)

type UtilizationController struct {
	web.Controller
}

func (ul *UtilizationController) Get() {

	ul.Data["json"] = models.GetUtilization()
	ul.ServeJSON() 
}

func (ul *UtilizationController) Post() {
	logs.Debug("handle utilization POST request")

	jsonStr := ul.Ctx.Input.RequestBody
	ul.Data["json"] = models.SetUtilization(jsonStr)
	ul.ServeJSON()
}

func (ul *UtilizationController) Delete() {
	logs.Debug("handle utilization DELETE request")

	jsonStr := ul.Ctx.Input.RequestBody
	ul.Data["json"] = models.DelUtilization(jsonStr)
	ul.ServeJSON()
}


 