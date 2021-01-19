package models

import (
	"errors"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"openkylin.com/ha-api/utils"
)

func GetNodesInfo() ([]map[string]string, error) {
	var result []map[string]string

	out, err := utils.RunCommand("crm_mon --as-xml")
	if err != nil {
		nodeOffline, err2 := GetHeartBeatHosts()
		if err2 != nil {
			return nil, errors.New("请确定集群节点已认证")
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
		return nil, errors.New("Get node failed")
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return nil, errors.New("parse xml failed")
	}
	nodes := doc.SelectElement("crm_mon").SelectElements("nodes")
	for _, node := range nodes {
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
	return nil, errors.New("Get node failed")
}

func GetNodeIDInfo(nodeID string) (map[string][]string, error) {
	var cmd string
	cmd = "cat /etc/hosts|grep " + nodeID + "|awk -F ' ' '{print $1}'"
	out, err := utils.RunCommand(cmd)

	var ips []string
	ips = strings.Split(strings.TrimSpace(string(out)), "\n")
	logs.Debug(ips)

	if err != nil || len(ips) == 0 {
		return nil, err
	}

	var nodeInfo = make(map[string][]string)
	nodeInfo["ips"] = ips
	return nodeInfo, nil
}

func DoNodeAction(nodeID, action string) map[string]interface{} {
	var cmd string
	result := map[string]interface{}{}

	if action == "standby" {
		cmd = "pcs node standby " + nodeID
	} else if action == "unstandby" {
		cmd = "pcs node unstandby " + nodeID
	} else if action == "start" {
		cmd = "pcs cluster start " + nodeID + " &sleep 5"
	} else if action == "stop" {
		cmd = "pcs cluster stop " + nodeID
	} else if action == "restart" {
		cmd = "pcs cluster restart " + nodeID
	}

	if _, err := utils.RunCommand(cmd); err != nil {
		logs.Error("run command error: ", err)
		result["action"] = false
		result["error"] = "Change node status Failed"
	}

	result["action"] = true
	result["error"] = "Change node status success"
	return result
}
