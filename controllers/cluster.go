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
	"openkylin.com/ha-api/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type ClustersController struct {
	web.Controller
}

type MultipleClustersController struct {
	web.Controller
}

type Sync_configController struct {
	web.Controller
}

type ClusterSetupController struct {
	web.Controller
}

type ClusterDestroyController struct {
	web.Controller
}

type ClustersStatusController struct {
	web.Controller
}

type ClusterRemoveController struct {
	web.Controller
}

func (mcc *MultipleClustersController) Post() {
	logs.Debug("Handle post request in MultipleClustersController.")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	//need to check whether we want to add or remove cluster
	if err := json.Unmarshal(mcc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.ClusterAdd(reqData)
	}

	mcc.Data["json"] = &result
	mcc.ServeJSON()

}

func (sc *Sync_configController) Post() {
	logs.Debug("handle post request in ClustersController.")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	if err := json.Unmarshal(sc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.SyncConfig(reqData)
	}
	sc.Data["json"] = &result
	sc.ServeJSON()
}

func (csc *ClusterSetupController) Post() {
	logs.Debug("handle post request in ClustersController.")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	if err := json.Unmarshal(csc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.ClusterSetup(reqData)
	}

	csc.Data["json"] = &result
	csc.ServeJSON()
}

//
//func (mcc *MultipleClustersController) Get() {
//	logs.Debug("Handle get request in MultipleClustersController.")
//	//need to return cluster data through res.
//
//}

func (cc *ClustersController) Get() {
	logs.Debug("handle get request in ClustersController.")
	result := models.GetClusterPropertiesInfo()
	cc.Data["json"] = &result
	cc.ServeJSON()
}

func (csc *ClustersStatusController) Get() {
	logs.Debug("handle get request in ClustersController.")
	result := models.GetClusterInfo()
	csc.Data["json"] = &result
	csc.ServeJSON()
}

func (cc *ClustersController) Post() {
	logs.Debug("handle post request in ClustersController.")
	cc.ServeJSON()
}

func (cd *ClusterDestroyController) Post() {
	logs.Debug("handle post request in ClustersController.")
	result := models.ClustersDestroy()
	cd.Data["json"] = &result
	// return result of destroying cluster back to user.
	cd.ServeJSON()
}

func (crc *ClusterRemoveController) Post() {
	logs.Debug("handle post request in ClustersController.")
	var Result models.RemoveRet
	var ReqData models.RemoveData
	body := crc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		Result.Action = false
		Result.Error = "invalid input data"
		crc.Data["json"] = &Result
		crc.ServeJSON()
	} else {
		Result2 := models.ClusterRemove(ReqData)
		crc.Data["json"] = Result2
		crc.ServeJSON()
	}
}

func (cc *ClustersController) Put() {
	logs.Debug("handle put request in ClustersController.")
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
