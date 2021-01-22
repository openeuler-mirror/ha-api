package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type MetaController struct {
	web.Controller
}

func (mc *MetaController) Get() {
	mc.Data["json"] = models.getResourceMetas()
	mc.ServeJSON()
}

type MetasController struct {
	web.Controller
}

func (mc *MetasController) Get() {
	rscClass := mc.Ctx.Input.Param(":rsc_class")
	rscType := mc.Ctx.Input.Param(":rsc_type")
	rscProvider := mc.Ctx.Input.Param(":rsc_provider")
	mc.Data["json"] = models.getAllResourceMetas(rscClass, rscType, rscProvider)
	mc.ServeJSON()
}
