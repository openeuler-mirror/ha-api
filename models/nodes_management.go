/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
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

// 定义集群管理界面的Get返回值结构体
// {action:true, cluster_exist:true, cluster_name:"hacluster",currentNode:"ha1",
//
//	data:[name:"ha1",nameid:"0",ring_data:[{ringname:"ring0_addr"，ip:"", status:""},{ringname:"ring1_addr"，ip:"", status:""}]}
type NodeManageClusterInfo struct {
	Action        bool                        `json:"action"`
	Cluster_Exist bool                        `json:"cluster_exist"`
	Cluster_Name  string                      `json:"cluster_name"`
	CurrentNode   string                      `json:"currentNode,omitempty"`
	Data          []NodeManageClusterInfoData `json:"data,omitempty"`
	Error         string                      `json:"error,omitempty"`
}

type NodeManageClusterInfoData struct {
	Name     string     `json:"name"`
	NodeId   string     `json:"nodeid,omitempty"`
	Type     string     `json:"type"`
	RingData []RingData `json:"ring_data"`
}

type RingData struct {
	RingName string `json:"ring_name,omitempty"`
	IP       string `json:"ip,omitempty"`
	Status   string `json:"status,omitempty"`
}

func GetNodesList() ([]string, error) {
	// 读取文件
	filename := "/etc/corosync/corosync.conf"
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("read config from /etc/corosync/corosync.conf failed")
	}
	defer file.Close()

	var nodeList []string // 存储文件内数据
	scanner := bufio.NewScanner(file)

	// 跳过直到找到 "nodelist {"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "nodelist {") {
			break
		}
	}

	// 读取直到 "quorum {"
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "quorum {") {
			break
		}
		if line != "" {
			nodeList = append(nodeList, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 过滤不需要的行
	var result []string
	for _, line := range nodeList {
		if line != "nodelist {" && line != "}" && line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

func GetRemoteInfo() ([]NodeManageClusterInfoData, []NodeManageClusterInfoData, error) {
	//声明remote、guest节点
	RemoteNodes := []NodeManageClusterInfoData{}
	GuestNodes := []NodeManageClusterInfoData{}

	//执行crm_mon -1 --as-xml命令
	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	if err != nil {
		return RemoteNodes, GuestNodes, errors.New("crm_mon -1 --as-xml failed")
	}

	doc := etree.NewDocument()
	if err1 := doc.ReadFromBytes(out); err1 != nil {
		return RemoteNodes, GuestNodes, errors.New(gettext.Gettext("parse xml failed"))
	}

	nodes := doc.SelectElement("crm_mon").SelectElement("nodes")

	for _, node := range nodes.SelectElements("node") {
		remoteNode := NodeManageClusterInfoData{}
		if node.SelectAttr("type").Value == "remote" {
			if node.SelectAttr("id_as_resource") == nil {
				remoteNode.Name = node.SelectAttr("name").Value
				remoteNode.NodeId = node.SelectAttr("id").Value
				remoteNode.RingData = nil
				remoteNode.Type = "remote"
				RemoteNodes = append(RemoteNodes, remoteNode)
			} else {
				remoteNode.Name = node.SelectAttr("name").Value
				remoteNode.NodeId = node.SelectAttr("id").Value
				remoteNode.RingData = nil
				remoteNode.Type = "guest"
				GuestNodes = append(GuestNodes, remoteNode)
			}
		}
	}

	return RemoteNodes, GuestNodes, nil
}

// getClusterInfo retrieves cluster information, including cluster nodes and their properties.
// Returns the cluster information in a structured map.
// 将社区获取节点函数utils.GetNodeList()替代为GetNodesList()
func GetClusterInfo1() NodeManageClusterInfo {
	var ret NodeManageClusterInfo
	currentNode, _ := utils.RunCommand(utils.CmdHostName)
	currentNodeContent := strings.ReplaceAll(string(currentNode), "\n", "")
	clusterName := getClusterName()
	var errorInfo string
	var nodeList []string
	var err error
	var nodes []map[string]string
	var ringIDs []string

	if !IsClusterExist() {
		errorInfo = "Cluster not established!"
		ret = NodeManageClusterInfo{
			Action:        false,
			Cluster_Exist: false,
			Cluster_Name:  clusterName,
			CurrentNode:   currentNodeContent,
			Error:         errorInfo,
		}
		return ret
	}
	nodeList, _ = GetNodesList()
	//nodeList返回示例：['node {', 'ring0_addr: 10.44.66.144', 'name: host141', 'nodeid: 1', 'ring1_addr: 192.168.122.202',
	//'node {', 'ring0_addr: 10.44.66.145', 'name: host142', 'nodeid: 2', 'ring1_addr: 192.168.122.24']
	index := 0
	for index < len(nodeList) {
		if nodeList[index] == "node {" {
			//按照node { 进行分割遍历每个node
			index += 1
			nodeInfo := make(map[string]string)
			nodeInfoRet := make(map[string]string)

			//解析单个节点信息
			for index < len(nodeList) && nodeList[index] != "node {" {
				line := nodeList[index]
				if strings.Contains(line, "ring7_addr") && strings.Contains(line, "disk") {
					index++
					continue
				}

				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					nodeInfo[key] = value
				}
				index++
			}

			// 处理必要字段
			for _, field := range []string{"nodeid", "name"} {
				if value, exists := nodeInfo[field]; exists {
					nodeInfoRet[field] = value
					delete(nodeInfo, field)
				} else {
					// haLogger.Warning("Get cluster info failed, cluster was not established.")
					errorInfo = "Cluster not established!"
					ret = NodeManageClusterInfo{
						Action:        false,
						Cluster_Exist: false,
						Cluster_Name:  clusterName,
						CurrentNode:   currentNodeContent,
						Error:         errorInfo,
					}
					return ret
				}
			}

			// 处理其他字段
			for key, value := range nodeInfo {
				if key != "nodeid" && key != "name" {
					if isIPv4(value) {
						nodeInfoRet[key] = value
					} else {
						nodeInfoRet[key] = ""
					}
				}
			}

			nodeInfoRet["type"] = "primitive"
			nodes = append(nodes, nodeInfoRet)
		} else {
			index++
		}
	}
	// 获取心跳状态
	cmd := "corosync-cfgtool -n"
	output, err := utils.RunCommand(cmd)
	var res []NodeManageClusterInfoData
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "LINK:") && !strings.Contains(line, "disk") {
				parts := strings.Fields(line)
				ringIDs = append(ringIDs, parts[1])
			}
		}
		seen := make(map[string]bool)
		ringIDs_test := []string{}

		for _, value := range ringIDs {
			if _, exists := seen[value]; !exists {
				seen[value] = true
				ringIDs_test = append(ringIDs_test, value)
			}
		}
		ringIDs = ringIDs_test
		for _, node_info := range nodes {
			var data NodeManageClusterInfoData
			var ring_data RingData
			data.Name = node_info["name"]
			data.NodeId = node_info["nodeid"]
			data.Type = "primitive"
			if currentNodeContent == node_info["name"] {
				if len(ringIDs) == 0 {
					for key := range node_info {
						if strings.HasPrefix(key, "ring") {
							ring_data.RingName = key
							ring_data.IP = node_info[key]
							ring_data.Status = "localhost"
							data.RingData = append(data.RingData, ring_data)
						}
					}
				} else {
					for _, id := range ringIDs {
						key := "ring" + string(id) + "_addr"
						ring_data.RingName = key
						ring_data.IP = node_info[key]
						ring_data.Status = "localhost"
						data.RingData = append(data.RingData, ring_data)
					}
				}
			} else {
				if len(ringIDs) == 0 {
					for key := range node_info {
						if strings.HasPrefix(key, "ring") {
							ring_data.RingName = key
							ring_data.IP = node_info[key]
							ring_data.Status = "down"
							data.RingData = append(data.RingData, ring_data)
						}
					}
				}

				for _, id := range ringIDs {
					key := "ring" + string(id) + "_addr"
					ring_data.IP = node_info[key]
					ring_data.RingName = key
					remote_key_str := "->" + ring_data.IP + ") enabled connected"
					local_key_str := ring_data.IP + "->"
					if strings.Contains(string(output), remote_key_str) {
						ring_data.Status = "up"
					} else if strings.Contains(string(output), local_key_str) {
						ring_data.Status = "localhost"
					} else {
						ring_data.Status = "down"
					}
					data.RingData = append(data.RingData, ring_data)
				}
			}

			res = append(res, data)
		}
	} else {
		// 获取心跳编号列表
		var ringNames []string
		for key := range nodes[0] {
			if strings.HasPrefix(key, "ring") {
				ringNames = append(ringNames, key)
			}
		}

		// 按照心跳编号排序
		sort.Strings(ringNames)

		for _, node_info := range nodes {
			var data NodeManageClusterInfoData
			data.Name = node_info["name"]
			data.NodeId = node_info["nodeid"]
			data.Type = "primitive"

			var ringDataList []RingData
			// 基于有序心跳组织数据
			for _, ring := range ringNames {
				var ringData RingData
				ringData.RingName = ring
				ringData.IP = node_info[ring]
				if currentNodeContent == node_info["name"] {
					ringData.Status = "localhost"
				} else {
					ringData.Status = "down"
				}
				ringDataList = append(ringDataList, ringData)
			}

			data.RingData = ringDataList
			res = append(res, data)
		}
	}
	RemoteNodes, GuestNodes, err := GetRemoteInfo()
	if err == nil {
		if len(RemoteNodes) != 0 {
			for _, node := range RemoteNodes {
				res = append(res, node)
			}
		}
		if len(GuestNodes) != 0 {
			for _, node := range GuestNodes {
				res = append(res, node)
			}
		}
	}
	ret = NodeManageClusterInfo{
		Action:        true,
		Cluster_Exist: true,
		Cluster_Name:  clusterName,
		CurrentNode:   currentNodeContent,
		Data:          res,
	}
	return ret
}

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

type AuthRequest struct {
	NodeList  []string `json:"node_list"`
	Passwords []string `json:"password"`
}

type AuthResp struct {
	Action     bool     `json:"action"`
	FailedNode string   `json:"failed_node,omitempty"`
	Details    []string `json:"details,omitempty"`
}

func ClusterSetupPre(addNodes ClusterData) map[string]interface{} {
	// 初始化认证信息
	authInfo := AuthRequest{
		NodeList:  make([]string, 0, len(addNodes.Data)),
		Passwords: make([]string, 0, len(addNodes.Data)),
	}

	// 提取节点认证数据
	for _, node := range addNodes.Data {
		authInfo.NodeList = append(authInfo.NodeList, node.Name)
		authInfo.Passwords = append(authInfo.Passwords, node.Password)
	}

	// 执行主机认证
	authRes := hostAuth(authInfo)
	if !authRes.Action {
		return map[string]interface{}{
			"action": authRes.Action,
			"error":  authRes.Error,
		}
	}

	// 执行集群创建并返回结果
	return clusterSetup(addNodes)
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
			hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v)
			nodeStr += hbIPCmd
		}
		cmd.WriteString(" " + nodeStr)
	}
	return cmd.String()
}

func LocalClusterDestroy() map[string]interface{} {
	res := map[string]interface{}{}
	// 集群摧毁之前先检测节点中pcsd是否正常运行或者联通
	offLinesNodes, err := utils.CheckPcsdOfflineNodes()
	if err != nil || len(offLinesNodes) != 0 {
		res["action"] = false
		res["error"] = gettext.Gettext("There are nodes in the cluster where the pcsd service is running abnormally, and continuing to destroy the cluster poses risks")
		slog.Warn("There are nodes in the cluster where the pcsd service is running abnormally, and continuing to destroy the cluster poses risks. please check the pcsd status.")
		return res
	}

	cmd := utils.CmdDestroyClusterForce
	out, err := utils.RunCommand(cmd)
	if err != nil {
		res["action"] = false
		res["error"] = string(out)
		return res
	}
	res["action"] = true
	res["info"] = string(out)
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

func LocalAddNodes(addNodes AddNodesData) map[string]interface{} {
	addNodesInfo := addNodes.Data
	//获取添加的节点类型
	addNodesType := addNodesInfo[0].Type

	if addNodesType == "remote" {
		//认证操作
		nodeList := make([]string, 0)
		password := make([]string, 0)
		var authrequest AuthRequest
		for _, node := range addNodesInfo {
			nodeList = append(nodeList, node.Name)
			password = append(password, node.Password)
		}
		authrequest.NodeList = nodeList
		authrequest.Passwords = password
		hostres := hostAuth(authrequest)
		if !hostres.Action {
			// return hostres
			return map[string]interface{}{
				"action":     hostres.Action,
				"error":      gettext.Gettext("Add node failed"),
				"detailInfo": hostres.Error,
			}
		}

		//添加remote节点操作
		addnodecmd := "pcs cluster node add-remote " + nodeList[0]
		out, err := utils.RunCommand(addnodecmd)
		if err != nil {
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("Add node failed"),
				"detailInfo": string(out),
			}
		}
	}

	if addNodesType == "guest" {
		//认证操作
		nodeList := make([]string, 0)
		password := make([]string, 0)
		resourceName := make([]string, 0)
		var authrequest AuthRequest
		for _, node := range addNodesInfo {
			nodeList = append(nodeList, node.Name)
			password = append(password, node.Password)
			resourceName = append(resourceName, node.ResourceName)
		}
		authrequest.NodeList = nodeList
		authrequest.Passwords = password
		hostres := hostAuth(authrequest)
		if !hostres.Action {
			return map[string]interface{}{
				"action":     hostres.Action,
				"error":      gettext.Gettext("Add node failed"),
				"detailInfo": hostres.Error,
			}
		}

		//添加guest节点操作
		addnodecmd := "pcs cluster node add-guest " + nodeList[0] + " " + resourceName[0]
		out, err := utils.RunCommand(addnodecmd)
		if err != nil {
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("Add node failed"),
				"detailInfo": string(out),
			}
		}
	}

	if addNodesType == "primitive" {
		var authInfo AuthInfo
		nodeList := make([]string, 0)
		password := make([]string, 0)
		ip := make([]string, 0)
		// var authInfo AuthInfo
		for _, node := range addNodesInfo {
			nodeList = append(nodeList, node.Name)
			password = append(password, node.Password)
			ip = append(ip, node.RingAddr...)
		}
		authInfo.nodeList = nodeList
		authInfo.passWord = password
		authInfo.ip = ip

		authres := hostAuthWithAddr(authInfo)
		if !authres.Action {
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("host auth failed"),
				"detailInfo": gettext.Gettext("host auth failed"),
			}
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
						hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v)
						nodeStr += hbIPCmd
					}
					addNodeCmd = fmt.Sprintf(utils.CmdNodeAddStart, nodeStr) + "device=" + sbdDevice
				}
			} else {
				for _, nodeInfo := range addNodesInfo {
					nodeStr := nodeInfo.Name
					for _, v := range nodeInfo.RingAddr {
						hbIPCmd := ""
						hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v)
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
	}
	return map[string]interface{}{
		"action": true,
		"info":   gettext.Gettext("Add node success"),
	}
}
