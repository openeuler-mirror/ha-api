package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type NodesController struct {
	web.Controller
}

func (nsc *NodesController) Get() {
	result := models.GetNodesInfo()

	nsc.Data["json"] = &result
	nsc.ServeJSON()
}

type NodeController struct {
	web.Controller
}

func (nc *NodeController) Get() {
	result := []string{}

	nc.Data["json"] = &result
	nc.ServeJSON()
}

type NodeActionController struct {
	web.Controller
}

func (nac *NodeActionController) Put() {
	nodeID := nac.Ctx.Input.Param(":nodeID")
	action := nac.Ctx.Input.Param(":action")
	result := models.DoNodeAction(nodeID, action)

	nac.Data["json"] = &result
	nac.ServeJSON()
}