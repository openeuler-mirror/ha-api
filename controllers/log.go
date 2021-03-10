package controllers

import (
	"net/http"

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

type LogDownloadController struct {
	web.Controller
}

func (ldc *LogDownloadController) Get() {
	logs.Debug("handle log download GET request")
	fileTail := ldc.Ctx.Input.Param(":filetail")
	// result := models.GenerateLog()
	// lc.Data["json"] = &result
	// lc.ServeJSON()

	const filePath = "/usr/share/heartbeat-gui/ha-api/static/"
	filePrefix := "kylinha-log-"
	http.ServeFile(ldc.Ctx.ResponseWriter, ldc.Ctx.Request, filePath+filePrefix+fileTail)
}
