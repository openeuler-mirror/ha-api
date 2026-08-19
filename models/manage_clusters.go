/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/beevik/etree"
	"github.com/pkg/errors"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

// ClustersInfo is a structure representing information about clusters.
type ClustersInfo struct {
	Text     map[string]interface{}
	Version  int
	Clusters []Cluster
}

type Cluster struct {
	ClusterName string                   `json:"cluster_name"`
	Nodes       []string                 `json:"nodes"`
	Nodeid      []string                 `json:"nodeid"`
	Ip          []map[string]interface{} `json:"ip"`
}

// 集群添加接口
type ClusterData struct {
	Cluster_name string `json:"cluster_name"`
	Data         []NodeData
}

// 节点添加接口数据
type NodeData struct {
	Type         string   `json:"type,omitempty"`
	NodeID       int      `json:"nodeid,omitempty"`
	Name         string   `json:"name"`
	Password     string   `json:"password,omitempty"`
	RingAddr     []string `json:"ring_data,omitempty"`
	ResourceName string   `json:"resource_name,omitempty"` // for remote/guest node
	// RingAddr     []RingAddrData `json:"ring_addr,omitempty"`
}

// type RingAddrData struct {
// 	Ring string `json:"ring"`
// 	Ip   string `json:"ip"`
// }

type RemoveData struct {
	Cluster_name []string
}

type RemoveRet struct {
	Action        bool     `json:"action"`
	Error         string   `json:"error,omitempty"`
	FailedCluster []string `json:"faild_cluster"`
	Data          []bool   `json:"data"`
}

type AddNodesData struct {
	Cluster_name string     `json:"cluster_name"`
	Data         []NodeData `json:"data"`
}

type DeleteNodesData struct {
	Cluster_name string     `json:"cluster_name"`
	Data         []NodeData `json:"data"`
}

type AuthRetA struct {
	Action     bool   `json:"action"`
	Error      string `json:"error,omitempty"`
	DetailInfo string `json:"detailInfo,omitempty"`
	Message    string `json:"message,omitempty"`
}

// NewClustersInfo creates a new ClustersInfo instance using the provided text data.
// If the text data is nil or empty, default values are initialized.
func NewClustersInfo(text map[string]interface{}) *ClustersInfo {
	c := &ClustersInfo{
		Text: text,
	}

	if len(text) == 0 {
		c.Text = make(map[string]interface{})
		c.Version = 0
		c.Clusters = make([]Cluster, 0)
		c.Text["version"] = c.Version
		c.Text["clusters"] = c.Clusters
	} else {
		c.Version = int(text["version"].(float64))
		clustersInterface, ok := text["clusters"].([]interface{})
		if !ok {
			slog.Error("clusters is not a slice of interface{}") // 空时进入；
		}
		for _, clusterInterface := range clustersInterface {
			clusterMap, ok := clusterInterface.(map[string]interface{})
			if !ok {
				continue
			}

			cluster, err := MapToCluster(clusterMap)
			if err != nil {
				continue
			}
			c.Clusters = append(c.Clusters, cluster)
		}
	}

	return c
}

// mapToStruct 将map转换为指定的结构体
func MapToCluster(m map[string]interface{}) (Cluster, error) {
	// 将map转换为JSON字符串
	bytes, err := json.Marshal(m)
	if err != nil {
		return Cluster{}, err
	}

	// 将JSON字符串解码到Cluster结构体
	var cluster Cluster
	err = json.Unmarshal(bytes, &cluster)
	if err != nil {
		slog.Error(fmt.Sprintf("json.Unmarshal failed: %s", err.Error()))
		return Cluster{}, err
	}

	return cluster, nil
}

// Save updates the version, performs a backup, and saves the ClustersInfo to a file in JSON format.
func (ci *ClustersInfo) Save() error {
	ci.Version++
	ci.Backup()
	saveConf := ci.UpdateText()
	file, err := os.Create(settings.ClustersConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(saveConf); err != nil {
		return err
	}
	return nil
}

// Backup creates a backup of the cluster information file with a timestamp.
func (ci *ClustersInfo) Backup() error {
	// Implement backup functionality
	// You can use the time and date in the file name
	cureTime := time.Now().Unix()
	backFile := fmt.Sprintf("%s.%d", settings.ClustersConfigFile, cureTime)
	backCount, err := BackCount(ci)
	if err == nil && backCount < settings.MaxBackTimes {
		os.Rename(settings.ClustersConfigFile, backFile)
		return nil
	}
	return err
}

func BackCount(ci *ClustersInfo) (int, error) {
	if out, err := utils.RunCommand(utils.CmdCountClustersConfigsBackuped); err != nil {
		return 0, err
	} else {
		return strconv.Atoi(strings.TrimSpace(string(out))), nil
	}
}

// UpdateText updates the version and clusters in the Text field and returns it.
func (ci *ClustersInfo) UpdateText() map[string]interface{} {
	ci.Text["version"] = ci.Version
	ci.Text["clusters"] = ci.Clusters
	return ci.Text
}

// AddCluster adds cluster information to the Clusters field.
func (ci *ClustersInfo) AddCluster(clusterInfo Cluster) {
	ci.Clusters = append(ci.Clusters, clusterInfo)
}

// IsClusterNameInUse checks if a cluster name is already in use.
func (ci *ClustersInfo) IsClusterNameInUse(clusterName string) bool {
	for _, c := range ci.Clusters {
		if c.ClusterName == clusterName {
			return true
		}
	}
	return false
}

// SetVersion sets the version of the ClustersInfo.
func (ci *ClustersInfo) SetVersion(version int) {
	ci.Version = version
}

// DeleteCluster deletes the Cluster from ClustersInfo.
func (ci *ClustersInfo) DeleteCluster(clusterNameJson string) bool {
	for i, c := range ci.Clusters {
		if c.ClusterName == clusterNameJson {
			ci.Clusters = append(ci.Clusters[:i], ci.Clusters[i+1:]...)
			return true
		}
	}
	return false
}

func (ci *ClustersInfo) UpdateCluster(clusterNameJson string, clusterInfo Cluster) {
	for i, c := range ci.Clusters {
		if c.ClusterName == clusterNameJson {
			ci.Clusters[i].Nodes = clusterInfo.Nodes
			ci.Clusters[i].Nodeid = clusterInfo.Nodeid
			ci.Clusters[i].Ip = clusterInfo.Ip
		}
	}
}

// GetNodes  gets nodes information
func (ci *ClustersInfo) GetNodes(clusterNameJson string) []string {
	for _, c := range ci.Clusters {
		if c.ClusterName == clusterNameJson {
			return c.Nodes
		}
	}
	return []string{}
}

func (ci *ClustersInfo) GetClusterNameOfNode(nodeName string) string {
	for _, cluster := range ci.Clusters {
		nodes := cluster.Nodes
		for _, node := range nodes {
			if node == nodeName {
				return cluster.ClusterName
			}
		}
	}
	return ""
}

func ClusterInfo() map[string]interface{} {
	localConf := getLocalConf()
	clusterSum := len(localConf.Clusters)

	if clusterSum == 0 {
		return map[string]interface{}{
			"action":       false,
			"cluster_list": []interface{}{},
		}
	} else {
		return map[string]interface{}{
			"action":       false,
			"cluster_list": checkClusterExist(),
		}
	}
}

func ClusterOverview() map[string]interface{} {
	_ = ClusterInfo()
	clusterExist := false
	localClusterName := ""
	clusterExistInfo := CheckIsClusterExist()

	if clusterExistInfo["action"] == true {
		clusterExist = true
		localClusterName = clusterExistInfo["cluster_name"].(string)
	}
	localConf := getLocalConf()
	clusters := localConf.Clusters
	clusterSum := len(clusters)
	if clusterSum == 0 {
		return map[string]interface{}{
			"action":             true,
			"cluster_exist":      clusterExist,
			"local_cluster_name": localClusterName,
			"cluster_data":       []interface{}{},
		}
	}
	var (
		list []oneClusterOverviewRes
		mu   sync.Mutex
		wg   sync.WaitGroup
	)
	if len(localConf.Clusters) > 0 {

		for _, cluster := range localConf.Clusters {
			ips := extractIPs(cluster)

			wg.Add(1)
			go func(cluster Cluster) {
				// checkOneClusterExist 内部通过 defer wg.Done() 负责计数器递减，
 	 			// 此处不能再调用 wg.Done()，否则会造成 WaitGroup 计数为负触发 panic。
				res := oneClusterOverview(cluster, localConf, ips, &wg)
				// 处理nil问题
				if res.NodeList == nil {
					res.NodeList = make([]Node, 0)
				}
				if res.ResourceList == nil {
					res.ResourceList = make([]Resource, 0)
				}
				mu.Lock()
				list = append(list, res)
				mu.Unlock()

			}(cluster)
		}
		wg.Wait()
	}
	return map[string]interface{}{
		"action":             true,
		"cluster_exist":      clusterExist,
		"local_cluster_name": localClusterName,
		"cluster_data":       list,
	}
}

func extractIPs(clusters Cluster) []IP {
	var ips []IP
	for _, ipEntry := range clusters.Ip {
		newIP := IP{Addrs: make(map[string]string)}
		for key, value := range ipEntry {
			if strings.HasPrefix(key, "ring") {
				ringNum := key[4 : len(key)-5]
				newIP.Addrs["ring"+ringNum] = value.(string)
			}
		}
		ips = append(ips, newIP)
	}
	return ips
}

func oneClusterOverview(cluster Cluster, localconf *ClustersInfo, ips []IP, wg *sync.WaitGroup) oneClusterOverviewRes {
	var singleClusterInfo oneClusterOverviewRes
	singleClusterInfo.ClusterName = cluster.ClusterName
	singleClusterInfo.NodeSum = len(cluster.Nodes)
	nodeList := cluster.Nodes
	connectNode := 0
	clusterConnect := false
	for _, node := range nodeList {
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/managec/local_cluster_overview")
		slog.Info("Starting to request local cluster overview", "node", node)
		resp, err := utils.SendRequest(url, "GET", nil)
		if err != nil {
			// 连接失败异常捕获部分
			slog.Warn("Failed to send request, possibly a network or connection issue", "node", node, "url", "managec/local_cluster_overview", "error", err.Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			clusterConnect = true
			connectNode = connectNode + 1
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Error reading response body", "node", node, "err", err.Error())
				continue
			}
			var resInfo localClusterOverviewRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				slog.Error("Error unmarshal response json", "node", node, "err", err.Error())
				continue
			}
			if resInfo.Action {
				if resInfo.ClusterStart {
					var oneClusterOverviewRes oneClusterOverviewRes
					oneClusterOverviewRes.ClusterName = resInfo.Data.ClusterName
					oneClusterOverviewRes.NodeList = resInfo.Data.NodeList
					oneClusterOverviewRes.NodeSum = resInfo.Data.NodeSum
					oneClusterOverviewRes.ResourceList = resInfo.Data.ResourceList
					oneClusterOverviewRes.ClusterOnline = resInfo.Data.ClusterOnline
					oneClusterOverviewRes.Ip = ips
					return oneClusterOverviewRes
				} else {
					connectNode = connectNode - 1
				}
			}
		}
	}

	if connectNode == 0 {
		slog.Warn("No active and started cluster found on any node. Preparing fallback cluster info.")
		var singleNodeInfo Node
		for _, node := range nodeList {
			singleNodeInfo.Name = node
			singleNodeInfo.Online = "false"
			singleClusterInfo.NodeList = append(singleClusterInfo.NodeList, singleNodeInfo)
		}
		if !clusterConnect {
			singleClusterInfo.ClusterOnline = "false"
		} else {
			singleClusterInfo.ClusterOnline = "stop"
		}
		singleClusterInfo.Ip = ips
		var EmptyResourceList []Resource
		singleClusterInfo.ResourceList = EmptyResourceList
		return singleClusterInfo
	}
	return singleClusterInfo
}

func LocalClusterOverview() localClusterOverviewRes {
	var result localClusterOverviewRes
	var data localClusterOverviewData
	data.ClusterOnline = "false"
	if !IsClusterExist() {
		result.Action = false
		result.ClusterStart = false
		result.Data = data
		return result
	}
	data.ClusterName = getClusterName()
	nodeStatus := nodeStatus()
	if !nodeStatus["action"].(bool) {
		result.Action = true
		result.ClusterStart = false
		result.Data = data
		return result
	} else {
		data.NodeSum = len(nodeStatus["data"].([]Node))
		data.NodeList = nodeStatus["data"].([]Node)
		for _, node := range nodeStatus["data"].([]Node) {
			if node.Online == "true" {
				data.ClusterOnline = "true"
			}
		}
		data.ResourceList = resourceStatus()
	}
	result.Action = true
	result.ClusterStart = true
	result.Data = data
	return result
}

func getResourceStatus() map[string]interface{} {
	result := make(map[string]interface{}) // 初始化 map
	clusterStatus := GetClusterStatus()
	if clusterStatus != 0 {
		result["action"] = true
		result["data"] = []string{}
		return result
	}

	resList := GetAllResourceStatusForNew()
	result["action"] = true
	result["data"] = resList
	return result
}

func ExtractAKey(data map[string]interface{}) string {
	for key, value := range data {
		if _, ok := value.(map[string]interface{}); ok { // 类型断言检测嵌套map
			return key
		}
	}
	return ""
}
func resourceStatus() []Resource {
	var resourceList []Resource

	resourceData := getResourceStatus()
	failedList := GetResourceFailedList()
	for _, curRes := range resourceData["data"].([]map[string]interface{}) {
		var json Resource
		json.ID = curRes["id"].(string)
		if subRscs, ok := curRes["subrscs"].(bool); ok && subRscs {
			json.Status = subResourceStatus(curRes, failedList)
		} else {
			status, ok := curRes["status"].(string)
			if ok && status != "Failed" {
				json.Status = status
			} else {
				runningNode, ok := curRes["running_node"].([]interface{})
				if !ok || len(runningNode) == 0 {
					json.Status = "Stop"
				} else {
					json.Status = "Running but failed"
				}
			}
		}
		resourceList = append(resourceList, json)
	}
	return resourceList

}

func hasSubResources(subrsc map[string]interface{}) bool {
	if _, ok := subrsc["subrscs"].([]interface{}); ok {
		return true
	}
	return false
}

func subResourceStatus(resource map[string]interface{}, failedList []string) string {
	status := resource["status"].(string)

	if status == "Running" {
		// subrscs := resource["subrscs"].([]map[string]interface{})
		if subResources, ok := resource["subrscs"].([]interface{}); ok && len(subResources) > 0 {
			return checkSubResources(subResources, failedList)
		}
	} else if status == "Stopped" {
		status = "Stopped"
	} else {
		status = "Not Running"
	}
	return status
}

func containsPrefix(list []string, prefix string) bool {
	for _, v := range list {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}
	return false
}

func containsFailedID(res map[string]interface{}, list []string) bool {
	idVal, _ := res["id"].(string)
	return containsPrefix(list, strings.Split(idVal, ":")[0])
}

func isFailedState(sub map[string]interface{}, status string, failedList []string) bool {
	// 失败条件三元组判断‌:ml-citation{ref="8" data="citationList"}
	return containsFailedID(sub, failedList) ||
		status == "Stopped" ||
		status == "RunningButFailed"
}

func checkSubResources(subs []interface{}, failedList []string) string {
	for _, item := range subs {
		subMap, _ := item.(map[string]interface{})
		subStatus := subResourceStatus(subMap, failedList)

		// 失败条件判断‌:ml-citation{ref="7,8" data="citationList"}
		if isFailedState(subMap, subStatus, failedList) {
			return "RunningButFailed"
		}
	}
	return "Running"
}

func nodeStatus() map[string]interface{} {
	out, err := utils.RunCommand(utils.CmdCrmMonAsXML)
	var nodeList []Node
	if err != nil {
		result := make(map[string]interface{})
		result["action"] = false
		result["error"] = "cluster is offline"
		return result
	} else {
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(out); err != nil {
			slog.Error("Failed to parse XML output from crm_mon", "err", err.Error())
			return map[string]interface{}{"action": false, "error": err.Error()}
		}

		for _, nodes := range doc.FindElements("/crm_mon/nodes") {
			for _, node := range nodes.FindElements("node") {
				var n Node
				n.Name = node.SelectAttr("name").Value
				n.Online = node.SelectAttr("online").Value
				nodeList = append(nodeList, n)
			}
		}
	}
	return map[string]interface{}{"action": true, "data": nodeList}
}

func checkClusterExist() []Cluster {
	localConf := getLocalConf()
	var wg sync.WaitGroup
	if len(localConf.Clusters) > 0 {
		for _, cluster := range localConf.Clusters {
			wg.Add(1)
			go func(cluster Cluster) {
				// checkOneClusterExist 内部通过 defer wg.Done() 负责计数器递减，
				// 此处不能再调用 wg.Done()，否则会造成 WaitGroup 计数为负触发 panic。
				checkOneClusterExist(localConf, cluster, &wg)
			}(cluster)
		}
		wg.Wait()
	}
	return localConf.Clusters
}

type checkClusterExistRes struct {
	Action      bool    `json:"action"`
	ClusterName string  `json:"cluster_name"`
	ClusterConf Cluster `json:"cluster_conf"`
}

type localClusterOverviewRes struct {
	Action       bool                     `json:"action"`
	ClusterStart bool                     `json:"cluster_start"`
	Data         localClusterOverviewData `json:"data"`
}

type Node struct {
	Name   string `json:"name"`
	Online string `json:"online"`
}

type Resource struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type localClusterOverviewData struct {
	ClusterName   string     `json:"cluster_name"`
	NodeSum       int        `json:"node_sum"`
	NodeList      []Node     `json:"node_list"`
	ResourceList  []Resource `json:"resource_list"`
	ClusterOnline string     `json:"cluster_online"`
}

type oneClusterOverviewRes struct {
	ClusterName   string     `json:"cluster_name"`
	NodeSum       int        `json:"node_sum"`
	NodeList      []Node     `json:"node_list"`
	ResourceList  []Resource `json:"resource_list"`
	ClusterOnline string     `json:"cluster_online"`
	Ip            []IP       `json:"ip"`
}

type IP struct {
	Addrs map[string]string `json:"-"`
}

func checkOneClusterExist(localConf *ClustersInfo, cluster Cluster) {
	// defer wg.Done()
	slog.Info("check cluster exist", "clusterName", cluster.ClusterName)
	connectNode := 0
	confNodeSum := len(cluster.Nodes)
	realNodeNum := 0
	var clusterConf Cluster
	for _, node := range cluster.Nodes {
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/managec/is_cluster_exist")
		resp, err := utils.SendRequest(url, "GET", nil)
		if err != nil {
			slog.Warn("check node exist, failed: http error", "node", node, "err", err.Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Info(fmt.Sprintf("Error reading response body: %v", err))
				continue
			}
			var resInfo checkClusterExistRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				slog.Info(fmt.Sprintf("Error Unmarshal response json: %v", err))
				continue
			}
			if resInfo.Action {
				if resInfo.ClusterName == cluster.ClusterName {
					connectNode++
					clusterConf = resInfo.ClusterConf
					realNodeNum = len(clusterConf.Nodes)
					slog.Info("check node exist, passed.", "node", node)
				} else {
					confNodeSum--
					slog.Info("check node exist, failed: inconsistent cluster name", "node", node)
				}

			} else {
				confNodeSum--
				slog.Info("check node exist, failed: cluster not exist", "node", node)
			}
		} else {
			slog.Error("check node exist, failed: request error", "node", node, "url", url, "status", resp.StatusCode)
		}
	}
	handleExistClusterConf(realNodeNum, confNodeSum, clusterConf, cluster, localConf, cluster.ClusterName)
}

var CheckIsClusterExist = func() map[string]interface{} {
	result := map[string]interface{}{}
	_, err := os.Stat(settings.CorosyncConfFile)
	if err == nil {
		cmd := "cat /etc/corosync/corosync.conf | grep \"cluster_name\""
		out, err := utils.RunCommand(cmd)
		var clusterName string
		if err != nil {
			result["action"] = false
			result["error"] = "Get cluster name failed"
			return result
		}
		clusterName = strings.Split(string(out), ": ")[1]
		clusterName = strings.ReplaceAll(clusterName, "\n", "")

		allInfo := GetClusterInfo()
		if allInfo["cluster_exist"] == true {
			clusterInfo := clusterInfoParse(allInfo)
			result["action"] = true
			result["cluster_name"] = clusterName
			result["cluster_conf"] = clusterInfo
			return result

		}
	}
	result["action"] = false
	return result
}

func handleExistClusterConf(realNodeNum, confNodeSum int, clusterConf Cluster, cluster Cluster, localConf *ClustersInfo, clusterName string) {
	if realNodeNum != 0 && realNodeNum >= confNodeSum {
		if !reflect.DeepEqual(cluster, clusterConf) || IsNotSet(clusterConf) {
			localConf.UpdateCluster(cluster.ClusterName, clusterConf)
			localConf.Save()
			syncClusterConfFile(localConf)
		}
	} else if confNodeSum == 0 {
		localConf.DeleteCluster(clusterName)
		localConf.Save()
		syncClusterConfFile(localConf)
	} else if IsNotSet(clusterConf) {
		clusterStatus := checkClusterStatus(clusterConf)
		if !clusterStatus {
			localConf.DeleteCluster(clusterName)
			localConf.Save()
			syncClusterConfFile(localConf)
		}
	}

}

func IsNotSet(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			return false
		}
	}
	return true
}
func checkClusterStatus(clusterConf Cluster) bool {
	nodeList := clusterConf.Nodes
	clusterName := clusterConf.ClusterName
	nodeSum := len(nodeList)
	clusterExist := true
	connectNode := 0
	for _, node := range nodeList {
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/managec/is_cluster_exist")
		resp, err := utils.SendRequest(url, "GET", nil)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Info(fmt.Sprintf("Error reading response body: %v", err))
				continue
			}
			var resInfo checkClusterExistRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				slog.Info(fmt.Sprintf("Error Unmarshal response json: %v", err))
				continue
			}
			connectNode = connectNode + 1
			if resInfo.Action {
				if resInfo.ClusterName != clusterName {
					clusterExist = false
				}
			} else {
				clusterExist = false
			}
		}
	}
	if connectNode == nodeSum && !clusterExist {
		return false
	}
	return true
}

// localClusterInfo retrieves the corosync cluster information locally and returns it as a map.
// If no cluster exists, an empty map is returned.
func LocalClusterInfo() Cluster {
	allInfo := GetClusterInfo()
	if allInfo["cluster_exist"] == true {
		clusterInfo := clusterInfoParse(allInfo)
		return clusterInfo
	}
	var EmptyCluster Cluster
	return EmptyCluster
}

// clusterInfoParse takes cluster information as input and parses it into a map of string to interface
func clusterInfoParse(clusterInfo map[string]interface{}) Cluster {
	var clusterParse Cluster
	if clusterName, ok := clusterInfo["cluster_name"].(string); ok {
		clusterParse.ClusterName = clusterName
	}

	nodes := make([]string, 0)
	nodeIDs := make([]string, 0)
	ips := make([]map[string]interface{}, 0)
	nodesInfo := clusterInfo["data"].([]map[string]interface{})
	for _, nodeInfo := range nodesInfo {
		ip := make(map[string]interface{})
		for k, v := range nodeInfo {
			if k == "name" {
				nodes = append(nodes, v.(string))
			} else if k == "nodeid" {
				nodeIDs = append(nodeIDs, v.(string))
			} else {
				ip[k] = v
			}
		}
		ips = append(ips, ip)
	}

	clusterParse.Nodes = nodes
	clusterParse.Nodeid = nodeIDs
	clusterParse.Ip = ips
	return clusterParse
}

func GetLocalConf() *ClustersInfo {
	return getLocalConf()
}

// getLocalConf reads the local cluster configuration from a file and returns a ClustersInfo instance.
func getLocalConf() *ClustersInfo {
	localConf, _ := readFile(settings.ClustersConfigFile)
	return NewClustersInfo(localConf)
}

func GetRemoteCluster(node string) (remoteCluster *Cluster, err error) {
	url := utils.GenerateRemoteRequestURL(node, "/api/v1/managec/local_cluster_info")
	httpResp, err := utils.SendRequest(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("can not request the remote cluster info in %s : %w", node, err)
	}
	httpRespData, _ := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	var remoteClusterInfo Cluster
	if err := json.Unmarshal(httpRespData, &remoteClusterInfo); err != nil {
		return nil, fmt.Errorf("parse remote cluster info failed: %w", err)
	}
	return &remoteClusterInfo, nil
}

var getRemoteNodes = func(clusterName string) interface{} {
	localConf := getLocalConf()
	nodeList := localConf.GetNodes(clusterName)
	return nodeList
}

// readFile reads a JSON file, decodes its content, and returns it as a map.
func readFile(filename string) (map[string]interface{}, error) {
	var newDict map[string]interface{}

	infile, err := os.Open(filename)
	if err != nil {
		slog.Error("Error opening file:", "file", filename, "err", err.Error())
		return newDict, err
	}
	defer infile.Close()

	data, err := io.ReadAll(infile)
	if err != nil {
		slog.Error("Error reading file:", "file", filename, "err", err.Error())
		return newDict, err
	}

	if err := json.Unmarshal(data, &newDict); err != nil {
		slog.Error("Error decoding JSON:", "file", filename, "err", err.Error())
		return newDict, err
	}

	return newDict, nil
}

// comment out due to type error as localconf could not be {}, it should be of type *ClustersInfo
// SyncConfig synchronizes the local configuration with remote configuration.
// Returns appropriate results indicating the synchronization status.
func SyncConfig(remoteConf map[string]interface{}) map[string]interface{} {
	localConf := getLocalConf()
	remoteClusterInfo := NewClustersInfo(remoteConf)
	if remoteClusterInfo.Version >= localConf.Version { //|| localConf == ({}) {
		remoteClusterInfo.SetVersion(remoteClusterInfo.Version - 1)
		remoteClusterInfo.Save()
		return map[string]interface{}{
			"result": "receive",
			"conf":   remoteClusterInfo.Text,
		}
	} else {
		return map[string]interface{}{
			"result": "refuse",
		}
	}
}

// syncClusterConfFile synchronizes the cluster configuration file with all nodes in the cluster.
func syncClusterConfFile(conf *ClustersInfo) {
	// Get local cluster info
	clusterInfo := LocalClusterInfo()

	// If the current node has no cluster config, save the provided config
	if clusterInfo.ClusterName == "" {
		conf.Save()
		return
	}

	// Sync config file with all nodes in the cluster
	nodeList := clusterInfo.Nodes
	var failedNodes []string
	for _, node := range nodeList {
		// Node-to-node config file sync operation
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/sync_config")
		confJSON, err := json.Marshal(conf.Text)
		if err != nil {
			slog.Error("Failed to marshal config to JSON", "err", err.Error())
			return
		}

		resp, err := utils.SendRequest(url, "POST", confJSON)
		if err != nil {
			failedNodes = append(failedNodes, node)
		} else {
			resp.Body.Close()
		}

	}
	if len(failedNodes) != 0 {
		slog.Warn("Sync config to some nodes failed", "failedNodes", failedNodes)
	}
}

// hostAuth performs host authentication using the provided information.
func hostAuth(authInfo AuthRequest) utils.GeneralResponse {
	nodeList := authInfo.NodeList
	passwordList := authInfo.Passwords
	for i := 0; i < len(nodeList); i++ {
		authCmd := fmt.Sprintf(utils.CmdHostAuthNode, nodeList[i], passwordList[i])
		_, err := utils.RunCommand(authCmd)
		if err != nil {
			return utils.GeneralResponse{
				Action: false,
				Error:  fmt.Sprintf(gettext.Gettext("%s host auth failed. Username or Password incorrect"), nodeList[i]),
			}
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Error:  gettext.Gettext("host auth success"),
	}
}

func hostAuthWithAddr(authInfo AuthInfo) AuthRetA {
	authFailed := false
	authFailedInfo := ""

	authCmd := fmt.Sprintf(utils.CmdHostAuthNodeWithAddr, authInfo.nodeList[0], authInfo.ip[0], authInfo.passWord[0])
	out, err := utils.RunCommand(authCmd)
	if err != nil {
		authFailed = true
		authFailedInfo = string(out)
	}
	if authFailed {
		return AuthRetA{
			Action:     false,
			Error:      gettext.Gettext("host auth failed"),
			DetailInfo: authFailedInfo,
		}
	}
	return AuthRetA{
		Action:  true,
		Message: gettext.Gettext("host auth success"),
	}
}

type ClusterAddReq struct {
	NodeName string `json:"node_name"`
	PassWord string `json:"password"`
}

// ClusterAdd adds a new cluster using the provided node information.
// Returns results indicating the success or failure of the operation.
func ClusterAdd(nodeInfo ClusterAddReq) utils.GeneralResponse {
	// authInfo := make(map[string]interface{})
	// nodeList := make([]string, 0)
	// passwords := make([]string, 0)
	authInfo := AuthRequest{
		NodeList:  make([]string, 0, 1),
		Passwords: make([]string, 0, 1),
	}
	count, _, _ := utils.GrepHostsFile(nodeInfo.NodeName)
	if count == 0 {
		slog.Error("node authentication failed, incorrect node or node not in /etc/hosts file", "node", nodeInfo.NodeName)
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Incorrect name, node authentication failed"),
		}
	}

	if count > 1 {
		slog.Error("multiple entries found for node in /etc/hosts", "node", nodeInfo.NodeName, "count", count)
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Multiple entries found for this node in /etc/hosts, please ensure it is unique"),
		}
	}

	authInfo.NodeList = append(authInfo.NodeList, nodeInfo.NodeName)
	authInfo.Passwords = append(authInfo.Passwords, nodeInfo.PassWord)

	authRes := hostAuth(authInfo)

	if !authRes.Action {
		return handleAuthError(authRes.Error)
	}

	url := utils.GenerateRemoteRequestURL(nodeInfo.NodeName, "/api/v1/managec/local_cluster_info")
	resp, err := utils.SendRequest(url, "GET", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Cluster Add failed: %s", err.Error()))
		if resp != nil {
			defer resp.Body.Close()
		}
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("add cluster failed"),
		}

	}
	if resp != nil {
		defer resp.Body.Close()
	}

	var NewClusterInfo Cluster
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &NewClusterInfo)
	if err != nil {
		slog.Error(fmt.Sprintf("Cluster Add failed: %s", err.Error()))
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("add cluster failed"),
		}
	}
	localConf := getLocalConf()

	if localConf.IsClusterNameInUse(NewClusterInfo.ClusterName) {
		slog.Warn(fmt.Sprintf("Cluster Add failed: %s already exists, please check the ClustersInfo or remote cluster name", NewClusterInfo.ClusterName))
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("The cluster already exists, please do not add it again"),
		}
	}

	localConf.AddCluster(NewClusterInfo)
	localConf.Save()
	syncClusterConfFile(localConf)
	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("add cluster success"),
	}
}

func handleAuthError(errStr string) utils.GeneralResponse {
	errInfo := strings.ToLower(errStr)
	var errorMsg string

	if strings.Contains(errInfo, "username and/or password is incorrect") {
		errorMsg = gettext.Gettext("Incorrect password, node authentication failed")
	} else if strings.Contains(errInfo, "unable to synchronize and save known-hosts on nodes") {
		errorMsg = gettext.Gettext("Node host exception, unable to synchronize and save known-hosts")
	} else {
		errorMsg = gettext.Gettext("Node communication exception, unable to continue adding")
	}

	slog.Error(fmt.Sprintf("Cluster Add failed: %s", errStr))
	return utils.GeneralResponse{
		Action: false,
		Error:  errorMsg,
	}
}

func getAuthInfoFromClusterData(clusterSetInfo ClusterData) map[string]interface{} {
	authInfo := make(map[string]interface{})
	nodeList := make([]string, 0)
	passwords := make([]string, 0)

	nodesData := clusterSetInfo.Data
	for _, node := range nodesData {
		nodeList = append(nodeList, node.Name)
		passwords = append(passwords, node.Password)
	}
	authInfo["node_list"] = nodeList
	authInfo["password"] = passwords
	return authInfo
}

func getNodeListFromClusterData(clusterSetInfo ClusterData) []string {
	nodeList := make([]string, 0)
	nodesData := clusterSetInfo.Data

	for _, node := range nodesData {
		nodeList = append(nodeList, node.Name)
	}
	return nodeList
}

func ClusterSetup(clusterSetInfo ClusterData) map[string]interface{} {
	localClusters := getLocalConf()
	if localClusters.IsClusterNameInUse(clusterSetInfo.Cluster_name) {
		return map[string]interface{}{
			"action":     false,
			"error":      gettext.Gettext("ClusterName has been used"),
			"detailInfo": gettext.Gettext("ClusterName has been used")}
	}
	nodeList := getNodeListFromClusterData(clusterSetInfo)

	for _, node := range nodeList {
		slog.Info(fmt.Sprintf("Send the request of setup_cluster, node: %s", node))
		httpResp, err := setupInNode(node, clusterSetInfo)
		if err != nil {
			slog.Error(fmt.Sprintf("setup cluster in %s failed: %v", node, err))
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("Cluster cannot connect"),
				"detailInfo": gettext.Gettext("Cluster cannot connect")}
		}

		defer httpResp.Body.Close()
		httpRespData, _ := io.ReadAll(httpResp.Body)
		var httpRespJson map[string]interface{}
		json.Unmarshal(httpRespData, &httpRespJson)
		if !httpRespJson["action"].(bool) {
			// TODO: continue or return?
			// return map[string]interface{}{
			// 	"action":     false,
			// 	"error":      httpRespJson["error"].(bool),
			// 	"detailInfo": gettext.Gettext("Node cannot connect")}
			return httpRespJson
		}

		// 创建集群成功，将集群信息更新到配置文件中
		slog.Info(fmt.Sprintf("Create Cluster in %s success, sync the ClustersInfo file", node))
		remoteCluster, err := GetRemoteCluster(node)
		if err != nil {
			slog.Error(err.Error())
			return map[string]interface{}{
				"action":     false,
				"error":      gettext.Gettext("sync cluster info failed"),
				"detailInfo": gettext.Gettext("sync cluster info failed")}
		}

		localConf := getLocalConf()
		localConf.AddCluster(*remoteCluster)
		localConf.Save()
		syncClusterConfFile(localConf)
		return httpRespJson

	}
	return map[string]interface{}{
		"action":     false,
		"error":      gettext.Gettext("Cluster cannot connect"),
		"detailInfo": gettext.Gettext("Cluster cannot connect")}
}

func setupInNode(node string, data ClusterData) (*http.Response, error) {
	url := utils.GenerateRemoteRequestURL(node, "/remote/api/v1/setup_cluster")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json Marshal failed: %w", err)
	}
	resp, err := utils.SendRequest(url, "POST", jsonData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ClusterRemove(RemoveInfo RemoveData) *RemoveRet {
	clusters := RemoveInfo.Cluster_name
	localConf := getLocalConf()
	removeRes := make([]bool, 0)
	failedCluster := make([]string, 0)
	for _, cluster := range clusters {
		res := localConf.DeleteCluster(cluster)
		removeRes = append(removeRes, res)
		if !res {
			failedCluster = append(failedCluster, cluster)
		}
		localConf.Save()
		syncClusterConfFile(localConf)
	}
	var RetData RemoveRet
	RetData.Action = true
	RetData.FailedCluster = failedCluster
	RetData.Data = removeRes
	return &RetData
}


// 获取集群节点列表
func getClusterNodes(clusters []Cluster, name string) []string {
	for _, cluster := range clusters {
		if name == cluster.ClusterName {
			return cluster.Nodes
		}
	}
	return nil
}

// 尝试销毁集群
func tryDestroyCluster(nodes []string) (bool, string) {
	var (
		success    bool
		detailInfo string
	)

	for _, node := range nodes {
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/destroy_cluster")
		resp, err := utils.SendRequest(url, "GET", nil)
		if err != nil {
			detailInfo = (fmt.Sprintf(gettext.Gettext("The node %s cannot be connected, failed to destroy the cluster"), node))
			success = false
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			detailInfo = err.Error()
			continue
		}

		if result["action"].(bool) {
			success = true
			break
		} else {
			success = false
			detailInfo = result["error"].(string)
		}
	}

	return success, detailInfo
}

// 记录失败信息
func recordFailure(res *[]bool, clusterName string, failedClusters *[]string, details *[]string, msg string) {
	// (*res)[index] = false
	(*res) = append(*res, false)
	*failedClusters = append(*failedClusters, clusterName)
	*details = append(*details, msg)
}

func ClusterDestroy(clustersJSON map[string]interface{}) map[string]interface{} {
	// 1. 获取集群信息
	localConf := getLocalConf()
	clusters := clustersJSON["cluster_name"].([]interface{})

	type result struct {
		success     bool
		clusterName string
		detailInfo  string
	}
	resultChan := make(chan result, len(clusters))

	// 由于集群摧毁可能执行事件比较长， 所以启动goroutine并行处理每个集群
	for _, desCluster := range clusters {
		desClusterName := desCluster.(string)
		go func(name string) {
			nodes := getRemoteNodes(desClusterName).([]string)
			if len(nodes) == 0 {
				resultChan <- result{
					success:     false,
					clusterName: name,
					detailInfo:  gettext.Gettext("Cluster not found"),
				}
				return
			}
			success, detailInfo := tryDestroyCluster(nodes)
			resultChan <- result{
				success:     success,
				clusterName: name,
				detailInfo:  detailInfo,
			}
		}(desClusterName)

	}

	// 收集结果进行返回
	res := make([]bool, 0)
	failedClusterList := make([]string, 0)
	detailInfos := make([]string, 0)
	for i := 0; i < len(clusters); i++ {
		r := <-resultChan
		res = append(res, r.success)
		if !r.success {
			failedClusterList = append(failedClusterList, r.clusterName)
			detailInfos = append(detailInfos, r.detailInfo)
		} else {
			// 成功则更新配置
			localConf.DeleteCluster(r.clusterName)
		}

	}

	localConf.Save()
	syncClusterConfFile(localConf)

	return map[string]interface{}{
		"action":     true,
		"data":       res,
		"clusters":   failedClusterList,
		"detailInfo": detailInfos,
	}
}

func ClusterDestroy2(clustersJSON map[string]interface{}) map[string]interface{} {
	localConf := getLocalConf()
	clusters := clustersJSON["cluster_name"].([]interface{})

	res := make([]bool, 0)
	failedClusterList := make([]string, 0)
	detailInfos := make([]string, 0)
	remoteUiPath := "/api/v1/destroy_cluster"
	// TODO: 异步执行
	for _, desCluster := range clusters {
		clusterName := desCluster.(string)
		nodes := getRemoteNodes(clusterName).([]string)
		if len(nodes) == 0 {
			// TODO: log
			recordFailure(&res, clusterName, &failedClusterList, &detailInfos, gettext.Gettext("Cluster not found"))
			continue
		}
		resultRemote, err := UrlRedirect(clusterName, remoteUiPath, "GET", nil, nil)
		if !resultRemote["action"].(bool) {
			if err == nil {
			recordFailure(&res, clusterName, &failedClusterList, &detailInfos, err.Error())
		        }else{
                        recordFailure(&res, clusterName, &failedClusterList, &detailInfos, resultRemote["error"])
			}	
			continue
		}

		// 删除集群
		localConf.DeleteCluster(desCluster.(string))

	}

	// 刷到本地配置文件
	localConf.Save()
	syncClusterConfFile(localConf)

	//TODO: always return true
	return map[string]interface{}{
		"action":     true,
		"data":       res,
		"clusters":   failedClusterList,
		"detailInfo": detailInfos,
	}
}

// UrlRedirect
func UrlRedirect(clusterName string, uiPath string, requestMethod string, requestData interface{}, postRequestHook func(node string, respData []byte) error) (map[string]interface{}, error) {
	remoteNodes := getRemoteNodes(clusterName).([]string)
	if len(remoteNodes) == 0 {
		return map[string]interface{}{
			"action": false,
			"error":  gettext.Gettext("Please reselect the cluster in the top operation area"),
		}, errors.New("no remote nodes")
	}

	for _, node := range remoteNodes {
		url := utils.GenerateRemoteRequestURL(node, uiPath)
		resp, err := utils.SendRequest(url, requestMethod, requestData)
		if err != nil {
			slog.Warn("request to node failed (retry next node)", "url", url, "node", node, "error", err)
			continue
		}
		respData, err := io.ReadAll(resp.Body)
		resp.Body.Close() // 立即关闭
		if err != nil {
			slog.Error("failed to read response body", "url", uiPath, "node", node, "error", err)
			continue
		}

		remoteClusterInfo := make(map[string]interface{})
		if err := json.Unmarshal(respData, &remoteClusterInfo); err != nil {
			slog.Error(fmt.Sprintf("parse response failed: %s", err.Error()))
			continue
		}

		action, ok := remoteClusterInfo["action"].(bool)
		if !ok {
			slog.Warn("invalid or missing 'action' field in response", "url", uiPath, "response", remoteClusterInfo)
			continue
		}

		// 远程集群中节点响应，但是执行失败（action为false）， 直接返回
		if !action {
			slog.Info("Remote cluster operation failed", "url", uiPath, "response", remoteClusterInfo)
			return remoteClusterInfo, fmt.Errorf("remote operation failed: %v", remoteClusterInfo["error"])
		}

		// 如果有差异逻辑，执行对应回调函数
		if postRequestHook != nil {
			if err := postRequestHook(node, respData); err != nil {
				slog.Error(fmt.Sprintf("postRequestHook failed : %s", err.Error()))
				continue
			}
		}

		return remoteClusterInfo, nil
	}
	slog.Error("UrlRedirect failed (all nodes )", "cluster", clusterName, "url", uiPath)
	return map[string]interface{}{"action": false, "error": gettext.Gettext("Please reselect the cluster in the top operation area")}, errors.New("no nodes succeeded")
}

func UrlRedirectWithSyncConfig(clusterName string, uiPath string, requestMethod string, requestData interface{}) (map[string]interface{}, error) {
	syncHook := func(node string, _ []byte) error {
		url := utils.GenerateRemoteRequestURL(node, "/api/v1/managec/local_cluster_info")
		httpResp, err := utils.SendRequest(url, "GET", nil)
		if err != nil {
			return fmt.Errorf("update config to ClustersInfo failed: %w", err)
		}
		httpRespData, _ := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		var remoteClusterInfo Cluster
		if err := json.Unmarshal(httpRespData, &remoteClusterInfo); err != nil {
			return fmt.Errorf("parse cluster info failed: %w", err)
		}
		err = UpdateClusterConfFile(remoteClusterInfo)
		if err != nil {
			return fmt.Errorf("update config to ClustersInfo failed: %w", err)
		}
		return nil
	}
	// 调用远程并且执行同步操作
	return UrlRedirect(clusterName, uiPath, requestMethod, requestData, syncHook)

}

// 更新并同步本地集群配置文件
func UpdateClusterConfFile(cluster Cluster) error {
	// localCluster := LocalClusterInfo()
	clusterName := cluster.ClusterName
	localClusters := getLocalConf()
	localClusters.UpdateCluster(clusterName, cluster)
	if err := localClusters.Save(); err != nil {
		return fmt.Errorf("update config to ClustersInfo failed: %w", err)
	}
	syncClusterConfFile(localClusters)
	return nil
}
