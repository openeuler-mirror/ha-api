/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"errors"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

func GetNodesInfo() ([]map[string]string, error) {
	result := []map[string]string{}

	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	if err != nil {
		nodeOffline, err2 := GetHeartBeatHosts()
		if err2 != nil {
			return nil, errors.New(gettext.Gettext("Please make sure that Cluster nodes has been authenticated"))
		}
		for _, node := range nodeOffline {
			infoMap := map[string]string{}
			infoMap["id"] = node.NodeID
			infoMap["status"] = "Not Running"
			infoMap["is_dc"] = "false"
			result = append(result, infoMap)
		}
		if len(result) > 0 {
			return result, nil
		}
		return nil, errors.New(gettext.Gettext("get node failed"))
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return nil, errors.New(gettext.Gettext("parse xml failed"))
	}
	nodes := doc.SelectElement("crm_mon").SelectElement("nodes")
	for _, node := range nodes.SelectElements("node") {
		name := node.SelectAttr("name").Value
		online := node.SelectAttr("online").Value
		standby := node.SelectAttr("standby").Value
		isDc := node.SelectAttr("is_dc").Value
		var status string

		if isDc == "true" {
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
		}

		infoMap := map[string]string{}
		infoMap["id"] = name
		infoMap["status"] = status
		infoMap["is_dc"] = isDc
		result = append(result, infoMap)
	}

	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New(gettext.Gettext("get node failed"))
}

func GetNodeIDInfo(nodeID string) (map[string][]string, error) {
	cmd := "cat /etc/hosts|grep " + nodeID + "|awk -F ' ' '{print $1}'"
	out, err := utils.RunCommand(cmd)

	ips := strings.Split(strings.TrimSpace(string(out)), "\n")
	logs.Debug(ips)

	if err != nil || len(ips) == 0 {
		return nil, err
	}

	nodeInfo := make(map[string][]string)
	nodeInfo["ips"] = ips
	return nodeInfo, nil
}

func DoNodeAction(nodeID, action string) map[string]interface{} {
	var cmd string
	result := map[string]interface{}{}

	if action == "standby" {
		cmd = utils.CmdNodeStandby + nodeID
	} else if action == "unstandby" {
		cmd = utils.CmdNodeUnStandby + nodeID
	} else if action == "start" {
		cmd = utils.CmdStartCluster + nodeID + " &sleep 5"
	} else if action == "stop" {
		cmd = utils.CmdStopCluster + nodeID
	} else if action == "restart" {
		cmd = utils.CmdStopCluster + nodeID + "||" + utils.CmdStopCluster + nodeID
	}

	if _, err := utils.RunCommand(cmd); err != nil {
		logs.Error("run command error: ", err)
		result["action"] = false
		result["error"] = gettext.Gettext("Change node status Failed")
	}

	result["action"] = true
	result["info"] = gettext.Gettext("Change node status success")
	return result
}
