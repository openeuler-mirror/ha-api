/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-19 11:19:33
 * LastEditTime: 2024-03-19 11:27:38
 * Description: 脚本生成器模块
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
