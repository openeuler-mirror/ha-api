package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"openkylin.com/ha-api/models"
)

type ResourceController struct {
	web.Controller
}

func (rc *ResourceController) Get() {
	logs.Debug("handle resource GET request")

	rc.Data["json"] = models.GerResourceInfo()
	rc.ServeJSON()
}

func (rc *ResourceController) Post() {
	logs.Debug("handle resource POST request")

	jsonStr := rc.Ctx.Input.RequestBody
	rc.Data["json"] = models.CreateResource(jsonStr)
	rc.ServeJSON()
}

func (rc *ResourceController) PUT() {
	logs.Debug("handle resource PUT request")
	logs.Debug("do nothing")

	rc.ServeJSON()
}

type ResourceActionController struct {
	web.Controller
}

func (rac *ResourceActionController) PUT() {
	logs.Debug("handle resource action PUT request")

	rscID := rac.Ctx.Input.Param(":rscID")
	action := rac.Ctx.Input.Param(":action")
	data := rac.Ctx.Input.RequestBody

	var result map[string]interface{}
	if err := models.ResourceAction(rscID, action, data); err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = false
		result["info"] = "Action on resource success"
	}

	rac.Data["json"] = &result
	rac.ServeJSON()
}
