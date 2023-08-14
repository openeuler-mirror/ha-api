package models

import (
	"bufio"
	"fmt"
	"openkylin.com/ha-api/utils"
	"os"
	"strconv"
	"strings"
)

//nodes_info格式
//{"nodeid": "1", "name": "HA1", "ring0_addr": "192.168.11.1", "ring1_addr": "192.168.11.4"},
//{"nodeid": "2", "name": "HA2", "ring0_addr": "192.168.11.2", "ring1_addr": "192.168.11.5"}

// getClusterName reads the cluster name from the corosync configuration file.
// Returns a map indicating the result and the extracted cluster name, if available.
func getClusterName() map[string]interface{} {
	filename := "/etc/corosync/corosync.conf"
	clusterName := ""
	result := make(map[string]interface{})

	file, err := os.Open(filename)
	if err != nil {
		result["action"] = false
		result["error"] = "File " + filename + " doesn't exist!"
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			str1 := strings.TrimSpace(parts[0])
			str2 := strings.TrimSpace(parts[1])
			if str1 == "cluster_name" {
				clusterName = str2
				break
			}
		}
	}

	result["action"] = true
	result["clusterName"] = clusterName
	return result
}

// getClusterInfo retrieves cluster information, including cluster nodes and their properties.
// Returns the cluster information in a structured map.
func getClusterInfo() map[string]interface{} {
	_, currentNode := utils.RunCommand("cat /etc/hostname")
	currentNodeStr := strings.ReplaceAll(fmt.Sprintf("%s", currentNode), "\n", "")

	if isClusterExist() {
		nodeList := getNodeList()
		nodes := make([]map[string]string, 0)

		index := 0
		for index < len(nodeList) {
			if nodeList[index] == "node {" {
				index++
				nodeInfo := make(map[string]string)
				nodeInfoRet := make(map[string]string)

				for index < len(nodeList) && nodeList[index] != "node {" {
					d := make(map[string]string)
					n := strings.Split(nodeList[index], ":")
					d[n[0]] = strings.TrimSpace(n[1])
					for k, v := range d {
						nodeInfo[k] = v
					}
					index++
				}

				count := 0
				for count < 2 {
					for k, v := range nodeInfo {
						item := map[string]string{k: v}
						for key, value := range item {
							nodeInfoRet[key] = value
						}
						count++
						delete(nodeInfo, k)
					}
				}

				for k, v := range nodeInfo {
					if k != "nodeid" && k != "name" {
						if isIPv4(v) {
							nodeInfoRet[k] = v
						} else {
							nodeInfoRet[k] = ""
						}
					}
				}

				nodes = append(nodes, nodeInfoRet)
			}
		}

		data := map[string]interface{}{
			"action":        true,
			"cluster_exist": true,
			"cluster_name":  getClusterName(),
			"currentNode":   currentNodeStr,
			"data":          nodes,
		}
		return data
	} else {
		data := map[string]interface{}{
			"action":        false,
			"cluster_exist": false,
			"cluster_name":  getClusterName(),
			"currentNode":   currentNodeStr,
			"error":         "Cluster not established!",
		}
		return data
	}
}

// clusterSetup sets up a cluster with the provided node information.
// Returns results indicating the success or failure of the cluster setup.
func clusterSetup(addNodes map[string]interface{}) map[string]interface{} {
	clusterName := ""
	for k, v := range addNodes {
		if k == "cluster_name" {
			clusterName = v.(string)
		}
	}
	if clusterName == "" {
		clusterName = "hacluster"
	}

	nodeCmdStr := generateNodeCmdStr(addNodes["data"].([]interface{}))
	cmd := "pcs cluster setup " + clusterName + nodeCmdStr + " totem token=8000 --start"
	output, err := utils.RunCommand(cmd)
	outputStr := string(output[:])
	if err != nil {
		if strings.Contains(outputStr, "Running cluster services") {
			return map[string]interface{}{"action": false, "error": "添加的部分节点已在集群中，请先将这些节点从集群中移除，或从已在集群中进行添加节点操作。"}
		}
		return map[string]interface{}{"action": false, "error": "集群创建失败"}
	} else {
		return map[string]interface{}{"action": true, "message": "集群创建成功"}
	}
}

// generateNodeCmdStr generates the command string for adding nodes to the cluster.
// Returns the generated command string.
func generateNodeCmdStr(nodesInfo []interface{}) string {
	hbIPPrefix := "addr="
	var cmd strings.Builder
	hbIPCmd := ""
	for _, nodeInfo := range nodesInfo {
		nodeInfoV := nodeInfo.(map[string]interface{})
		nodeStr := fmt.Sprintf("%v", nodeInfoV["name"])
		for k, v := range nodeInfoV {
			if k != "nodeid" && k != "name" && k != "port" && k != "password" {
				hbIPCmd = fmt.Sprintf(" %s%s", hbIPPrefix, v)
				nodeStr += hbIPCmd
			}
		}
		cmd.WriteString(" " + nodeStr)
	}

	return cmd.String()
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
