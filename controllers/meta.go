/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: 赵忠章 <jinzi120021@sina.com>
 * Date: Fri Jan 22 11:18:35 2021 +0800
 */
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
)

type MetaController struct {
	web.Controller
}

func (mc *MetaController) Get() {
	rscClass := mc.Ctx.Input.Param(":rsc_class")
	rscType := mc.Ctx.Input.Param(":rsc_type")
	rscProvider := mc.Ctx.Input.Param(":rsc_provider")
	mc.Data["json"] = models.GetResourceMetas(rscClass, rscType, rscProvider)
	mc.ServeJSON()
}

type MetasController struct {
	web.Controller
}

func (mc *MetasController) Get() {
	mc.Data["json"] = models.GetAllResourceMetas()
	mc.ServeJSON()
}
