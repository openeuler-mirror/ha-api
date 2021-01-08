package models

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"openkylin.com/ha-api/utils"
)

func GetNodesInfo() map[string]interface{} {
	var result map[string]interface{}

	out, err := utils.RunCommand("crm_mon --as-xml")
	if err != nil {
		// get nodes info from hosts
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		// TODO: not finished
	}

	return result
}

func DoNodeAction(nodeID, action string) map[string]interface{} {
	var cmd string
	result := map[string]interface{}{}

	if action == "standby" {
		// 备用
		cmd = "pcs node standby " + nodeID
	} else if action == "unstandby" {
		// 不备用
		cmd = "pcs node unstandby " + nodeID
	} else if action == "start" {
		// 启动
		cmd = "pcs cluster start " + nodeID + " &sleep 5"
	} else if action == "stop" {
		// 停止
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
