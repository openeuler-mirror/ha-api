package controllers

import (
	"github.com/beego/beego/v2/server/web"

	"encoding/json"

	"openkylin.com/ha-api/models"
)

type AlarmConfig struct {
	web.Controller
}

func (ac *AlarmConfig) Get() {
	ac.Data["json"] = models.AlarmsGet()
	ac.ServeJSON()
}
func (ac *AlarmConfig) Post() {
	var result map[string]interface{}

	reqData := make(map[string]string)
	if err := json.Unmarshal(ac.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.AlarmsSet(reqData)
	}
}
