package controllers

import (
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
