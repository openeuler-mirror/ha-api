/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
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
	// "github.com/beego/beego/v2/core/logs"
	"gitee.com/openeuler/ha-api/utils"
)

type UtilizationData struct{
	Name string `json:"name"`
	Attri map[string]interface{}  `json:"attri"`
}

type GetUtilResult struct{
	Action bool `json:"action"`
	Data *UtilResponseData `json:"data"`
}

type UtilResponseData struct{
	NodeUtilization []UtilizationData `json:"NodeUtilization"`
	ResUtilization []UtilizationData `json:"ResUtilization"`
 }


type UtilizationResult struct {
    Action bool   `json:"action"`
    Error  string `json:"error,omitempty"`
    Info   string `json:"info,omitempty"`
}


func GetUtilization() GetUtilResult {
	result := GetUtilResult{}
	nodeUtilization := GetOneTypeUtilization("node")
	resUtilization := GetOneTypeUtilization("resource")
	result.Action = true
	result.Data = &UtilResponseData{
		NodeUtilization : nodeUtilization,
		ResUtilization : resUtilization,
	}

	return result
}

func GetOneTypeUtilization(Uti_type string) []UtilizationData {
	data := []UtilizationData{}
	cmd := "pcs " + Uti_type + " utilization"
	output, _ := utils.RunCommand(cmd)
	Util := strings.Split(string(output), "Utilization:")[1]
	
	Util = strings.TrimSpace(Util)
	if Util == ""{
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

func SetUtilization(data []byte) UtilizationResult {
	var result UtilizationResult
	if data == nil || len(data) == 0 {
        result.Action = false
        result.Error = "No input data"
        return result
    }

	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
        result.Action = false
        result.Error = "Cannot convert data to json map"
        return result
    }

	utype, _ := jsonData["type"]
	name, _ := jsonData["name"]
	cmd := "pcs " + utype + " utilization " + name + " "
	for k, v := range jsonData {
		if k != "name" && k != "type" {
			cmd = cmd + k + "=" + v
		}
	}
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = "利用率设置成功！"

	} else {
		result.Action = false
		if strings.Contains(string(out), "Error: Unable to find") {
			if utype == "node" {
				result.Error = "节点名称填写错误"
			} else if utype == "resource" {
				result.Error = "资源名称填写错误"
			}
		} else {
			result.Error = string(out)
		}
	}
	return result
}

func DelUtilization(data []byte) UtilizationResult {
	var result UtilizationResult
	if data == nil || len(data) == 0 {
        result.Action = false
        result.Error = "No input data"
        return result
    }

	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
        result.Action = false
        result.Error = "Cannot convert data to json map"
        return result
    }
	
	utype, _ := jsonData["type"]
	name, _ := jsonData["name"]
	cmd := "pcs " + utype + " utilization " + name + " "
	for k, _ := range jsonData {
		if k != "name" && k != "type" {
			cmd = cmd + k + "="
		}
	}
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = "利用率删除成功！"
	} else {
		result.Action = false
		result.Error = string(out)
	}
	return result
}
