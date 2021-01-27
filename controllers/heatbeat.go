package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type HeartBeatController struct {
	web.Controller
}

func (hbc *HeartBeatController) Get() {
	result := map[string]interface{}{}

	data, err := models.GetHeartBeatConfig()
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = data
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

func (hbc *HeartBeatController) Post() {
	result := map[string]interface{}{}

	data := hbc.Ctx.Input.RequestBody
	if err := models.EditHeartbeatInfo(data); err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = "Change cluster success"
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

type HeartBeatStatusController struct {
	web.Controller
}

func (hbsc *HeartBeatStatusController) Get() {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = 0

	hbsc.Data["json"] = &result
	hbsc.ServeJSON()
}
