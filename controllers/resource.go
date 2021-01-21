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

	rc.Data["json"] = models.GetResourceInfo()
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

type ResourceMetaAttributesController struct {
	web.Controller
}

func (rc *ResourceMetaAttributesController) Get() {
	catagory := rc.Ctx.Input.Param(":catagory")
	result := models.GetResourceMetaAttributes(catagory)
	rc.Data["json"] = &result
	rc.ServeJSON()
}

type ResourceOpsById struct {
	web.Controller
}

func (robi *ResourceOpsById) Get() {
	rscID := robi.Ctx.Input.Param(":rscID")
	result := models.GetResourceInfoByrscID(rscID)
	robi.Data["json"] = &result
	robi.ServeJSON()
}

type ResourceRelationsController struct {
	web.Controller
}

func (rrc *ResourceRelationsController) Get() {
	logs.Debug("handle resource relation GET request")
	rscID := rrc.Ctx.Input.Param(":rscID")
	relation := rrc.Ctx.Input.Param(":relation")

	result := map[string]interface{}{}

	retData, err := models.GetResourceConstraints(rscID, relation)
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = retData
	}

	rrc.Data["json"] = &result
	rrc.ServeJSON()
}
