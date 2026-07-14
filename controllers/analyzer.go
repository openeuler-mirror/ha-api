/*
 * Copyright (c) KylinSoft  Co., Ltd. 2027.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: xuxiaojuan <xuxiaojuan@kylinos.cn>
 * Date: Wed July 8 13:56:40 2026 +0800
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/server/web"
)

func urlNodeRedirect(clusterName string, uiPath string, node string, requestMethod string, requestData interface{}) (map[string]interface{}, error) {
	port, err := utils.ReadPortFromConfig()
	if err != nil {
		slog.Warn("ReadPortFromConfig failed, using default port", "err", err)
	}

	allowedNodes := models.GetLocalConf().GetNodes(clusterName)
	if !containsNode(allowedNodes, node) {
		errMsg := fmt.Sprintf("node %q is not a known member of cluster %q", node, clusterName)
		slog.Error("SSRF attempt blocked", "node", node, "cluster", clusterName)
		return map[string]interface{}{"action": false, "error": errMsg}, fmt.Errorf("%s", errMsg)
	}

	url := "https://" + node + ":" + port + "/remote" + uiPath
	resp, err := utils.SendRequest(url, requestMethod, requestData)
	if err != nil {
		slog.Error("failed to request remote node", "node", node, "err", err)
		errMsg := "Can not request the remote node " + node
		return map[string]interface{}{"action": false, "error": errMsg}, fmt.Errorf("%s", errMsg)
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response from remote node", "node", node, "err", err)
		errMsg := "Can not read from the remote node " + node
		return map[string]interface{}{"action": false, "error": errMsg}, fmt.Errorf("%s", errMsg)
	}

	remoteClusterInfo := make(map[string]interface{})
	if err := json.Unmarshal(respData, &remoteClusterInfo); err != nil {
		slog.Error("parse response failed", "err", err)
		return map[string]interface{}{"action": false, "error": "Parse response failed"}, fmt.Errorf("parse response failed")
	}

	return remoteClusterInfo, nil
}

func containsNode(nodes []string, target string) bool {
	for _, n := range nodes {
		if n == target {
			return true
		}
	}
	return false
}

func remoteaction(ctrl web.Controller, funcName string) map[string]interface{} {
	result := map[string]interface{}{}
	clusterName := ctrl.Ctx.Input.Param(":cluster_name")
	nodeName := ctrl.Ctx.Input.Param(":node_name")

	if utils.IsLocalCluster(clusterName) && nodeName == "localhost" {
		switch funcName {
		case "system_overview":
			result = models.GetSystemInfo()
		case "service_status":
			result = models.GetServiceStatus()
		default:
			result["action"] = false
			result["error"] = fmt.Sprintf("unknown action: %s", funcName)
		}
	} else {
		segments := strings.Split(ctrl.Ctx.Request.RequestURI, "/")
		if len(segments) > 5 {
			segments[5] = "localhost"
		}
		newurl := strings.Join(segments, "/")

		if nodeName == "localhost" {
			nodes := models.GetLocalConf().GetNodes(clusterName)
			if len(nodes) == 0 {
				result["action"] = false
				result["error"] = fmt.Sprintf("no nodes found for cluster %q", clusterName)
				return result
			}
			nodeName = nodes[0]
		}
		data, err := urlNodeRedirect(clusterName, newurl, nodeName, ctrl.Ctx.Request.Method, nil)
		if err != nil {
			result["action"] = false
			result["error"] = err.Error()
		} else {
			result = data
		}
	}
	return result
}

type SystemInfoController struct {
	web.Controller
}

func (c *SystemInfoController) Get() {
	sysInfo := remoteaction(c.Controller, "system_overview")
	c.Data["json"] = sysInfo
	c.ServeJSON()
}

type ServiceStatusController struct {
	web.Controller
}

func (sst *ServiceStatusController) Get() {
	serviceStatus := remoteaction(sst.Controller, "service_status")
	sst.Data["json"] = serviceStatus
	sst.ServeJSON()
}
