/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Wed Mar 13 15:34:21 2024 +0800
 */
package controllers

import (
	"encoding/json"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type ScriptsController struct {
	web.Controller
}

type ScriptsRemoteController struct {
	web.Controller
}

func (sc *ScriptsController) Get() {
	scriptName := sc.GetString("filename")
	sc.Data["json"] = models.IsScriptExist(scriptName)
	sc.ServeJSON()
}

func (sc *ScriptsController) Post() {
	logs.Debug("handle scripts POST request")
	data := make(map[string]string)
	if err := json.Unmarshal(sc.Ctx.Input.RequestBody, &data); err != nil {
		result := utils.HandleJsonError(err.Error(), false)
		logs.Error("RequestBody Json parsing failed")
		sc.Data["json"] = &result
		sc.ServeJSON()
	}

	result := models.GenerateScript(data)
	sc.Data["json"] = result
	sc.ServeJSON()
}

func (sc *ScriptsRemoteController) Post() {
	logs.Debug("handle remote scripts POST request")
	data := make(map[string]string)
	if err := json.Unmarshal(sc.Ctx.Input.RequestBody, &data); err != nil {
		result := utils.HandleJsonError(err.Error(), false)
		logs.Error("RequestBody Json parsing failed")
		sc.Data["json"] = &result
		sc.ServeJSON()
	}
	sc.Data["json"] = models.GenerateLocalScript(data)
	sc.ServeJSON()
}
