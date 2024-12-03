/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Jason011125 <zic022@ucsd.edu>
 * Date: Mon Aug 14 13:38:45 2023 +0800
 */
package controllers

import (
	"encoding/json"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"
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

type LocalClusterDestroyController struct {
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

type AddNodesController struct {
	web.Controller
}

type LocalAddNodesController struct {
	web.Controller
}

type LocalClusterInfoController struct {
	web.Controller
}

type IsClusterExistController struct {
	web.Controller
}

func (mcc *MultipleClustersController) Post() {
	logs.Debug("Handle post request in MultipleClustersController.")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	//need to check whether we want to add or remove cluster
	if err := json.Unmarshal(mcc.Ctx.Input.RequestBody, &reqData); err != nil {
		// result = make(map[string]interface{})
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.ClusterAdd(reqData)
	}

	mcc.Data["json"] = &result
	mcc.ServeJSON()

}

func (sc *Sync_configController) Post() {
	logs.Debug("handle post request in Sync_configController.")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	if err := json.Unmarshal(sc.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.SyncConfig(reqData)
	}
	sc.Data["json"] = &result
	sc.ServeJSON()
}

func (csc *ClusterSetupController) Post() {
	logs.Debug("handle post request in ClusterSetupController.")
	result := make(map[string]interface{})
	//reqData := make(map[string]interface{})
	var reqData models.ClusterData
	if err := json.Unmarshal(csc.Ctx.Input.RequestBody, &reqData); err != nil {
		// result = make(map[string]interface{})
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.ClusterSetup(reqData)
	}

	csc.Data["json"] = &result
	csc.ServeJSON()
}

func (cc *ClustersController) Get() {
	logs.Debug("handle get request in ClustersController.")
	clusterName := cc.Ctx.Input.Param(":cluster_name")
	result := make(map[string]interface{})
	if !utils.IsLocalCluster(clusterName) {
		result, _ = models.UrlRedirect(clusterName, cc.Ctx.Input.URL(), cc.Ctx.Request.Method, cc.Ctx.Input.RequestBody)
	} else {
		result = models.GetClusterPropertiesInfo()
	}

	cc.Data["json"] = &result
	cc.ServeJSON()
}

func (csc *ClustersStatusController) Get() {
	var result map[string]interface{}
	logs.Debug("handle get request in ClustersStatusController.")
	clusterName := csc.Ctx.Input.Param(":cluster_name")

	if !utils.IsLocalCluster(clusterName) {
		result, _ = models.UrlRedirect(clusterName, csc.Ctx.Input.URL(), csc.Ctx.Request.Method, csc.Ctx.Input.RequestBody)
	} else {
		result = models.GetClusterInfo()
	}

	csc.Data["json"] = &result
	csc.ServeJSON()
}

func (lci *LocalClusterInfoController) Get() {
	logs.Debug("handle get request in LocalClusterInfoController.")
	result := models.LocalClusterInfo()
	lci.Data["json"] = &result
	lci.ServeJSON()
}

func (ice *IsClusterExistController) Get() {
	logs.Debug("handle get request in IsClusterExistController.")
	result := models.IsClusterExist()
	ice.Data["json"] = &result
	ice.ServeJSON()
}

func (lcd *LocalClusterDestroyController) Get() {
	logs.Debug("handle post request in LocalClusterDestroyController.")
	result := models.LocalClustersDestroy()
	lcd.Data["json"] = &result
	// return result of destroying cluster back to user.
	lcd.ServeJSON()
}

func (cd *ClusterDestroyController) Post() {
	logs.Debug("handle post request in ClusterDestroyController.")
	//var Result models.AddNodesRet
	result := make(map[string]interface{})
	ReqData := make(map[string]interface{})
	body := cd.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.ClusterDestroy(ReqData)
	}
	cd.Data["json"] = &result
	cd.ServeJSON()
}

func (crc *ClusterRemoveController) Post() {
	logs.Debug("handle post request in ClusterRemoveController.")
	var Result models.RemoveRet
	var ReqData models.RemoveData
	body := crc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		Result.Action = false
		Result.Error = gettext.Gettext("invalid input data")
		crc.Data["json"] = &Result
		crc.ServeJSON()
	} else {
		Result2 := models.ClusterRemove(ReqData)
		crc.Data["json"] = Result2
		crc.ServeJSON()
	}
}

func (anc *AddNodesController) Post() {
	logs.Debug("handle post request in AddNodesController.")
	//var Result models.AddNodesRet
	result := map[string]interface{}{}
	var ReqData models.AddNodesData
	body := anc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.AddNodes(ReqData).(map[string]interface{})
	}
	anc.Data["json"] = &result
	anc.ServeJSON()
}

func (lanc *LocalAddNodesController) Post() {
	logs.Debug("handle post request in LocalAddNodesController.")
	//var Result models.AddNodesRet
	result := make(map[string]interface{})
	var ReqData models.AddNodesData
	body := lanc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.LocalAddNodes(ReqData).(map[string]interface{})
	}
	lanc.Data["json"] = &result
	lanc.ServeJSON()
}

func (cc *ClustersController) Put() {
	logs.Debug("handle put request in ClustersController.")
	result := make(map[string]interface{})
	clusterName := cc.Ctx.Input.Param(":cluster_name")

	if !utils.IsLocalCluster(clusterName) {
		result, _ = models.UrlRedirect(clusterName, cc.Ctx.Input.URL(), cc.Ctx.Request.Method, cc.Ctx.Input.RequestBody)
	} else {
		reqData := make(map[string]interface{})
		if err := json.Unmarshal(cc.Ctx.Input.RequestBody, &reqData); err != nil {
			result["action"] = false
			result["error"] = gettext.Gettext("invalid input data")
		} else {
			result = models.UpdateClusterProperties(reqData)
		}
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
