/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Jason011125 <zic022@ucsd.edu>
 * Date: Mon Aug 14 15:53:52 2023 +0800
 */
package models

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/pkg/errors"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

var port, _ = utils.ReadPortFromConfig()

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
	Cluster_name string
	Data         []NodeData
}

// 节点添加接口数据
type NodeData struct {
	NodeID   int            `json:"nodeid"`
	Name     string         `json:"name"`
	Password string         `json:"password"`
	RingAddr []RingAddrData `json:"ring_addr"`
}

type RingAddrData struct {
	Ring string `json:"ring"`
	Ip   string `json:"ip"`
}

type RemoveData struct {
	Cluster_name []string
}

type RemoveRet struct {
	Action        bool     `json:"action,omitempty"`
	Error         string   `json:"error,omitempty"`
	Faild_cluster []string `json:"faild_cluster,omitempty"`
	Data          []bool   `json:"data,omitempty"`
}

type AddNodesData struct {
	Cluster_name string
	Data         []NodeData
}

type AddNodesRet struct {
	Action bool   `json:"action,omitempty"`
	Error  string `json:"error,omitempty"`
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
			logs.Error("clusters is not a slice of interface{}")
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
		logs.Error("json.Unmarshal failed: %s", err.Error())
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
		return int(binary.LittleEndian.Uint16(out)), nil
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
			fmt.Printf("Before assignment: c = %+v\n", c)
			ci.Clusters[i].Ip = clusterInfo.Ip
			fmt.Printf("After assignment: c = %+v\n", c)
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
			"action":             false,
			"cluster_exist":      clusterExist,
			"local_cluster_name": localClusterName,
			"cluster_data":       []interface{}{},
		}
	}
	var wg sync.WaitGroup
	if len(localConf.Clusters) > 0 {

		for _, cluster := range localConf.Clusters {
			// ips := []IP
			// for _, ipInfo := range cluster.Ip {
			// 	oneNodeIp := make(map[string]interface{})
			// 	for k, v := range ipInfo {
			// 		if k != "type" && k != "status" {
			// 			oneNodeIp[k] = v
			// 		}
			// 	}
			// 	ips = append(ips, oneNodeIp)
			// }
			// ips := []*IP{}
			// for _, ipInfo := range cluster.Ip {
			// 	ip := newIP()
			// 	for k, v := range ipInfo {
			// 		if strings.HasPrefix(k, "ring") && strings.HasSuffix(k, "_addr") {
			// 			ringNum := k[4 : len(k)-5]

			// 			if addr, ok := v.(string); ok {
			// 				ip.Addrs[ringNum] = addr
			// 			}
			// 		}
			// 	}
			// 	ips = append(ips, ip)
			// }
			ips := extractIPs(cluster)

			wg.Add(1)
			go func(cluster Cluster) {
				defer wg.Done()
				oneClusterOverview(cluster, localConf, ips, &wg)
			}(cluster)
		}
		wg.Wait()
	}
	return map[string]interface{}{
		"action":             true,
		"cluster_exist":      clusterExist,
		"local_cluster_name": localClusterName,
		"cluster_data":       []interface{}{},
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
		url := fmt.Sprintf(("https://%s/remote/api/v1/managec/local_cluster_overview"), node)
		resp, err := http.Get(url)
		if err != nil {
			// 连接失败异常捕获部分
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			clusterConnect = true
			connectNode = connectNode + 1
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logs.Info("Error reading response body: %v", err)
				continue
			}
			var resInfo localClusterOverviewRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				logs.Info("Error Unmarshal response json: %v", err)
				continue
			}
			if resInfo.Action {
				if resInfo.ClusterStart {
					var oneClusterOverviewRes oneClusterOverviewRes
					oneClusterOverviewRes.ClusterName = resInfo.Data.ClusterName
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
		var singleNodeInfo Node
		for _, node := range nodeList {
			singleNodeInfo.Name = node
			singleNodeInfo.Online = "false"
			singleClusterInfo.NodeList = append(singleClusterInfo.NodeList, singleNodeInfo)
		}
		if clusterConnect == false {
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

func checkClusterExist() []Cluster {
	localConf := getLocalConf()
	var wg sync.WaitGroup
	if len(localConf.Clusters) > 0 {
		for _, cluster := range localConf.Clusters {
			wg.Add(1)
			go func(cluster Cluster) {
				defer wg.Done()
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
	Addrs map[string]string `json: "-"`
}

func checkOneClusterExist(localConf *ClustersInfo, cluster Cluster, wg *sync.WaitGroup) {
	defer wg.Done()
	connectNode := 0
	confNodeSum := len(cluster.Nodes)
	realNodeNum := 0
	var clusterConf Cluster
	for _, node := range cluster.Nodes {
		url := fmt.Sprintf(("https://%s/remote/api/v1/managec/is_cluster_exist"), node)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logs.Info("Error reading response body: %v", err)
				continue
			}
			var resInfo checkClusterExistRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				logs.Info("Error Unmarshal response json: %v", err)
				continue
			}
			if resInfo.Action {
				if resInfo.ClusterName == cluster.ClusterName {
					connectNode++
					clusterConf = resInfo.ClusterConf
					realNodeNum = len(clusterConf.Nodes)
					logs.Info("Node %s information check passed.", node)

				} else {
					confNodeSum--
					logs.Info("Node %s information check failed: inconsistent cluster name.", node)
				}

			} else {
				confNodeSum--
				logs.Info("Node %s information check failed: cluster not exist", node)
			}
		} else {
			logs.Info("Get %s failed: status is %d.", url, resp.StatusCode)
		}
	}
	handleExistClusterConf(realNodeNum, confNodeSum, clusterConf, cluster, localConf, cluster.ClusterName)
}

func CheckIsClusterExist() map[string]interface{} {
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
		} else {
			clusterName = strings.Split(string(out), ": ")[1]
			clusterName = strings.ReplaceAll(clusterName, "\n", "")

		}

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
		if clusterStatus == false {
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
		url := fmt.Sprintf(("https://%s/remote/api/v1/managec/is_cluster_exist"), node)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logs.Info("Error reading response body: %v", err)
				continue
			}
			var resInfo checkClusterExistRes
			err = json.Unmarshal(body, &resInfo)
			if err != nil {
				logs.Info("Error Unmarshal response json: %v", err)
				continue
			}
			connectNode = connectNode + 1
			if resInfo.Action == true {
				if resInfo.ClusterName != clusterName {
					clusterExist = false
				}
			} else {
				clusterExist = false
			}
		}
	}
	if connectNode == nodeSum && clusterExist == false {
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
	localConf := readFile(settings.ClustersConfigFile)
	return NewClustersInfo(localConf)
}

func getRemoteNodes(clusterName string) interface{} {
	localConf := getLocalConf()
	nodeList := localConf.GetNodes(clusterName)
	return nodeList
}

// readFile reads a JSON file, decodes its content, and returns it as a map.
func readFile(filename string) map[string]interface{} {
	var newDict map[string]interface{}

	infile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return newDict
	}
	defer infile.Close()

	data, err := io.ReadAll(infile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return newDict
	}

	if err := json.Unmarshal(data, &newDict); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return newDict
	}

	return newDict
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
	for _, node := range nodeList {
		// Node-to-node config file sync operation
		url := fmt.Sprintf("http://%s:%s/remote/api/v1/sync_config", node, port)
		confJSON, err := json.Marshal(conf.Text)
		if err != nil {
			fmt.Println("Error marshaling config to JSON:", err)
			return
		}

		_, err = http.Post(url, "application/json", bytes.NewBuffer(confJSON))
		if err != nil {
			fmt.Println("Error syncing config to node:", err)
		}
	}

	fmt.Println("Sync complete")
}

// hostAuth performs host authentication using the provided information.
func hostAuth(authInfo map[string]interface{}) map[string]interface{} {
	authFailed := false
	nodeList := authInfo["node_list"].([]string)
	passwordList := authInfo["password"].([]string)
	fmt.Println(nodeList, passwordList)
	for i := 0; i < len(nodeList); i++ {
		authCmd := fmt.Sprintf(utils.CmdHostAuthNode, nodeList[i], passwordList[i])
		_, err := utils.RunCommand(authCmd)
		if err != nil {
			authFailed = true
			break
		}
	}

	if authFailed {
		return map[string]interface{}{
			"action": false,
			"error":  gettext.Gettext("host auth failed"),
		}
	}

	return map[string]interface{}{
		"action":  true,
		"message": gettext.Gettext("host auth success"),
	}
}

func hostAuthWithAddr(authInfo AuthInfo) AuthRetA {
	authFaild := false
	authFaildInfo := ""

	authCmd := fmt.Sprintf(utils.CmdHostAuthNodeWithAddr, authInfo.nodeList[0], authInfo.ip[0], authInfo.passWord[0])
	out, err := utils.RunCommand(authCmd)
	if err != nil {
		authFaild = true
		authFaildInfo = string(out)
	}
	if authFaild {
		return AuthRetA{
			Action:     false,
			Error:      gettext.Gettext("host auth failed"),
			DetailInfo: authFaildInfo,
		}
	}
	return AuthRetA{
		Action:  true,
		Message: gettext.Gettext("host auth success"),
	}
}

// ClusterAdd adds a new cluster using the provided node information.
// Returns results indicating the success or failure of the operation.
func ClusterAdd(nodeInfo map[string]interface{}) map[string]interface{} {
	authInfo := make(map[string]interface{})
	nodeList := make([]string, 0)
	passwords := make([]string, 0)
	nodeList = append(nodeList, nodeInfo["node_name"].(string))
	passwords = append(passwords, nodeInfo["password"].(string))

	authInfo["node_list"] = nodeList
	authInfo["password"] = passwords

	authRes := hostAuth(authInfo)

	if !authRes["action"].(bool) {
		return authRes
	}
	fmt.Println("send get cluster info request")
	url := fmt.Sprintf("http://%s:%s/remote/api/v1/managec/local_cluster_info", authInfo["node_list"], port)
	resp, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{
			"action": false,
			"error":  gettext.Gettext("add cluster failed")}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var NewClusterInfo Cluster
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(body, &NewClusterInfo)
		if err != nil {
			return map[string]interface{}{
				"action": false,
				"error":  gettext.Gettext("add cluster failed")}
		}
		localConf := getLocalConf()

		if localConf.IsClusterNameInUse(NewClusterInfo.ClusterName) {
			return map[string]interface{}{
				"action": false,
				"error":  gettext.Gettext("The cluster already exists, please do not add it again")}
		}

		localConf.AddCluster(NewClusterInfo)
		localConf.Save()
		syncClusterConfFile(localConf)
		return map[string]interface{}{
			"action": true,
			"error":  gettext.Gettext("add cluster success")}
	}

	return map[string]interface{}{
		"action": false,
		"error":  gettext.Gettext("add cluster failed")}
}

func ConvertClusterDataToSetupMap(clusterSetInfo ClusterData) map[string]interface{} {
	return convertClusterDataToSetupMap(clusterSetInfo)
}

// convertClusterDataToSetupMap convert ClusterData to setup map
func convertClusterDataToSetupMap(clusterSetInfo ClusterData) map[string]interface{} {
	clusterInfo := make(map[string]interface{})

	nodesData := clusterSetInfo.Data
	var data []map[string]interface{}
	for _, node := range nodesData {
		nodeMap := make(map[string]interface{})
		nodeMap["name"] = node.Name
		nodeMap["nodeid"] = node.NodeID
		nodeMap["password"] = node.Password

		// 将每个RingAddrData映射到一个新的map并添加到RingAddr切片中
		var ringAddr []map[string]string
		for _, addrData := range node.RingAddr {
			addrMap := make(map[string]string)
			addrMap["ring"] = addrData.Ring
			addrMap["ip"] = addrData.Ip
			ringAddr = append(ringAddr, addrMap)
		}

		nodeMap["ring_addr"] = ringAddr
		data = append(data, nodeMap)
	}

	clusterInfo["cluster_name"] = clusterSetInfo.Cluster_name
	clusterInfo["data"] = data
	return clusterInfo
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

// ClusterSetup performs the setup of a cluster using the provided cluster information.
func ClusterSetup(clusterSetInfo ClusterData) map[string]interface{} {
	authInfo := getAuthInfoFromClusterData(clusterSetInfo)

	// first: host auth
	authRes := hostAuth(authInfo)
	if !authRes["action"].(bool) {
		return authRes
	}
	// second: cluster setup
	res := clusterSetup(clusterSetInfo)
	if res["action"].(bool) {
		// third: cluster conf sync
		localConf := getLocalConf()
		clusterInfo := convertClusterDataToSetupMap(clusterSetInfo)
		localConf.AddCluster(clusterInfoParse(clusterInfo))
		localConf.Save()
		syncClusterConfFile(localConf) //TODO:check sync
	}
	return res
}

func ClusterRemove(RemoveInfo RemoveData) *RemoveRet {
	clusters := RemoveInfo.Cluster_name
	localConf := getLocalConf()
	removeRes := make([]bool, 0)
	faildCluster := make([]string, 0)
	for _, cluster := range clusters {
		res := localConf.DeleteCluster(cluster)
		removeRes = append(removeRes, res)
		if !res {
			faildCluster = append(faildCluster, cluster)
		}
		localConf.Save()
		syncClusterConfFile(localConf)
	}
	var RetData RemoveRet
	RetData.Action = true
	RetData.Faild_cluster = faildCluster
	RetData.Data = removeRes
	return &RetData
}

func AddNodes(AddNodesinfo AddNodesData) interface{} {
	localConf := getLocalConf()
	clusterName := AddNodesinfo.Cluster_name
	localClusterName := getClusterName()

	if localClusterName == clusterName {
		return LocalAddNodes(AddNodesinfo)
	}
	remoteNodeList := getRemoteNodes(clusterName).([]interface{})
	if len(remoteNodeList) > 0 {
		for _, node := range remoteNodeList {
			url := fmt.Sprintf("http://%s:%s/remote/api/v1/nodes/add_nodes", node, port)

			httpResp, _ := utils.SendRequest(url, "POST", AddNodesinfo.Data)
			if httpResp.StatusCode == http.StatusOK {
				httpRespData, _ := io.ReadAll(httpResp.Body)
				httpResp.Body.Close()
				var httpRespMessage map[string]interface{}
				json.Unmarshal(httpRespData, &httpRespMessage)

				url = fmt.Sprintf("http://%s:%s/remote/api/v1/managec/local_cluster_info", node, port)
				httpResp, _ = utils.SendRequest(url, "GET", nil)
				httpRespData, _ = io.ReadAll(httpResp.Body)
				httpResp.Body.Close()
				var remoteClusterInfo Cluster
				json.Unmarshal(httpRespData, &remoteClusterInfo)

				localConf.UpdateCluster(remoteClusterInfo.ClusterName, remoteClusterInfo)
				localConf.Save()
				syncClusterConfFile(localConf)

				return httpRespMessage
			}
		}
	}

	return map[string]interface{}{
		"action":     false,
		"error":      gettext.Gettext("No cluster found"),
		"detailInfo": gettext.Gettext("Unable to connect to cluster"),
	}
}

func ClusterDestroy(clustersJSON map[string]interface{}) map[string]interface{} {
	localConf := getLocalConf()
	clusterList := localConf.Clusters
	res := make([]bool, 0)
	failedClusterList := make([]string, 0)
	detailInfos := make([]string, 0)
	clusters := clustersJSON["cluster_name"].([]interface{})
	for _, desCluster := range clusters {
		nodeList := make([]string, 0)
		for _, cluster := range clusterList {
			if desCluster == cluster.ClusterName {
				nodeList = cluster.Nodes
			}
		}
		des := false
		detailInfo := gettext.Gettext("Unable to connect to cluster")
		for _, node := range nodeList {
			url := fmt.Sprintf("http://%s:%s/remote/api/v1/destroy_cluster", node, port)
			success := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println(r) // 处理异常
					}
				}()
				response, err := http.Get(url)
				if err != nil {
					panic(err) // 触发异常
				}
				defer response.Body.Close()
				var requestMessage map[string]interface{}
				err = json.NewDecoder(response.Body).Decode(&requestMessage)
				if err != nil {
					panic(err) // 触发异常
				}
				if requestMessage["action"].(bool) {
					success = true
				} else {
					detailInfo = requestMessage["detailInfo"].(string)
				}
			}()
			if success {
				des = true
				break
			}
		}
		if !des {
			res = append(res, false)
			failedClusterList = append(failedClusterList, desCluster.(string))
			detailInfos = append(detailInfos, detailInfo)
		} else {
			res = append(res, true)
		}
		localConf.DeleteCluster(desCluster.(string))
		localConf.Save()
		syncClusterConfFile(localConf)
	}
	return map[string]interface{}{
		"action":     true,
		"data":       res,
		"clusters":   failedClusterList,
		"detailInfo": detailInfos,
	}
}

// UrlRedirect
func UrlRedirect(clusterName string, uiPath string, requestMethod string, requestData interface{}) (map[string]interface{}, error) {
	remoteNodes := getRemoteNodes(clusterName).([]interface{})
	if len(remoteNodes) == 0 {
		return map[string]interface{}{
			"action":  false,
			"message": gettext.Gettext("Please reselect the cluster in the top operation area"),
		}, errors.New("no remote nodes")
	}

	// the first remote node searched
	node := remoteNodes[0]
	url := generateRemoteRequestURL(node.(string), uiPath)
	resp, err := utils.SendRequest(url, requestMethod, requestData)
	if err != nil {
		return map[string]interface{}{"action": false, "message": gettext.Gettext("Request remote Cluster info failed")}, err
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"action":  false,
			"message": gettext.Gettext("Please reselect the cluster in the top operation area"),
		}, err
	}
	defer resp.Body.Close()
	remoteClusterInfo := make(map[string]interface{})
	json.Unmarshal(respData, &remoteClusterInfo)
	return remoteClusterInfo, nil
}

func generateRemoteRequestURL(node string, uri string) string {
	if strings.HasPrefix(uri, "/remote") {
		return "http://" + node + ":" + port + uri
	}
	return "http://" + node + ":" + port + "/remote" + uri
}
