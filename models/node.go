/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

func GetNodesInfo() ([]map[string]string, error) {
	result := []map[string]string{}

	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	if err != nil {
		if !IsClusterExist() {
			nodeOffline, err2 := GetHeartBeatHosts()
			if err2 != nil {
				return nil, errors.New(gettext.Gettext("Please make sure that Cluster nodes has been authenticated"))
			}
			for _, node := range nodeOffline {
				infoMap := map[string]string{}
				infoMap["id"] = node.NodeID
				infoMap["status"] = "Not Running"
				infoMap["is_dc"] = "false"
				infoMap["type"] = "primitive"
				infoMap["res"] = ""
				result = append(result, infoMap)
			}
			if len(result) > 0 {
				return result, nil
			}
			return nil, errors.New(gettext.Gettext("get node failed"))
		} else {
			nodes_info := GetClusterInfo1()
			nodesOfflineData := nodes_info.Data
			for _, node := range nodesOfflineData {
				infoMap := map[string]string{}
				infoMap["id"] = node.Name
				infoMap["status"] = "Not Running"
				infoMap["is_dc"] = "false"
				infoMap["health_status"] = "healthy"
				infoMap["type"] = "primitive"
				infoMap["res"] = ""
				result = append(result, infoMap)
			}
			if len(result) > 0 {
				return result, nil
			}
			return nil, errors.New(gettext.Gettext("get node failed"))
		}
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return nil, errors.New(gettext.Gettext("parse xml failed"))
	}
	nodes := doc.SelectElement("crm_mon").SelectElement("nodes")
	for _, node := range nodes.SelectElements("node") {
		nodetype := "primitive"
		noderes := ""

		if node.SelectAttr("type").Value == "remote" {
			if node.SelectAttr("id_as_resource") == nil {
				nodetype = "remote"
			} else {
				nodetype = "guest"
				noderes = node.SelectAttr("id_as_resource").Value
			}
		}
		name := node.SelectAttr("name").Value
		online := node.SelectAttr("online").Value
		standby := node.SelectAttr("standby").Value
		isDc := node.SelectAttr("is_dc").Value
		healthStatus := GetHealthInfo(name)
		var status string

		if isDc == "true" {
			if healthStatus == "healthy" {
				if standby == "true" {
					if online == "true" {
						status = "Master/Standby"
					} else {
						status = "Not Running"
					}
				} else {
					if online == "true" {
						status = "Master"
					} else {
						status = "Not Running"
					}
				}
			} else {
				if standby == "true" {
					if online == "true" {
						status = "Master/Standby"
					} else {
						status = "Not Running"
					}
				} else {
					if online == "true" {
						status = "Master/unhealthy"
					} else {
						status = "Not Running"
					}
				}
			}
		} else {
			if healthStatus == "healthy" {
				if standby == "true" {
					if online == "true" {
						status = "Standby"
					} else {
						status = "Not Running"
					}
				} else {
					if online == "true" {
						status = "Running"
					} else {
						status = "Not Running"
					}
				}
			} else {
				if standby == "true" {
					if online == "true" {
						status = "Standby"
					} else {
						status = "Not Running"
					}
				} else {
					if online == "true" {
						status = "unhealthy"
					} else {
						status = "Not Running"
					}
				}
			}
		}

		infoMap := map[string]string{}
		infoMap["id"] = name
		infoMap["status"] = status
		infoMap["is_dc"] = isDc
		infoMap["type"] = nodetype
		infoMap["res"] = noderes
		result = append(result, infoMap)
	}

	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New(gettext.Gettext("get node failed"))
}

// 获取普通节点的数量
func GetPrimitiveNodesInfo() int {
	nodesInfo, err := GetNodesInfo()
	if err != nil {
		return 0
	}
	nodeNum := 0
	for _, node := range nodesInfo {
		if node["type"] == "primitive" {
			nodeNum += 1
		}
	}
	return nodeNum
}

func GetNodeIDInfo(nodeID string) (map[string][]string, error) {
	cmd := "cat /etc/hosts|grep " + nodeID + "|awk -F ' ' '{print $1}'"
	out, err := utils.RunCommand(cmd)

	ips := strings.Split(strings.TrimSpace(string(out)), "\n")

	if err != nil || len(ips) == 0 {
		return nil, err
	}

	nodeInfo := make(map[string][]string)
	nodeInfo["ips"] = ips
	return nodeInfo, nil
}

func DoNodeAction(nodeID string, action string, data map[string]string) map[string]interface{} {
	var cmd string
	result := map[string]interface{}{}
	nodeType := data["type"]
	nodeRes := data["res"]
	if action == "standby" {
		cmd = utils.CmdNodeStandby + nodeID
	} else if action == "unstandby" {
		cmd = utils.CmdNodeUnStandby + nodeID
	} else if action == "start" {
		cmd = handleNodeAction(action, nodeID, nodeType, nodeRes)
	} else if action == "stop" {
		cmd = handleNodeAction(action, nodeID, nodeType, nodeRes)
	} else if action == "restart" {
		cmd = utils.CmdStopCluster + nodeID + " && " + utils.CmdStopCluster + nodeID
	}
	// TODO --force
	if _, err := utils.RunCommand(cmd); err != nil {
		slog.Error(fmt.Sprintf("run command error: %s", err))
		result["action"] = false
		result["error"] = gettext.Gettext("Change node status Failed")
	}

	result["action"] = true
	result["info"] = gettext.Gettext("Change node status success")
	return result
}

// handleNodeAction 根据节点类型和操作生成相应的命令
func handleNodeAction(action string, nodeID, nodeType, nodeRes string) string {
	commands := map[string]map[string]string{
		"start": {
			"primitive": fmt.Sprintf(utils.CmdStartClusterNode, nodeID),
			"remote":    fmt.Sprintf(utils.CmdStartClusterRemoteNode, nodeID),
			"guest":     fmt.Sprintf(utils.CmdStartClusterRemoteNode, nodeRes),
		},
		"stop": {
			// b#371753 add --force for: 3 nodes when 1 node has been stopped
			// In a three-node cluster, if one node is stopped first, and then another node is attempted to be stopped, the second node cannot be stopped normally.
			// The command reports an error, but $?=0 after executing the command, which is considered as no problem in the code
			"primitive": fmt.Sprintf(utils.CmdStopClusterNodeForce, nodeID),
			"remote":    fmt.Sprintf(utils.CmdStopClusterRemoteNode, nodeID),
			"guest":     fmt.Sprintf(utils.CmdStopClusterRemoteNode, nodeRes),
		},
	}
	sleepStr := utils.DefaultSleep

	if _, ok := commands[action]; !ok {
		return fmt.Sprintf("invalid action: %s", action)
	}

	if _, ok := commands[action][nodeType]; !ok {
		return fmt.Sprintf("invalid node type: %s for action: %s", nodeType, action)
	}

	if action == "stop" && (nodeType == "remote" || nodeType == "guest") {
		sleepStr = ""
	}
	cmd := commands[action][nodeType] + sleepStr
	return cmd
}
