package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"openkylin.com/ha-api/models"
)

type HAClustersController struct {
	web.Controller
}

func (hcc *HAClustersController) Get() {
	// TODO: handle get here
	logs.Debug("handle get request in HAClustersController.")
	result := models.GetClusterPropertiesInfo()
	hcc.Data["json"] = &result
	hcc.ServeJSON()
}

func (hcc *HAClustersController) Post() {
	logs.Debug("handle post request in HAClustersController.")
	// do nothing here
	hcc.ServeJSON()
}

func (hcc *HAClustersController) Put() {
	logs.Debug("handle put request in HAClustersController.")
	result := map[string]interface{}{}

	reqData := make(map[string]string)
	if err := json.Unmarshal(hcc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.UpdateClusterProperties(reqData)
	}

	hcc.Data["json"] = &result
	hcc.ServeJSON()
}

type LocalHaOperation struct {
	web.Controller
}

func (lho *LocalHaOperation) Put() {
	action := lho.Ctx.Input.Param("action")
	lho.Data["json"] = models.OperationClusterAction(action)
	lho.ServeJSON()
}
