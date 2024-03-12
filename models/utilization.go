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

	"gitee.com/openeuler/ha-api/utils"
)

func GetUtilization() map[string]interface{} {
	result := map[string]interface{}{}
	ul := make(map[string]interface{})

	ul["NodeUtilization"] = GetOneTypeUtilization("node")
	ul["ResUtilization"] = GetOneTypeUtilization("resource")
	result["action"] = true
	result["data"] = ul
	return result
}

func GetOneTypeUtilization(Uti_type string) []map[string]interface{} {
	data := []map[string]interface{}{}
	cmd := "pcs " + Uti_type + " utilization"
	output, _ := utils.RunCommand(cmd)
	Util := strings.Split(string(output), "Utilization:")[1]
	info := map[string]interface{}{}
	if Util == "" {
		return nil
	} else {
		attri := map[string]string{}
		UtilList := strings.Split(Util, "\n")
		if UtilList[1] == "" {
			return nil
		}
		Data := strings.Split(UtilList[1], ":")
		for _, j := range strings.Split(strings.TrimSpace(Data[1]), " ") {
			parts := strings.Split(string(j), "=")
			res_key := strings.TrimSpace(parts[0])
			res_value := strings.TrimSpace(parts[1])
			attri[res_key] = res_value
		}
		info["attri"] = attri
		// for k,OneUtil := range UtilList{

		data = append(data, info)
	}
	return data
}

func SetUtilization(data []byte) map[string]interface{} {
	if data == nil || len(data) == 0 {
		return map[string]interface{}{"action": false, "error": "No input data"}
	}
	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return map[string]interface{}{"action": false, "error": "Cannot convert data to json map"}
	}

	result := map[string]interface{}{}
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
		result["action"] = true
		result["info"] = "利用率设置成功！"

	} else {
		result["action"] = false
		if strings.Contains(string(out), "Error: Unable to find") {
			if utype == "node" {
				result["error"] = "节点名称填写错误"
			} else if utype == "resource" {
				result["error"] = "资源名称填写错误"
			}
		} else {
			result["error"] = out
		}
	}
	return result
}

func DelUtilization(data []byte) map[string]interface{} {
	if data == nil || len(data) == 0 {
		return map[string]interface{}{"action": false, "error": "No input data"}
	}
	jsonData := map[string]string{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return map[string]interface{}{"action": false, "error": "Cannot convert data to json map"}
	}

	result := map[string]interface{}{}
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
		result["action"] = true
		result["info"] = "利用率删除成功！"
	} else {
		result["action"] = false
		result["error"] = out
	}
	return result
}
