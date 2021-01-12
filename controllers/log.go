package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"openkylin.com/ha-api/models"
)

type LogController struct {
	web.Controller
}

func (lc *LogController) Get() {
	logs.Debug("handle resource GET request")
	result := models.GenerateLog()
	lc.Data["json"] = &result
	lc.ServeJSON()
}

func (lc *LogController) Put() {
	logs.Debug("handle resource PUT request")

	lc.ServeJSON()
}

func (lc *LogController) Post() {
	logs.Debug("handle resource POST request")

	lc.ServeJSON()
}
