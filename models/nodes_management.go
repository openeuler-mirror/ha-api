/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Jason011125 <zic022@ucsd.edu>
 * Date: Mon Aug 14 15:53:52 2023 +0800
 */
package models

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

// nodes_info格式
// {"nodeid": "1", "name": "HA1", "ring0_addr": "192.168.11.1", "ring1_addr": "192.168.11.4"},
// {"nodeid": "2", "name": "HA2", "ring0_addr": "192.168.11.2", "ring1_addr": "192.168.11.5"}
type AuthInfo struct {
	nodeList []string
	passWord []string
	ip       []string
}

func getLocalClusterName() (string, error) {
	filename := settings.CorosyncConfFile
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("open corosync conf failed")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "cluster_name" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("not found cluster name info in corosync conf")
}

// getClusterName reads the cluster name from the corosync configuration file.
// Returns a map indicating the result and the extracted cluster name, if available.
// func getClusterName() map[string]interface{} {
// 	result := map[string]interface{}{
// 		"action":      false,
// 		"clusterName": "",
// 	}
// 	localClusterName, err := getLocalClusterName()
// 	if err != nil {
// 		return result
// 	}

//		result["action"] = true
//		result["clusterName"] = localClusterName
//		return result
//	}
func getClusterName() string {
	clusterName := ""
	localClusterName, err := getLocalClusterName()
	if err != nil {
		return clusterName
	}

	return localClusterName
}

// getClusterInfo retrieves cluster information, including cluster nodes and their properties.
// Returns the cluster information in a structured map.
func GetClusterInfo() map[string]interface{} {
	currentNode, _ := utils.RunCommand(utils.CmdHostName)
	currentNodeContent := strings.ReplaceAll(string(currentNode), "\n", "")
	clusterName := getClusterName()
	var errorInfo string
	var data map[string]interface{}
	var nodesInfo []map[string]string
	var err error
	processedNodesInfo := make([]map[string]interface{}, 0)

	if !IsClusterExist() {
		errorInfo = "Cluster not established!"
		goto ret
	}

	nodesInfo, err = utils.GetNodeList()
	if err != nil {
		errorInfo = "Cluster not established!"
		goto ret
	}

	for _, node := range nodesInfo {
		newNode := make(map[string]interface{})
		for k, v := range node {
			newNode[k] = v
		}
		newNode["type"] = "primitive"
		processedNodesInfo = append(processedNodesInfo, newNode)
	}

	// add node status
	updateNodeStatus(processedNodesInfo, currentNodeContent)

	data = map[string]interface{}{
		"action":        true,
		"cluster_exist": true,
		"cluster_name":  clusterName,
		"currentNode":   currentNodeContent,
		"data":          processedNodesInfo,
	}
	return data

ret:
	return map[string]interface{}{
		"action":        false,
		"cluster_exist": false,
		"cluster_name":  clusterName,
		"currentNode":   currentNodeContent,
		"error":         errorInfo,
	}
}

// add node status to clusterInfo
func updateNodeStatus(nodesInfo []map[string]interface{}, currentNode string) error {
	cmd := utils.CmdHbStatus
	out, err := utils.RunCommand(cmd)
	if err != nil {
		// Get node status when cluster stoped
		for _, node := range nodesInfo {
			node["status"] = make(map[string]string, 0)
			for k := range node {
				if strings.HasPrefix(k, "ring") {
					if currentNode == node["name"] {
						node["status"].(map[string]string)[k] = "localhost"
					}
					node["status"].(map[string]string)[k] = "down"
				}
			}
		}

	} else {
		// Get node status when cluster started
		lines := strings.Split(string(out), "\n")
		var ringIds []string
		for _, line := range lines {
			if strings.Contains(line, "LINK") && !strings.Contains(line, "disk") {
				parts := strings.Split(strings.TrimSpace(line), " ")
				ringIds = append(ringIds, parts[1])
			}
		}
		for _, node := range nodesInfo {
			node["status"] = make(map[string]string, 0)
			for _, id := range ringIds {
				ringIdKey := "ring" + id + "_addr"
				// current node
				if currentNode == node["name"] {
					if _, exists := node[ringIdKey]; exists {
						node["status"].(map[string]string)[ringIdKey] = "localhost"
					}
					continue
				}
				// other node in clusters
				outStr := string(out)
				if ringIp, exists := node[ringIdKey]; exists {
					remoteKeyStr := "->" + ringIp.(string) + ") enabled connected"
					localKeyStr := ringIp.(string) + "->"
					if strings.Contains(outStr, remoteKeyStr) {
						node["status"].(map[string]string)[ringIdKey] = "up"
					} else if strings.Contains(outStr, localKeyStr) {
						node["status"].(map[string]string)[ringIdKey] = "localhost"
					} else {
						node["status"].(map[string]string)[ringIdKey] = "down"
					}
				}
			}
		}
	}
	return nil
}

// clusterSetup sets up a cluster with the provided node information.
// Returns results indicating the success or failure of the cluster setup.
func clusterSetup(cluster ClusterData) map[string]interface{} {
	clusterName := cluster.Cluster_name
	if clusterName == "" {
		clusterName = settings.DefaultClusterName
	}
	// TODO: --force
	nodeCmdStr := generateNodeCmdStr(cluster.Data)
	cmd := fmt.Sprintf(utils.CmdSetupClusterStandard, clusterName, nodeCmdStr)
	output, err := utils.RunCommand(cmd)
	outputStr := string(output[:])
	if err != nil {
		return map[string]interface{}{"action": false, "error": gettext.Gettext("Create cluster failed"), "detailInfo": outputStr}
	}

	_, err = utils.RunCommand(utils.CmdCreateAlert)
	if err != nil {
		return map[string]interface{}{"action": true, "info": gettext.Gettext("Create cluster success"),
			"alertInfo": gettext.Gettext("Failed to configure the alarm function module. If you need an alarm log, please manually execute the following command: pcs alert create id=alert_log path=/usr/share/pacemaker/alerts/alert_log.sh")}
	}

	return map[string]interface{}{"action": true, "info": gettext.Gettext("Create cluster success"), "alertInfo": gettext.Gettext("Create alert_log success")}
}

// generateNodeCmdStr generates the command string for adding nodes to the cluster.
// Returns the generated command string.
func generateNodeCmdStr(nodesInfo []NodeData) string {
	hbIPPrefix := "addr="
	var cmd strings.Builder
	hbIPCmd := ""
	for _, nodeInfo := range nodesInfo {
		//nodeInfoV := nodeInfo.(map[string]interface{})
		nodeStr := fmt.Sprintf("%v", nodeInfo.Name)
		for _, v := range nodeInfo.RingAddr {
			hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v.Ip)
			nodeStr += hbIPCmd
		}
		cmd.WriteString(" " + nodeStr)
	}
	return cmd.String()
}

func LocalClustersDestroy() map[string]interface{} {
	res := map[string]interface{}{}
	cmd := utils.CmdDestroyCluster
	out, err := utils.RunCommand(cmd)
	if err != nil {
		res["action"] = false
		res["error"] = string(out)
		return res
	}
	res["action"] = true
	res["message"] = string(out)
	return res
}

// isIPv4 checks if the provided string is a valid IPv4 address.
// Returns true if the string is a valid IPv4 address, false otherwise.
func isIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return false
		}
		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

func LocalAddNodes(addNodes AddNodesData) interface{} {
	addNodesInfo := addNodes.Data
	nodeList := make([]string, 0)
	password := make([]string, 0)
	ip := make([]string, 0)
	var authInfo AuthInfo
	for _, node := range addNodesInfo {
		nodeList = append(nodeList, node.Name)
		password = append(password, node.Password)
		for _, v := range node.RingAddr {
			ip = append(ip, v.Ip)
		}
	}
	authInfo.nodeList = nodeList
	authInfo.passWord = password
	authInfo.ip = ip
	authres := hostAuthWithAddr(authInfo)
	if !authres.Action {
		return authres
	}

	if IsClusterExist() {
		hbIPPrefix := "addr="
		addNodeCmd := ""

		currentNodeData, _ := utils.RunCommand(utils.CmdHostName)
		currentNode := string(currentNodeData)
		currentNode = strings.Replace(currentNode, "\n", "", -1)

		cmd := fmt.Sprintf("echo \"`pcs stonith sbd status`\"| grep %s:", currentNode)
		out, _ := utils.RunCommand(cmd)
		curNodeSbdStat := strings.Split(string(out), ":")[1]
		curNodeRunSbd := strings.Split(curNodeSbdStat, "|")
		if curNodeRunSbd[1] == " YES " {
			out, _ := utils.RunCommand(utils.CmdGetSbdStatus)
			sbdHeader := strings.Split(string(out), "SBD header on device")
			deviceInfo := strings.Split(sbdHeader[1], ":")
			sbdDevice := strings.TrimSpace(deviceInfo[0])
			for _, nodeInfo := range addNodesInfo {
				nodeStr := nodeInfo.Name
				for _, v := range nodeInfo.RingAddr {
					hbIPCmd := ""
					hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v.Ip)
					nodeStr += hbIPCmd
				}
				addNodeCmd = fmt.Sprintf(utils.CmdNodeAddStart, nodeStr) + "device=" + sbdDevice
			}
		} else {
			for _, nodeInfo := range addNodesInfo {
				nodeStr := nodeInfo.Name
				for _, v := range nodeInfo.RingAddr {
					hbIPCmd := ""
					hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v.Ip)
					nodeStr += hbIPCmd
				}
				addNodeCmd = fmt.Sprintf(utils.CmdNodeAddStart, nodeStr)
			}
		}
		out, err := utils.RunCommand(addNodeCmd)
		if err != nil {
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("Add node failed"),
				"detailInfo": string(out),
			}
		}

	} else {
		var clusterInfo ClusterData
		clusterInfo.Cluster_name = addNodes.Cluster_name
		clusterInfo.Data = addNodesInfo
		return clusterSetup(clusterInfo)
	}
	return map[string]interface{}{
		"action":  true,
		"message": "Add node success",
	}
}
