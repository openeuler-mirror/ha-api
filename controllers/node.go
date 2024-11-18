/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-19 17:37:46
 * Description: 节点控制器
 ******************************************************************************/
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
