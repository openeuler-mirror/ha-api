package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type MetaController struct {
	web.Controller
}

func (mc *MetaController) Get() {
	rscClass := mc.Ctx.Input.Param(":rsc_class")
	rscType := mc.Ctx.Input.Param(":rsc_type")
	rscProvider := mc.Ctx.Input.Param(":rsc_provider")
	mc.Data["json"] = models.GetResourceMetas(rscClass, rscType, rscProvider)
	mc.ServeJSON()
}

type MetasController struct {
	web.Controller
}

func (mc *MetasController) Get() {
	mc.Data["json"] = models.GetAllResourceMetas()
	mc.ServeJSON()
}
