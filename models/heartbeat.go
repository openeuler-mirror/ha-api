/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 10:31:21
 * Description: 磁盘心跳相关功能的实现
 ******************************************************************************/
package models

import (
	"encoding/json"
	"strconv"

	"gitee.com/openeuler/ha-api/utils"

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
		return errors.New("no input data")
	}

	data := struct {
		Hbaddrs1        []HostInfo `json:"hbaddrs1"`
		Hbaddrs2        []HostInfo `json:"hbaddrs2"`
		Hbaddrs2Enabled int        `json:"hbaddrs2_enabled"`
	}{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return errors.New("invalid config data")
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

	return errors.New("there are running resources in the cluster, please close first")

ret:
	return errors.New("change cluster failed")
}
