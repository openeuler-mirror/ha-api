package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type MetaController struct {
	web.Controller
}

func (mc *MetaController) Get() {
	// TODO
	// mc.Data["json"] = models.GetResourceMetas()
	mc.ServeJSON()
}

type MetasController struct {
	web.Controller
}

func (mc *MetasController) Get() {
	// TODO
	// rscClass := mc.Ctx.Input.Param(":rsc_class")
	// rscType := mc.Ctx.Input.Param(":rsc_type")
	// rscProvider := mc.Ctx.Input.Param(":rsc_provider")
	// mc.Data["json"] = models.GetAllResourceMetas(rscClass, rscType, rscProvider)
	mc.ServeJSON()
}
