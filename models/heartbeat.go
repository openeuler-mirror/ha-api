/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Thu Jan 14 13:33:38 2021 +0800
 */
package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
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

func AddLink(linkIds string) error {
	cmd := fmt.Sprintf(utils.CmdAddLink, linkIds)
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
