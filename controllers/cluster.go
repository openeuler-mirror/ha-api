/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-19 17:37:56
 * Description: HA集群控制器
 ******************************************************************************/
package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"openkylin.com/ha-api/models"
)

type ClustersController struct {
	web.Controller
}

func (cc *ClustersController) Get() {
	logs.Debug("handle get request in HAClustersController.")
	result := models.GetClusterPropertiesInfo()
	cc.Data["json"] = &result
	cc.ServeJSON()
}

func (cc *ClustersController) Post() {
	logs.Debug("handle post request in HAClustersController.")
	// do nothing here
	cc.ServeJSON()
}

func (cc *ClustersController) Put() {
	logs.Debug("handle put request in HAClustersController.")
	result := map[string]interface{}{}

	reqData := make(map[string]interface{})
	if err := json.Unmarshal(cc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.UpdateClusterProperties(reqData)
	}

	cc.Data["json"] = &result
	cc.ServeJSON()
}

type LocalHaOperation struct {
	web.Controller
}

func (lho *LocalHaOperation) Put() {
	action := lho.Ctx.Input.Param("action")
	lho.Data["json"] = models.OperationClusterAction(action)
	lho.ServeJSON()
}
