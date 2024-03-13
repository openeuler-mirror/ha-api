/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-13 14:16:07
 * LastEditTime: 2024-03-13 15:32:03
 * Description:脚本生成器模块
 */
package controllers

import (
	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/server/web"
)

type ScriptsController struct {
	web.Controller
}

func (sc *ScriptsController) Get() {
	scriptName := sc.GetString("filename")
	sc.Data["json"] = models.IsScriptExist(scriptName)
	sc.ServeJSON()
}

// func (sc *ScriptsController) Post() {
// 	logs.Debug("handle scripts POST request")
// 	data := make(map[string]string)
// 	if err := json.Unmarshal(sc.Ctx.Input.RequestBody, &data); err != nil {
// 		// rc.handleJsonError(err.Error(), false)
// 	}
// 	sc.ServeJSON()
// }
