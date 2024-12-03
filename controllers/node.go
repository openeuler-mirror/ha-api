/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
)

type NodesController struct {
	web.Controller
}

func (nsc *NodesController) Get() {
	result := map[string]interface{}{}

	retData, err := models.GetNodesInfo()
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = retData
	}

	nsc.Data["json"] = &result
	nsc.ServeJSON()
}

type NodeController struct {
	web.Controller
}

func (nc *NodeController) Get() {
	result := map[string]interface{}{}

	nodeID := nc.Ctx.Input.Param(":nodeID")
	retData, err := models.GetNodeIDInfo(nodeID)
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = retData
	}

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
