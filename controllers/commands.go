/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */
package controllers

import (
	"strconv"

	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/server/web"
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
