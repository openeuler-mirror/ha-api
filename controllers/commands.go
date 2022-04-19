/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-19 17:37:58
 * Description: 命令控制器
 ******************************************************************************/
package controllers

import (
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type CommandsController struct {
	web.Controller
}

func (cc *CommandsController) Get() {
	result := models.GetCommandsList()
	cc.Data["json"] = &result
	cc.ServeJSON()
}

type CommandsRunnerController struct {
	web.Controller
}

func (crc *CommandsRunnerController) Get() {
	result := map[string]interface{}{}

	t := crc.Ctx.Input.Param(":cmd_type")
	cmdID, err := strconv.Atoi(t)
	if err != nil {
		result["action"] = false
		result["error"] = err
	} else {
		out, err := models.RunBuiltinCommand(cmdID)
		if err != nil {
			result["action"] = false
			result["error"] = err
		} else {
			result["action"] = true
			result["data"] = out
		}
	}
	crc.Data["json"] = &result
	crc.ServeJSON()
}
