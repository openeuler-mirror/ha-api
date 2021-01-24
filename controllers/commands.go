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
	var result map[string]interface{}

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
