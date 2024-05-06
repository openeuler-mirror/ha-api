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

// getClusterName reads the cluster name from the corosync configuration file.
// Returns a map indicating the result and the extracted cluster name, if available.
func getClusterName() map[string]interface{} {
	filename := settings.CorosyncConfFile
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
func GetClusterInfo() map[string]interface{} {
	_, currentNode := utils.RunCommand(utils.CmdHostName)
	currentNodeStr := strings.ReplaceAll(fmt.Sprintf("%s", currentNode), "\n", "")

	if IsClusterExist() {
		nodeList := getNodeList()
		nodes := make([]map[string]interface{}, 0)

		index := 0
		for index < len(nodeList) {
			if nodeList[index] == "node {" {
				index++
				nodeInfo := make(map[string]interface{})
				nodeInfo["ring_addr"] = make([]map[string]string, 0)

				for index < len(nodeList) && nodeList[index] != "node {" {
					//d := make(map[string]string)
					d := make(map[string]interface{})
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
						if k != "nodeid" && k != "name" && k != "ring_addr" {
							if isIPv4(v.(string)) {
								ringAddr := map[string]string{
									"ring": k,
									"ip":   v.(string),
								}
								nodeInfo["ring_addr"] = append(nodeInfo["ring_addr"].([]map[string]string), ringAddr)
							} else {
								nodeInfo["ring_addr"] = append(nodeInfo["ring_addr"].([]map[string]string), map[string]string{})
							}
							delete(nodeInfo, k)
						}
					}
					count++
				}

				nodes = append(nodes, nodeInfo)
			} else {
				index++
			}
		}
		for _, node := range nodes {
			nodeID, ok := node["nodeid"].(string)
			if !ok {
				fmt.Println("Invalid nodeid format")
				continue
			}

			convertedID, err := strconv.Atoi(nodeID)
			if err != nil {
				fmt.Println("Failed to convert nodeid:", err)
				continue
			}

			node["nodeid"] = convertedID
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
		return map[string]interface{}{"action": false, "error": gettext.Getdata("Create cluster failed"), "detailInfo": outputStr}
	}

	_, err = utils.RunCommand(utils.CmdCreateAlert)
	if err != nil {
		return map[string]interface{}{"action": true, "message": gettext.Getdata("Create cluster success"),
			"alertInfo": gettext.Getdata("Failed to configure the alarm function module. If you need an alarm log, please manually execute the following command: pcs alert create id=alert_log path=/usr/share/pacemaker/alerts/alert_log.sh")}
	}

	return map[string]interface{}{"action": true, "message": gettext.Getdata("Create cluster success")}
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
				"error":      gettext.Getdata("Add node failed"),
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
