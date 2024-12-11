/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Thu Jan 14 13:33:38 2021 +0800
 */
package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/chai2010/gettext-go"

	"errors"
)

type HostInfo struct {
	IP     string `json:"ip"`
	NodeID string `json:"nodeid"`
}

func GetHeartBeatHosts() ([]HostInfo, error) {
	knownHosts := []HostInfo{}

	out, err := utils.RunCommand(utils.CmdGetPcsdAuthFile)
	if err != nil {
		return nil, errors.New("no node in the Cluster, please run \"pcs host auth $nodename\" to add node")
	}

	jsonData := map[string]interface{}{}

	if err := json.Unmarshal(out, &jsonData); err != nil {
		return nil, errors.New("parse host info failed")
	}
	for k := range jsonData["known_hosts"].(map[string]interface{}) {
		hostInfo := HostInfo{
			IP:     "",
			NodeID: k,
		}
		knownHosts = append(knownHosts, hostInfo)
	}

	return knownHosts, nil
}

func GetHeartBeatDictionary() (interface{}, error) {
	nodeList, err := utils.GetNodeList()
	if err != nil {
		return nil, err
	}

	res := map[string][]map[string]string{}
	for _, node := range nodeList {
		name := node["name"]
		for k, addr := range node {
			if k != "name" && k != "nodeid" {
				info := map[string]string{}
				info["nodeid"] = name
				info["ip"] = addr

				if _, ok := res[k]; !ok {
					res[k] = []map[string]string{}
				}
				res[k] = append(res[k], info)
			}
		}
	}

	ret := map[string][]map[string]string{}
	count := 0
	for _, value := range res {
		count++
		hbStr := "hbaddrs" + strconv.Itoa(count)
		ret[hbStr] = value
	}

	return ret, nil
}

func GetHeartBeatConfig() (interface{}, error) {
	var rst interface{}
	var err error
	if utils.IsClusterExist() {
		rst, err = GetHeartBeatDictionary()
	} else {
		rst, err = GetHeartBeatHosts()
	}

	if err != nil {
		return nil, err
	}
	return rst, nil
}

func EditHeartbeatInfo(jsonData []byte) error {
	if len(jsonData) == 0 {
		return errors.New(gettext.Gettext("No input data"))
	}

	data := struct {
		Hbaddrs1        []HostInfo `json:"hbaddrs1"`
		Hbaddrs2        []HostInfo `json:"hbaddrs2"`
		Hbaddrs2Enabled int        `json:"hbaddrs2_enabled"`
	}{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return errors.New(gettext.Gettext("invalid config data"))
	}

	heartBeatInfos := map[string]([]string){}
	if len(data.Hbaddrs1) > 0 {
		for _, info := range data.Hbaddrs1 {
			if _, ok := heartBeatInfos[info.NodeID]; !ok {
				heartBeatInfos[info.NodeID] = []string{}
			}
			heartBeatInfos[info.NodeID] = append(heartBeatInfos[info.NodeID], info.IP)
		}
	}
	if len(data.Hbaddrs2) > 0 {
		for _, info := range data.Hbaddrs2 {
			if _, ok := heartBeatInfos[info.NodeID]; !ok {
				heartBeatInfos[info.NodeID] = []string{}
			}
			heartBeatInfos[info.NodeID] = append(heartBeatInfos[info.NodeID], info.IP)
		}
	}

	cmd := utils.CmdSetupCluster
	for key, value := range heartBeatInfos {
		cmd = cmd + " " + key
		for _, v := range value {
			addr := "addr=" + v
			cmd = cmd + " " + addr
		}
	}
	cmd = cmd + " --force --start"

	runResource := false
	// TODO: check logic
	if _, err := utils.RunCommand(utils.CmdClusterStatusAsXML); err != nil {
		// means a cluster is already running
		runResource = false
	} else {
		runResource = true
	}

	if runResource {
		if _, err := utils.RunCommand(utils.CmdSaveCIB); err != nil {
			goto ret
		}
		if _, err := utils.RunCommand(utils.CmdResourceCleanup); err != nil {
			goto ret
		}
		if _, err := utils.RunCommand(utils.CmdStopCluster); err != nil {
			goto ret
		}
		if _, err := utils.RunCommand(utils.CmdDestroyCluster); err != nil {
			goto ret
		}
		if _, err := utils.RunCommand(cmd); err != nil {
			goto ret
		}
		if _, err := utils.RunCommand(utils.CmdPushFileToCIB); err != nil {
			goto ret
		}

		return nil
	}

	return errors.New(gettext.Gettext("there are running resources in the cluster, please close first"))

ret:
	return errors.New(gettext.Gettext("Change cluster failed"))
}

func DeletLinks(linkIds string) error {
	cmd := fmt.Sprintf(utils.CmdDeleteLinks, linkIds)
	_, err := utils.RunCommand(cmd)
	return err
}

func AddLink(hbInfo map[string]string, linkId string) error {
	hbInfoStr := GenerateLinkStr(hbInfo)
	var cmd string
	if linkId == "" {
		cmd = fmt.Sprintf(utils.CmdAddLinkForce, hbInfoStr)
	} else {
		cmd = fmt.Sprintf(utils.CmdAddLinksWithLinkNum, hbInfoStr, linkId)
	}

	_, err := utils.RunCommand(cmd)
	return err
}

func EditLinks(hbInfo map[string]string, linkId string) error {
	hbInfoStr := GenerateLinkStr(hbInfo)
	cmd := fmt.Sprintf(utils.CmdUpdateLinkForce, hbInfoStr, linkId)
	_, err := utils.RunCommand(cmd)
	return err
}

func GenerateLinkStr(data map[string]string) string {
	var linkStr strings.Builder
	for k, v := range data {
		linkStr.WriteString(" ")
		linkStr.WriteString(k)
		linkStr.WriteString("=")
		linkStr.WriteString(v)
	}
	return linkStr.String()[1:]
}

// HBInfo represents the heartbeat information structure
type HBInfo struct {
	Data        []map[string]string `json:"data"`
	ClusterName string              `json:"cluster_name"`
}

// ExtractHbInfo reorganizes heartbeat information format
func ExtractHbInfo(hbInfo []map[string]string) ([]map[string]string, []string) {
	var hbDictList []map[string]string
	var ids []string

	if len(hbInfo) == 0 {
		return nil, nil
	}
	var numIds []string
	for key := range hbInfo[0] {
		if strings.HasPrefix(key, "ring") {
			numIds = append(numIds, key[4:5]) // ring0-ring7
		}
	}

	for _, id := range numIds {
		hbDict := map[string]string{}
		for _, nodeInfo := range hbInfo {
			nodeName := nodeInfo["name"]
			ringAddr := nodeInfo["ring"+id+"_addr"]
			if ringAddr != "" {
				hbDict[nodeName] = nodeInfo["ring"+id+"_addr"]
			}
		}
		hbDictList = append(hbDictList, hbDict)
	}
	return hbDictList, ids
}

// get net heartbeat info from corosync conf
func ExtractHbInfoFromConf() (map[string][]string, []string) {
	linksInfo := make(map[string][]string)
	hosts := []string{}

	file, err := os.Open(settings.CorosyncConfFile)
	if err != nil {
		logs.Error("Error opening file: %v", err)
		return linksInfo, hosts
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ring") && !strings.Contains(line, "disk") {
			id := strings.TrimSpace(line)[4:5]
			parts := strings.Fields(line)
			if len(parts) > 1 {
				ip := parts[1]
				linksInfo[id] = append(linksInfo[id], ip)
			}
		}
		if strings.Contains(line, "name:") && !strings.Contains(line, "cluster_name") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				hosts = append(hosts, parts[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logs.Error("Error reading file: %v", err)
	}

	return linksInfo, hosts
}

func GetRingIdFromIPOnline(ipAddress string) (int, error) {
	// IPv4 check
	if !isValidIPv4(ipAddress) {
		return -1, fmt.Errorf("%s is not an IPv4 address", ipAddress)
	}

	output, err := utils.RunCommand(utils.CmdHbStatus)
	if err != nil {
		return -1, err
	}

	// Parse the output
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "LINK") {
			linkID := strings.Fields(line)[1]
			ipStr := strings.Split(strings.Split(line, "udp (")[1], ")")[0]
			ips := strings.Split(ipStr, "->")
			for _, ip := range ips {
				if strings.TrimSpace(ip) == ipAddress {
					id, err := strconv.Atoi(linkID)
					if err != nil {
						return -1, err
					}
					return id, nil
				}
			}
		}
	}

	return -1, fmt.Errorf("unable to find the link for the corresponding IP information")
}

// getRingIDFromIPOffline attempts to get the heartbeat ID from an IP address using an offline method.
func GetRingIdFromIPOffline(ipAddress string) (int, error) {
	file, err := os.Open(settings.CorosyncConfFile)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ring") && strings.Contains(line, ipAddress) {
			idStr := strings.Split(strings.TrimSpace(line), "_")[0]
			id, err := strconv.Atoi(idStr[len(idStr)-1:])
			if err != nil {
				return -1, err
			}
			return id, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return -1, err
	}

	return -1, fmt.Errorf("unable to find the link for the corresponding IP information")
}

// isValidIPv4 checks if the given string is a valid IPv4 address.
func isValidIPv4(ip string) bool {
	return net.ParseIP(ip) != nil && strings.Contains(ip, ".")
}

func DeleteHeartbeat(hbInfo HBInfo) utils.GeneralResponse {
	hbData := hbInfo.Data
	minLinkNum := 1
	hbInfoList, _ := ExtractHbInfo(hbData)
	status := GetClusterStatus()
	currentLocalHbIds := GetCurrentLinkIds()
	var linkId int
	var linkIdSet []string
	if len(currentLocalHbIds) <= minLinkNum || len(currentLocalHbIds) <= len(hbInfoList) {
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("At least one heartbeat needs to be preserved and cannot be deleted further."),
		}
	}

	for _, hbInfo := range hbInfoList {
		ip := utils.Values(hbInfo)
		if status == 0 {
			linkId, _ = GetRingIdFromIPOnline(ip[0])
		} else {
			linkId, _ = GetRingIdFromIPOnline(ip[0])
		}

		if linkId == -1 {
			linkId, _ = GetRingIdFromIPOnline(ip[0])
			if linkId == -1 {
				return utils.GeneralResponse{
					Action: false,
					Error:  gettext.Gettext("Exception occurred while deleting the heartbeat. Please check the logs and the original heartbeat status."),
				}
			}
		}
		linkIdSet = append(linkIdSet, strconv.Itoa(linkId))
	}
	linkIdSetStr := strings.Join(linkIdSet, " ")
	cmd := fmt.Sprintf(utils.CmdDeleteLinks, linkIdSetStr)
	if _, err := utils.RunCommand(cmd); err != nil {
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Deletion heartbeat failed, the heartbeat ID does not exist"),
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Error:  gettext.Gettext("Delete heartbeat success"),
	}
}

func GetCurrentLinkIds() []string {
	linksInfo, _ := ExtractHbInfoFromConf()
	ids := make([]string, 0, len(linksInfo))
	for k := range linksInfo {
		ids = append(ids, k)
	}
	return ids
}
