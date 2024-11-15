/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bixiaoyan
 * Date: 2024-03-07 14:49:51
 * LastEditTime: 2024-03-07 17:37:48
 * Description: 集群利用率功能
 ******************************************************************************/
package models

import (
	"encoding/json"
	"strings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

type UtilizationData struct {
	Name  string                 `json:"name"`
	Attri map[string]interface{} `json:"attri"`
}

type GetUtilResult struct {
	Action bool              `json:"action"`
	Data   *UtilResponseData `json:"data"`
}

type UtilResponseData struct {
	NodeUtilization []UtilizationData `json:"NodeUtilization"`
	ResUtilization  []UtilizationData `json:"ResUtilization"`
}

func GetUtilization() GetUtilResult {
	result := GetUtilResult{}
	nodeUtilization := GetOneTypeUtilization("node")
	resUtilization := GetOneTypeUtilization("resource")
	result.Action = true
	result.Data = &UtilResponseData{
		NodeUtilization: nodeUtilization,
		ResUtilization:  resUtilization,
	}

	return result
}

func GetOneTypeUtilization(Uti_type string) []UtilizationData {
	data := []UtilizationData{}
	cmd := "pcs " + Uti_type + " utilization"
	output, _ := utils.RunCommand(cmd)
	Util := strings.Split(string(output), "Utilization:")[1]

	Util = strings.TrimSpace(Util)
	if Util == "" {
		return nil
	}

	UtilList := strings.Split(Util, "\n")

	for i := range UtilList {
		info := UtilizationData{}
		name := strings.Split(UtilList[i], ":")[0]
		Data := strings.Split(UtilList[i], ":")[1]
		info.Name = name
		info.Attri = make(map[string]interface{})
		for _, j := range strings.Split(strings.TrimSpace(Data), " ") {
			parts := strings.Split(string(j), "=")
			res_key := strings.TrimSpace(parts[0])
			res_value := strings.TrimSpace(parts[1])
			info.Attri[res_key] = res_value
		}
		data = append(data, info)
	}
	return data
}

func SetUtilization(data []byte) utils.GeneralResponse {
	var result utils.GeneralResponse
	if len(data) == 0 {
		result.Action = false
		result.Error = gettext.Gettext("No input data")
		return result
	}

	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		result.Action = false
		result.Error = gettext.Gettext("Cannot convert data to json map")
		return result
	}

	utype := jsonData["type"]
	name := jsonData["name"]
	cmd := "pcs " + utype + " utilization " + name + " "
	for k, v := range jsonData {
		if k != "name" && k != "type" {
			cmd = cmd + k + "=" + v
		}
	}
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = gettext.Gettext("Utilization set success")

	} else {
		result.Action = false
		if strings.Contains(string(out), "Error: Unable to find") {
			if utype == "node" {
				result.Error = gettext.Gettext("Node name error")
			} else if utype == "resource" {
				result.Error = gettext.Gettext("Resource name error")
			}
		} else {
			result.Error = string(out)
		}
	}
	return result
}

func DelUtilization(data []byte) utils.GeneralResponse {
	var result utils.GeneralResponse
	if len(data) == 0 {
		result.Action = false
		result.Error = gettext.Gettext("No input data")
		return result
	}

	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		result.Action = false
		result.Error = gettext.Gettext("Cannot convert data to json map")
		return result
	}

	utype := jsonData["type"]
	name := jsonData["name"]
	cmd := "pcs " + utype + " utilization " + name + " "
	for k := range jsonData {
		if k != "name" && k != "type" {
			cmd = cmd + k + "="
		}
	}
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = gettext.Gettext("Delete Utilization success")
	} else {
		result.Action = false
		result.Error = string(out)
	}
	return result
}
