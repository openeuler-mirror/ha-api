/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package controllers

import (
	"encoding/json"
	"log/slog"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
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

type ClusterOverviewController struct {
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

type DeleteNodesController struct {
	web.Controller
}

type LocalAddNodesController struct {
	web.Controller
}

type LocalDeleteNodesController struct {
	web.Controller
}
type LocalClusterInfoController struct {
	web.Controller
}

type LocalClusterOverviewController struct {
	web.Controller
}

type IsClusterExistController struct {
	web.Controller
}

type LocalClusterSetupController struct {
	web.Controller
}

func (coc *ClusterOverviewController) Get() {
	slog.Debug("handle get request in ClusterOverviewController.")
	result := models.ClusterOverview()
	coc.Data["json"] = &result
	coc.ServeJSON()
}

func (lcoc *LocalClusterOverviewController) Get() {
	slog.Debug("handle get request in ClusterOverviewController.")
	result := models.LocalClusterOverview()
	lcoc.Data["json"] = &result
	lcoc.ServeJSON()
}
func (mcc *MultipleClustersController) Post() {
	slog.Debug("Handle post request in MultipleClustersController.")
	var result utils.GeneralResponse

	var reqData models.ClusterAddReq
	//need to check whether we want to add or remove cluster
	if err := json.Unmarshal(mcc.Ctx.Input.RequestBody, &reqData); err != nil {
		result.Action = false
		result.Error = gettext.Gettext("invalid input data")
	} else {
		result = models.ClusterAdd(reqData)
	}

	mcc.Data["json"] = &result
	mcc.ServeJSON()

}

func (sc *Sync_configController) Post() {
	slog.Debug("handle post request in Sync_configController.")
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
	slog.Debug("handle post request in ClusterSetupController.")
	result := make(map[string]interface{})
	var reqData models.ClusterData
	if err := json.Unmarshal(csc.Ctx.Input.RequestBody, &reqData); err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.ClusterSetup(reqData)
	}

	csc.Data["json"] = &result
	csc.ServeJSON()
}

func (scc *LocalClusterSetupController) Post() {
	slog.Debug("handle post request in SetupClusterController.")
	result := make(map[string]interface{})
	var reqData models.ClusterData
	if err := json.Unmarshal(scc.Ctx.Input.RequestBody, &reqData); err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {

		result = models.ClusterSetupPre(reqData)
	}

	scc.Data["json"] = &result
	scc.ServeJSON()
}

func (cc *ClustersController) Get() {
	slog.Debug("handle get request in ClustersController.")
	result := make(map[string]interface{})
	result = models.GetClusterPropertiesInfo()
	cc.Data["json"] = &result
	cc.ServeJSON()
}

func (csc *ClustersStatusController) Get() {
	var result models.NodeManageClusterInfo
	slog.Debug("handle get request in ClustersStatusController.")
	result = models.GetClusterInfo1()
	csc.Data["json"] = &result
	csc.ServeJSON()
}

func (lci *LocalClusterInfoController) Get() {
	slog.Debug("handle get request in LocalClusterInfoController.")
	result := models.LocalClusterInfo()
	lci.Data["json"] = &result
	lci.ServeJSON()
}

func (ice *IsClusterExistController) Get() {
	slog.Debug("handle get request in IsClusterExistController.")
	result := models.CheckIsClusterExist()
	ice.Data["json"] = &result
	ice.ServeJSON()
}

func (lcd *LocalClusterDestroyController) Get() {
	slog.Debug("handle post request in LocalClusterDestroyController.")
	result := models.LocalClusterDestroy()
	lcd.Data["json"] = &result
	// return result of destroying cluster back to user.
	lcd.ServeJSON()
}

func (cd *ClusterDestroyController) Post() {
	slog.Debug("handle post request in ClusterDestroyController.")
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
	slog.Debug("handle post request in ClusterRemoveController.")
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
	var reqData models.AddNodesData
	var result utils.GeneralResponse
	var resultRemote map[string]interface{}
	if err := json.Unmarshal(anc.Ctx.Input.RequestBody, &reqData); err != nil {
		result.Action = false
		result.Error = gettext.Gettext("invalid input data")
		anc.Data["json"] = &result
		anc.ServeJSON()
		return
	}
	clusterName := reqData.Cluster_name
	if !utils.IsLocalCluster(clusterName) {
		remoteUiPath := "/api/v1/nodes/add_nodes"
		resultRemote, _ = models.UrlRedirectWithSyncConfig(clusterName, remoteUiPath, anc.Ctx.Request.Method, anc.Ctx.Input.RequestBody)
		anc.Data["json"] = &resultRemote
	} else {
		resultRemote = models.LocalAddNodes(reqData)
		anc.Data["json"] = &resultRemote
	}
	anc.ServeJSON()
}

func (lanc *LocalAddNodesController) Post() {
	slog.Debug("handle post request in LocalAddNodesController.")
	result := make(map[string]interface{})
	var ReqData models.AddNodesData
	body := lanc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		// result = models.LocalAddNodes(ReqData).(map[string]interface{})
		result = models.LocalAddNodes(ReqData)
	}
	lanc.Data["json"] = &result
	lanc.ServeJSON()
}

func (dnc *DeleteNodesController) Post() {
	slog.Debug("handle post request in DeleteNodesController.")
	result := make(map[string]interface{})
	var ReqData models.DeleteNodesData
	body := dnc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
		dnc.Data["json"] = &result
	} else {
		clusterName := ReqData.Cluster_name
		if !utils.IsLocalCluster(clusterName) {
			remoteUiPath := "/api/v1/nodes/delete_nodes"
			resultRemote, _ := models.UrlRedirectWithSyncConfig(clusterName, remoteUiPath, dnc.Ctx.Request.Method, dnc.Ctx.Input.RequestBody)
			dnc.Data["json"] = &resultRemote
		} else {
			result = models.LocalDeleteNodes(ReqData).(map[string]interface{})
			dnc.Data["json"] = &result
		}
	}
	dnc.ServeJSON()
}

func (ldnc *LocalDeleteNodesController) Post() {
	slog.Debug("handle post request in DeleteNodesController.")
	result := make(map[string]interface{})
	var ReqData models.DeleteNodesData
	body := ldnc.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &ReqData)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		result = models.LocalDeleteNodes(ReqData).(map[string]interface{})
	}
	ldnc.Data["json"] = &result
	ldnc.ServeJSON()
}

func (cc *ClustersController) Put() {
	slog.Debug("handle put request in ClustersController.")
	result := make(map[string]interface{})

	reqData := make(map[string]interface{})
	if err := json.Unmarshal(cc.Ctx.Input.RequestBody, &reqData); err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("invalid input data")
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
