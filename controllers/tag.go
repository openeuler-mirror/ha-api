/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Fri Mar 22 17:26:24 2024 +0800
 */
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/core/logs"
)

type TagController struct {
	web.Controller
}

func (tc *TagController) Get() {
	tc.Data["json"] = models.GetTag()
	tc.ServeJSON() 
}

func (tc *TagController) Post() {
	logs.Debug("handle Tag POST request")
	jsonStr := tc.Ctx.Input.RequestBody
	tc.Data["json"] = models.SetTag(jsonStr)
	tc.ServeJSON()
}

type TagUpdateController struct {
	web.Controller
}

func (tuc *TagUpdateController) Put() {
	logs.Debug("handle TagUpdate POST request")
	jsonStr := tuc.Ctx.Input.RequestBody
	tagName := tuc.Ctx.Input.Param(":tag_name")
	tuc.Data["json"] = models.UpdateTag(tagName, jsonStr)
	tuc.ServeJSON()
}

type TagActionController struct {
	web.Controller
}

func (tuc *TagActionController) Put() {
	tagName := tuc.Ctx.Input.Param(":tag_name")
	atcion := tuc.Ctx.Input.Param(":action") 
	tuc.Data["json"] = models.TagAction(tagName,atcion)
	tuc.ServeJSON()
}



 
  