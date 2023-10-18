package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"openkylin.com/ha-api/utils"
	"os"
)

var clustersFileName = "/usr/share/heartbeat-gui/ha-api/ClustersInfo.conf"
var port = 8088

// ClusterInfo is a structure representing information about clusters.
type ClusterInfo struct {
	Text     map[string]interface{}
	Version  int
	Clusters []interface{}
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

// NewClustersInfo creates a new ClusterInfo instance using the provided text data.
// If the text data is nil or empty, default values are initialized.
func NewClustersInfo(text map[string]interface{}) *ClusterInfo {
	c := &ClusterInfo{
		Text: text,
	}

	if text == nil || len(text) == 0 {
		c.Text = make(map[string]interface{})
		c.Version = 0
		c.Clusters = make([]interface{}, 0)
		c.Text["version"] = c.Version
		c.Text["clusters"] = c.Clusters
	} else {
		c.Version = int(text["version"].(float64))
		c.Clusters = text["clusters"].([]interface{})
	}

	return c
}

// Save updates the version, performs a backup, and saves the ClusterInfo to a file in JSON format.
func (ci *ClusterInfo) Save() {
	ci.Version++
	ci.Backup()
	saveConf := ci.UpdateText()
	file, err := os.Create(clustersFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(saveConf); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

// Backup creates a backup of the cluster information file with a timestamp.
func (ci *ClusterInfo) Backup() {
	// Implement backup functionality
	// You can use the time and date in the file name
}

// UpdateText updates the version and clusters in the Text field and returns it.
func (ci *ClusterInfo) UpdateText() map[string]interface{} {
	ci.Text["version"] = ci.Version
	ci.Text["clusters"] = ci.Clusters
	return ci.Text
}

// AddCluster adds cluster information to the Clusters field.
func (ci *ClusterInfo) AddCluster(clusterInfo map[string]interface{}) {
	ci.Clusters = append(ci.Clusters, clusterInfo)
}

// IsClusterNameInUse checks if a cluster name is already in use.
func (ci *ClusterInfo) IsClusterNameInUse(clusterName string) bool {
	for _, c := range ci.Clusters {
		cV := c.(map[string]interface{})
		if cV["cluster_name"].(string) == clusterName {
			return true
		}
	}
	return false
}

// SetVersion sets the version of the ClusterInfo.
func (ci *ClusterInfo) SetVersion(version int) {
	ci.Version = version
}

// DeleteCluster delete the Cluster from ClusterInfo.
func (ci *ClusterInfo) DeleteCluster(clusterNameJson string) bool {
	for i, c := range ci.Clusters {
		cV := c.(map[string]interface{})
		if cV["cluster_name"] == clusterNameJson {
			ci.Clusters = append(ci.Clusters[:i], ci.Clusters[i+1:]...)
			return true
		}
	}
	return false
}

// localClusterInfo retrieves the cluster information locally and returns it as a map.
// If no cluster exists, an empty map is returned.
func localClusterInfo() map[string]interface{} {
	allInfo := GetClusterInfo()
	if allInfo["cluster_exist"] == true {
		clusterInfo := clusterInfoParse(allInfo)
		return clusterInfo
	}
	return make(map[string]interface{})
}

// clusterInfoParse takes cluster information as input and parses it into a map of string to interface
func clusterInfoParse(clusterInfo map[string]interface{}) map[string]interface{} {
	clusterParse := make(map[string]interface{})
	clusterParse["cluster_name"] = clusterInfo["cluster_name"]
	nodes := make([]string, 0)
	nodeIDs := make([]interface{}, 0)
	ips := make([]map[string]interface{}, 0)
	dataInterface, _ := clusterInfo["data"].([]interface{})
	var nodesInfo []map[string]interface{}
	for _, item := range dataInterface {
		nodesInfo = append(nodesInfo, item.(map[string]interface{}))
	}
	for _, nodeInfo := range nodesInfo {
		nodes = append(nodes, nodeInfo["name"].(string))
		nodeIDs = append(nodeIDs, nodeInfo["nodeid"])

		ip := make(map[string]interface{})
		for k, v := range nodeInfo {
			if k != "name" && k != "nodeid" && k != "password" {
				ip[k] = v
			}
		}
		ips = append(ips, ip)
	}

	clusterParse["nodes"] = nodes
	clusterParse["nodeid"] = nodeIDs
	clusterParse["ip"] = ips
	return clusterParse
}

// getLocalConf reads the local cluster configuration from a file and returns a ClusterInfo instance.
func getLocalConf() *ClusterInfo {
	localFile := readFile(clustersFileName)

	localConf := NewClustersInfo(localFile)
	return localConf

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

	data, err := ioutil.ReadAll(infile)
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

// comment out due to type error as localconf could not be {}, it should be of type *ClusterInfo
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
func syncClusterConfFile(conf *ClusterInfo) {
	// Get local cluster info
	clusterInfo := localClusterInfo()

	// If the current node has no cluster config, save the provided config
	if len(clusterInfo) == 0 {
		conf.Save()
		return
	}

	// Sync config file with all nodes in the cluster
	nodeList := clusterInfo["nodes"].([]string)
	for _, node := range nodeList {
		// Node-to-node config file sync operation
		url := fmt.Sprintf("https://%s:%s/remote/api/v1/sync_config", node, port)
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

	for i := 0; i < len(nodeList); i++ {
		authCmd := fmt.Sprintf("pcs host auth %s -u hacluster -p '%s'", nodeList[i], passwordList[i])

		_, err := utils.RunCommand(authCmd)

		if err != nil {
			authFailed = true
			break
		}
	}

	if authFailed {
		return map[string]interface{}{
			"action": false,
			"error":  "host auth failed",
		}
	}

	return map[string]interface{}{
		"action":  true,
		"message": "host auth success",
	}
}

// ClusterAdd adds a new cluster using the provided node information.
// Returns results indicating the success or failure of the operation.
func ClusterAdd(nodeInfo map[string]interface{}) map[string]interface{} {
	authInfo := make(map[string]interface{})
	authInfo["node_list"] = nodeInfo["node_name"].(string)
	authInfo["password"] = nodeInfo["password"].(string)

	authRes := hostAuth(authInfo)

	if !authRes["action"].(bool) {
		return authRes
	}

	url := fmt.Sprintf("https://%s:%s/remote/api/v1/managec/local_cluster_info", authInfo["node_list"], port)
	resp, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{
			"action": false,
			"error":  "添加集群失败"}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var NewClusterInfo map[string]interface{}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(body, &NewClusterInfo)
		if err != nil {
			return map[string]interface{}{
				"action": false,
				"error":  "添加集群失败"}
		}

		localConf := getLocalConf()
		if localConf.IsClusterNameInUse(NewClusterInfo["cluster_name"].(string)) {
			return map[string]interface{}{
				"action": false,
				"error":  "请勿重复添加"}
		}

		localConf.AddCluster(NewClusterInfo)
		localConf.Save()
		syncClusterConfFile(localConf)
		return map[string]interface{}{
			"action": true,
			"error":  "添加集群成功"}
	}

	return map[string]interface{}{
		"action": false,
		"error":  "添加集群失败"}
}

// ClusterSetup performs the setup of a cluster using the provided cluster information.
func ClusterSetup(clusterInfo map[string]interface{}) map[string]interface{} {
	authInfo := make(map[string]interface{})
	nodeList := make([]string, 0)
	passwords := make([]string, 0)

	dataInterface, _ := clusterInfo["data"].([]interface{})

	var data []map[string]interface{}
	for _, item := range dataInterface {
		data = append(data, item.(map[string]interface{}))
	}

	for _, node := range data {
		nodeList = append(nodeList, node["name"].(string))
		passwords = append(passwords, node["password"].(string))
	}

	authInfo["node_list"] = nodeList
	authInfo["password"] = passwords

	authRes := hostAuth(authInfo)
	if !authRes["action"].(bool) {
		return authRes
	} else {
		res := clusterSetup(clusterInfo)
		if res["action"].(bool) {
			localConf := getLocalConf()
			localConf.AddCluster(clusterInfoParse(clusterInfo))
			localConf.Save()
			syncClusterConfFile(localConf)
		}
		return res
	}
}

func ClusterRemove(RemoveInfo RemoveData) *RemoveRet {
	clusters := RemoveInfo.Cluster_name
	localConf := getLocalConf()
	removeRes := make([]bool, 0)
	faildCluster := make([]string, 0)
	for _, cluster := range clusters {
		res := localConf.DeleteCluster(cluster)
		removeRes = append(removeRes, res)
		if res == false {
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
