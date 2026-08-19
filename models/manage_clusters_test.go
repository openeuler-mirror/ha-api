/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

const testClusterName = "test-cluster"
const testIP1 = "192.0.2.1"

// MockIsClusterExist 模拟 IsClusterExist 函数
var MockIsClusterExist = func() bool {
	return false
}

// MockGetClusterInfo 模拟 GetClusterInfo 函数
var MockGetClusterInfo = func() map[string]interface{} {
	return map[string]interface{}{
		"cluster_exist": false,
	}
}

// MockRunCommand 模拟 utils.RunCommand 函数
var MockRunCommand = func(cmd string) ([]byte, error) {
	return []byte("output"), nil
}

// MockGetRemoteNodes 模拟 getRemoteNodes 函数
func MockGetRemoteNodes(nodes []string) func(string) interface{} {
	return func(clusterName string) interface{} {
		return nodes
	}
}

// MockSendRequest 模拟 utils.SendRequest 函数
func MockSendRequest(handler http.HandlerFunc) func(string, string, interface{}) (*http.Response, error) {
	return func(url, method string, data interface{}) (*http.Response, error) {
		// 创建一个测试服务器
		ts := httptest.NewServer(handler)
		defer ts.Close()

		// 发送请求到测试服务器
		resp, err := http.Get(ts.URL)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func TestUrlRedirect_AllNodesFail(t *testing.T) {
	// 设置mock
	originalGetRemoteNodes := getRemoteNodes
	originalSendRequest := utils.SendRequest
	defer func() {
		getRemoteNodes = originalGetRemoteNodes
		utils.SendRequest = originalSendRequest
	}()

	getRemoteNodes = MockGetRemoteNodes([]string{"node1", "node2"})
	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	// 调用函数
	result, err := UrlRedirect(testClusterName, "/path", "GET", nil, nil)

	// 验证结果
	if err == nil || err.Error() != "no nodes succeeded" {
		t.Errorf("Expected 'no nodes succeeded' error, got %v", err)
	}

	if result["action"].(bool) != false {
		t.Error("Expected action to be false when all nodes fail")
	}
}

func TestUrlRedirect_NoRemoteNodes(t *testing.T) {
	// 保存原始函数并设置mock
	originalGetRemoteNodes := getRemoteNodes
	defer func() { getRemoteNodes = originalGetRemoteNodes }()
	getRemoteNodes = MockGetRemoteNodes([]string{})

	// 调用函数
	result, err := UrlRedirect(testClusterName, "/path", "GET", nil, nil)

	// 验证结果
	if err == nil || err.Error() != "no remote nodes" {
		t.Errorf("Expected 'no remote nodes' error, got %v", err)
	}

	if result["action"].(bool) != false {
		t.Error("Expected action to be false when no remote nodes")
	}
}

// 保存原始函数以便恢复
var originalGetRemoteNodes = getRemoteNodes
var originalSendRequest = utils.SendRequest

// setupTest 准备测试环境
func setupTest() *bytes.Buffer {
	// 捕获日志
	var logBuf bytes.Buffer
	h := slog.NewJSONHandler(&logBuf, nil)
	slog.SetDefault(slog.New(h))

	return &logBuf
}

// teardownTest 清理测试环境
func teardownTest() {
	getRemoteNodes = originalGetRemoteNodes
	utils.SendRequest = originalSendRequest
}

func TestActionField_False(t *testing.T) {
	logBuf := setupTest()
	defer teardownTest()

	getRemoteNodes = func(clusterName string) interface{} {
		return []string{"node1", "node2"} // 多个节点测试 continue 逻辑
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"action": false,
			"error":  "something wrong",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	UrlRedirect(testClusterName, "/path", "GET", nil, nil)

	// 验证日志
	if !strings.Contains(logBuf.String(), "something wrong") {
		t.Error("expected warning log about invalid action")
	}
}

func TestNewClustersInfo_EmptyText(t *testing.T) {
	// 测试空map的情况
	emptyText := make(map[string]interface{})
	ci := NewClustersInfo(emptyText)

	assert.Equal(t, 0, ci.Version, "Version should be 0 for empty text")
	assert.Empty(t, ci.Clusters, "Clusters should be empty for empty text")
	assert.Equal(t, int(0), ci.Text["version"], "Text version should be 0")
	assert.Equal(t, []Cluster{}, ci.Text["clusters"], "Text clusters should be empty slice")
}

func TestNewClustersInfo_WithValidData(t *testing.T) {
	jsonData := `{
		"clusters": [
		  {
			"cluster_name": "v11-1",
			"nodes": [
			  "host103",
			  "host104",
			  "host105"
			],
			"nodeid": [
			  "1",
			  "2",
			  "3"
			],
			"ip": [
			  {
				"ring0_addr": "192.168.122.103",
				"status": {
				  "ring0_addr": "up"
				},
				"type": "primitive"
			  },
			  {
				"ring0_addr": "192.168.122.104",
				"status": {
				  "ring0_addr": "localhost"
				},
				"type": "primitive"
			  },
			  {
				"ring0_addr": "192.168.122.105",
				"status": {
				  "ring0_addr": "down"
				},
				"type": "primitive"
			  }
			]
		  },
		  {
			"cluster_name": "HA_114_115",
			"nodes": [
			  "host114",
			  "host115"
			],
			"nodeid": [
			  "1",
			  "2"
			],
			"ip": [
			  {
				"ring0_addr": "192.168.122.114",
				"status": {
				  "ring0_addr": "localhost"
				},
				"type": "primitive"
			  },
			  {
				"ring0_addr": "192.168.122.115",
				"status": {
				  "ring0_addr": "up"
				},
				"type": "primitive"
			  }
			]
		  }
		],
		"version": 1
	  }`

	var data map[string]interface{}
	json.Unmarshal([]byte(jsonData), &data)
	ci := NewClustersInfo(data)
	clusterName := ci.GetClusterNameOfNode("host114")
	if clusterName != "HA_114_115" {
		t.Errorf("GetClusterNameOfNode failed")
	}
	nodes := ci.GetNodes("HA_114_115")
	expectedNodes := []string{"host114", "host115"}
	assert.ElementsMatch(t, expectedNodes, nodes)

	ci.SetVersion(2)
	ci.DeleteCluster("HA_114_115")
	ci.UpdateText()
	res := ci.IsClusterNameInUse("v11-1")
	if !res {
		t.Errorf("cluster name check failed")
	}
	// settings.ClustersConfigFile = "temp"
	assert.Equal(t, 2, ci.Version, "Version should be 1")
	assert.Len(t, ci.Clusters, 1, "Should have 1 clusters")
}

func TestNewClustersInfo_WithInvalidClustersData(t *testing.T) {
	// 测试clusters数据无效的情况
	testText := map[string]interface{}{
		"version":  float64(1),
		"clusters": "invalid_data", // 不是slice类型
	}

	ci := NewClustersInfo(testText)

	assert.Equal(t, 1, ci.Version, "Version should still be 1")
	assert.Empty(t, ci.Clusters, "Clusters should be empty due to invalid data")
}

func TestNewClustersInfo_WithInvalidClusterItems(t *testing.T) {
	// 测试包含无效cluster项的情况
	testText := map[string]interface{}{
		"version": float64(1),
		"clusters": []interface{}{
			"invalid_cluster_data", // 不是map类型
			map[string]interface{}{
				"id":   "valid_cluster",
				"name": "Valid Cluster",
			},
		},
	}

	ci := NewClustersInfo(testText)

	assert.Equal(t, 1, ci.Version, "Version should be 1")
	assert.Len(t, ci.Clusters, 1, "Should have only 1 valid cluster")
}

func TestMapToCluster(t *testing.T) {
	// 测试正常情况
	clusterMap := map[string]interface{}{
		"cluster_name": testClusterName,
		"nodes":        []string{"node1", "node2"},
		"nodeid":       []string{"1", "2"},
		"ip": []map[string]interface{}{
			{"ring0_addr": testIP1},
			{"ring0_addr": "192.0.2.2"},
		},
	}

	cluster, err := MapToCluster(clusterMap)
	assert.NoError(t, err)
	assert.Equal(t, testClusterName, cluster.ClusterName)
	assert.ElementsMatch(t, []string{"node1", "node2"}, cluster.Nodes)
	assert.ElementsMatch(t, []string{"1", "2"}, cluster.Nodeid)
	assert.Len(t, cluster.Ip, 2)

	// 测试错误情况
	_, err = MapToCluster(map[string]interface{}{"cluster_name": make(chan int)})
	assert.Error(t, err)
}

func TestGetRemoteCluster(t *testing.T) {
	// 设置mock
	originalSendRequest := utils.SendRequest
	defer func() { utils.SendRequest = originalSendRequest }()

	// 测试成功情况
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Cluster{
			ClusterName: testClusterName,
			Nodes:       []string{"node1", "node2"},
			Nodeid:      []string{"1", "2"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	cluster, err := GetRemoteCluster("node1")
	assert.NoError(t, err)
	assert.NotNil(t, cluster)
	assert.Equal(t, testClusterName, cluster.ClusterName)

	// 测试失败情况
	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	cluster, err = GetRemoteCluster("node1")
	assert.Error(t, err)
	assert.Nil(t, cluster)
}

func TestClusterAdd(t *testing.T) {
	// 设置mock
	originalSendRequest := utils.SendRequest
	originalRunCommand := utils.RunCommand
	defer func() {
		utils.SendRequest = originalSendRequest
		utils.RunCommand = originalRunCommand
	}()

	// 模拟节点认证成功
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte("output"), nil
	}

	// 模拟集群信息获取成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Cluster{
			ClusterName: testClusterName,
			Nodes:       []string{"node1", "node2"},
			Nodeid:      []string{"1", "2"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	// 测试集群添加
	nodeInfo := ClusterAddReq{
		NodeName: "node1",
		PassWord: "password",
	}

	result := ClusterAdd(nodeInfo)
	assert.True(t, result.Action)
}

func TestClusterSetupInManageClusters(t *testing.T) {
	// 设置mock
	originalSendRequest := utils.SendRequest
	defer func() {
		utils.SendRequest = originalSendRequest
	}()

	// 生成一个唯一的集群名称，避免冲突
	clusterName := "test-cluster-" + t.Name()

	// 模拟集群设置成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method == "POST" {
			// 处理setup_cluster请求
			response := map[string]interface{}{
				"action": true,
			}
			json.NewEncoder(w).Encode(response)
		} else if r.Method == "GET" {
			// 处理local_cluster_info请求
			clusterInfo := Cluster{
				ClusterName: testClusterName,
				Nodes:       []string{"node1", "node2"},
				Nodeid:      []string{"1", "2"},
			}
			json.NewEncoder(w).Encode(clusterInfo)
		}
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		// 所有请求都发送到测试服务器
		client := &http.Client{}
		req, err := http.NewRequest(method, ts.URL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// 测试集群设置
	clusterData := ClusterData{
		Cluster_name: clusterName,
		Data: []NodeData{
			{Name: "node1", Password: "password", RingAddr: []string{testIP1}},
		},
	}

	result := ClusterSetup(clusterData)
	// 添加调试信息
	t.Logf("ClusterSetup result: %v", result)
	// 检查result是否包含action字段，并且其值为true
	assert.Contains(t, result, "action", "Result should contain 'action' field")
	if action, ok := result["action"].(bool); ok {
		assert.True(t, action, "Action should be true")
	} else {
		t.Errorf("result['action'] is not a bool: %v", result["action"])
		if err, ok := result["error"].(string); ok {
			t.Errorf("Error message: %s", err)
		}
		if detail, ok := result["detailInfo"].(string); ok {
			t.Errorf("Detail info: %s", detail)
		}
	}
}

func TestClusterRemove(t *testing.T) {
	// 创建一个测试集群配置
	testCluster := Cluster{
		ClusterName: testClusterName,
		Nodes:       []string{"node1", "node2"},
		Nodeid:      []string{"1", "2"},
	}

	// 获取本地配置并添加测试集群
	localConf := getLocalConf()
	localConf.AddCluster(testCluster)

	// 测试集群移除
	removeData := RemoveData{
		Cluster_name: []string{testClusterName},
	}

	result := ClusterRemove(removeData)
	assert.True(t, result.Action)
	assert.Len(t, result.FailedCluster, 0)
	assert.ElementsMatch(t, []bool{true}, result.Data)
}

func TestClusterDestroy(t *testing.T) {
	// 设置mock
	originalGetRemoteNodes := getRemoteNodes
	originalSendRequest := utils.SendRequest
	defer func() {
		getRemoteNodes = originalGetRemoteNodes
		utils.SendRequest = originalSendRequest
	}()

	// 模拟远程节点
	getRemoteNodes = MockGetRemoteNodes([]string{"node1", "node2"})

	// 模拟集群销毁成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"action": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	// 测试集群销毁
	clustersJSON := map[string]interface{}{
		"cluster_name": []interface{}{testClusterName},
	}

	result := ClusterDestroy(clustersJSON)
	assert.True(t, result["action"].(bool))
	assert.ElementsMatch(t, []bool{true}, result["data"].([]bool))
	assert.Len(t, result["clusters"].([]string), 0)
}

func TestLocalClusterInfo(t *testing.T) {
	// 测试无集群情况
	result := LocalClusterInfo()
	assert.Empty(t, result.ClusterName)
	assert.Empty(t, result.Nodes)
}

func TestClusterInfo(t *testing.T) {
	result := ClusterInfo()
	assert.Contains(t, result, "action")
	assert.Contains(t, result, "cluster_list")
}

func TestClusterOverview(t *testing.T) {
	result := ClusterOverview()
	assert.Contains(t, result, "action")
	assert.Contains(t, result, "cluster_exist")
	assert.Contains(t, result, "local_cluster_name")
	assert.Contains(t, result, "cluster_data")
}

func TestSetupInNode(t *testing.T) {
	// 设置mock
	originalSendRequest := utils.SendRequest
	defer func() {
		utils.SendRequest = originalSendRequest
	}()

	// 模拟请求成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回成功响应
		response := map[string]interface{}{
			"action": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		// 验证URL和方法
		assert.Contains(t, url, "node1")
		assert.Equal(t, "POST", method)

		// 返回测试服务器的响应
		resp, err := http.Get(ts.URL)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// 测试数据
	clusterData := ClusterData{
		Cluster_name: testClusterName,
		Data: []NodeData{
			{Name: "node1", Password: "password", RingAddr: []string{testIP1}},
		},
	}

	// 调用函数
	resp, err := setupInNode("node1", clusterData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// 检查响应
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "action")
	if action, ok := response["action"].(bool); ok {
		assert.True(t, action)
	}
}

func TestGetAuthInfoFromClusterData(t *testing.T) {
	// 测试数据
	clusterData := ClusterData{
		Cluster_name: testClusterName,
		Data: []NodeData{
			{Name: "node1", Password: "pass1"},
			{Name: "node2", Password: "pass2"},
		},
	}

	// 调用函数
	authInfo := getAuthInfoFromClusterData(clusterData)

	// 验证结果
	assert.Contains(t, authInfo, "node_list")
	assert.Contains(t, authInfo, "password")

	nodeList := authInfo["node_list"].([]string)
	passwords := authInfo["password"].([]string)

	assert.Len(t, nodeList, 2)
	assert.Len(t, passwords, 2)
	assert.Contains(t, nodeList, "node1")
	assert.Contains(t, nodeList, "node2")
	assert.Contains(t, passwords, "pass1")
	assert.Contains(t, passwords, "pass2")
}

func TestGetNodeListFromClusterData(t *testing.T) {
	// 测试数据
	clusterData := ClusterData{
		Cluster_name: testClusterName,
		Data: []NodeData{
			{Name: "node1", Password: "pass1"},
			{Name: "node2", Password: "pass2"},
			{Name: "node3", Password: "pass3"},
		},
	}

	// 调用函数
	nodeList := getNodeListFromClusterData(clusterData)

	// 验证结果
	assert.Len(t, nodeList, 3)
	assert.Contains(t, nodeList, "node1")
	assert.Contains(t, nodeList, "node2")
	assert.Contains(t, nodeList, "node3")
}

func TestClusterDestroy2(t *testing.T) {
	// 设置mock
	originalGetRemoteNodes := getRemoteNodes
	originalSendRequest := utils.SendRequest
	defer func() {
		getRemoteNodes = originalGetRemoteNodes
		utils.SendRequest = originalSendRequest
	}()

	// 模拟远程节点
	getRemoteNodes = MockGetRemoteNodes([]string{"node1", "node2"})

	// 模拟集群销毁成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"action": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	// 测试数据
	clustersJSON := map[string]interface{}{
		"cluster_name": []interface{}{testClusterName},
	}

	// 调用函数
	result := ClusterDestroy2(clustersJSON)

	// 验证结果
	assert.True(t, result["action"].(bool))
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "clusters")
	assert.Contains(t, result, "detailInfo")
}

func TestUrlRedirectWithSyncConfig(t *testing.T) {
	// 设置mock
	originalGetRemoteNodes := getRemoteNodes
	originalSendRequest := utils.SendRequest
	defer func() {
		getRemoteNodes = originalGetRemoteNodes
		utils.SendRequest = originalSendRequest
	}()

	// 模拟远程节点
	getRemoteNodes = MockGetRemoteNodes([]string{"node1"})

	// 模拟请求成功
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "local_cluster_info") {
			// 返回集群信息
			clusterInfo := Cluster{
				ClusterName: testClusterName,
				Nodes:       []string{"node1", "node2"},
				Nodeid:      []string{"1", "2"},
			}
			json.NewEncoder(w).Encode(clusterInfo)
		} else {
			// 返回操作成功
			response := map[string]interface{}{
				"action": true,
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer ts.Close()

	utils.SendRequest = func(url, method string, data interface{}) (*http.Response, error) {
		return http.Get(ts.URL)
	}

	// 调用函数
	result, err := UrlRedirectWithSyncConfig(testClusterName, "/path", "GET", nil)

	// 验证结果
	assert.NoError(t, err)
	assert.Contains(t, result, "action")
}

func TestUpdateClusterConfFile(t *testing.T) {
	// 测试数据
	cluster := Cluster{
		ClusterName: testClusterName,
		Nodes:       []string{"node1", "node2"},
		Nodeid:      []string{"1", "2"},
	}

	// 调用函数
	err := UpdateClusterConfFile(cluster)

	// 验证结果
	assert.NoError(t, err)
}

func TestSyncConfig(t *testing.T) {
	// 测试数据 - 使用float64类型的version，因为JSON解析会返回float64
	remoteConf := map[string]interface{}{
		"version":  float64(1),
		"clusters": []interface{}{},
	}

	// 调用函数
	result := SyncConfig(remoteConf)

	// 验证结果
	assert.Contains(t, result, "result")
}

func TestHostAuth(t *testing.T) {
	// 设置mock
	originalRunCommand := utils.RunCommand
	defer func() {
		utils.RunCommand = originalRunCommand
	}()

	// 模拟命令执行成功
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte("output"), nil
	}

	// 测试数据
	authInfo := AuthRequest{
		NodeList:  []string{"node1", "node2"},
		Passwords: []string{"pass1", "pass2"},
	}

	// 调用函数
	result := hostAuth(authInfo)

	// 验证结果
	assert.True(t, result.Action)
}

func TestHostAuthWithAddr(t *testing.T) {
	// 设置mock
	originalRunCommand := utils.RunCommand
	defer func() {
		utils.RunCommand = originalRunCommand
	}()

	// 模拟命令执行成功
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte("output"), nil
	}

	// 测试数据
	authInfo := AuthInfo{
		nodeList: []string{"node1"},
		ip:       []string{testIP1},
		passWord: []string{"pass1"},
	}

	// 调用函数
	result := hostAuthWithAddr(authInfo)

	// 验证结果
	assert.True(t, result.Action)
}

func TestExtractIPs(t *testing.T) {
	// 测试数据
	cluster := Cluster{
		ClusterName: testClusterName,
		Ip: []map[string]interface{}{
			{"ring0_addr": testIP1},
			{"ring0_addr": "192.0.2.2"},
		},
	}

	// 调用函数
	ips := extractIPs(cluster)

	// 验证结果
	assert.Len(t, ips, 2)
	assert.Contains(t, ips[0].Addrs, "ring0")
	assert.Contains(t, ips[1].Addrs, "ring0")
}
