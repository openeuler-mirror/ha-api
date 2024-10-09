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
 * LastEditTime: 2022-04-20 10:30:18
 * Description: 集群相关功能的实现
 ******************************************************************************/
package models

import (
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
)

// TypeToSplit used in cluster global parameters that has a unit
var TypeToSplit = []string{"time", "percentage"}

func GetClusterPropertiesInfo() map[string]interface{} {
	result := make(map[string]interface{})

	clusterData, err := getClusterPropertiesDefinition()
	if err != nil {
		result["action"] = false
		result["error"] = ""
		return result
	}

	data := map[string]interface{}{}
	data["name"] = "Policy Engine"
	data["shortdesc"] = "Policy Engine Options"
	data["version"] = "1.0"
	data["nodecount"] = 2
	data["isconfig"] = true
	data["longdesc"] = "This is a fake resource that details the options that can be configured for the Policy Engine."
	data["parameters"] = clusterData
	result["data"] = data
	result["action"] = true

	return result
}

//
//func CreateCluster(clusterInfo map[string]interface{}) map[string]interface{} {
//	result := map[string]interface{}{}
//	logs.Debug(clusterInfo)
//	if len(clusterInfo) == 0 {
//		result["action"] = false
//		result["error"] = gettext.Gettext("No input data")
//		return result
//	}
//	authRes := hostAuth(clusterInfo)
//	if !authRes["action"] {
//		return authRes
//	} else {
//		nodeName := clusterInfo["node_name"]
//		url := "https://" + fmt.Sprint(nodeName) + ":" + strconv.Itoa(8080) + "/api/v1/managec/local_cluster_info"
//		resp, err := http.Get(url)
//		if resp.StatusCode == 200 {
//
//		}
//	}
//	return result
//
//}

func UpdateClusterProperties(newProp map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	logs.Debug(newProp)
	if len(newProp) == 0 {
		result["action"] = false
		result["error"] = gettext.Gettext("No input data")
		return result
	}

	for key, value := range newProp {
		var strValue string
		if t, ok := value.(string); ok {
			strValue = t
		} else if t, ok := value.(bool); ok {
			if t {
				strValue = "true"
			} else {
				strValue = "false"
			}
		} else if t, ok := value.(float64); ok {
			strValue = strconv.FormatInt(int64(t), 10)
		}

		var cmdStr string
		// special for getting resource-stickiness property
		if key == "resource-stickiness" {
			cmdStr = utils.CmdUpdateResourceStickness + strValue
		} else {
			cmdStr = utils.CmdUpdateCrmConfig + key + " -v " + strValue
		}

		out, err := utils.RunCommand(cmdStr)
		if err != nil {
			result["action"] = false
			result["error"] = string(out)
			return result
		}
	}

	result["action"] = true
	result["info"] = gettext.Gettext("Update cluster properties Success")
	return result
}

// GetClusterStatus returns crm_mon running status, 0 if normal, -1 if any error
func GetClusterStatus() int {
	_, err := utils.RunCommand(utils.CmdClusterStatus)
	if err != nil {
		return -1
	}
	return 0
}

func getClusterPropertiesDefinition() (map[string]interface{}, error) {
	clusterProperties, err := getClusterProperties()
	if err != nil {
		return nil, err
	}

	enableList := []string{"node-health-green", "stonith-enabled",
		"symmetric-cluster", "maintenance-mode", "node-health-yellow",
		"no-quorum-policy", "node-health-red", "node-health-strategy",
		"default-resource-stickiness", "start-failure-is-fatal",
		"stonith-action", "placement-strategy", // new properties
		"cluster-recheck-interval", "load-threshold",
		"node-action-limit", "transition-delay", "stonith-max-attempts",
		"enable-acl", "cluster-ipc-limit", "stop-all-resources",
		"priority-fencing-delay"}
	sources := []map[string]string{
		{
			"name": "pacemaker-schedulerd",
			"path": settings.PacemakerSchedulerd,
		}, {
			"name": "pacemaker-controld",
			"path": settings.PacemakerControld,
		}, {
			"name": "pacemaker-based",
			"path": settings.PacemakerBased,
		},
	}

	result := make(map[string]interface{})
	for _, source := range sources {
		cmd := source["path"] + " metadata "
		out, err := utils.RunCommand(cmd)
		if err != nil {
			logs.Error("run command failed: ", cmd, err)
			goto ret
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(out); err != nil {
			logs.Error("parse command xml failed: ", doc, err)
			goto ret
		}

		for _, e := range doc.FindElements("//parameters/parameter") {
			prop := getClusterPropertyFromXml(e)
			logs.Debug(prop)
			name := prop["name"].(string)
			if utils.IsInSlice(name, enableList) {
				if _, ok := clusterProperties[name]; ok {
					prop["value"] = clusterProperties[name]
				} else {
					// pacemaker-schedulerd
					if name == "node-health-green" {
						prop["value"] = 0
					}
					if name == "stonith-enabled" {
						prop["value"] = "true"
					}
					if name == "symmetric-cluster" {
						prop["value"] = "true"
					}
					if name == "maintenance-mode" {
						prop["value"] = "false"
					}
					if name == "node-health-yellow" {
						prop["value"] = "0"
					}
					if name == "node-health-red" {
						prop["value"] = "0"
					}
					if name == "no-quorum-policy" {
						prop["value"] = "ignore"
					}
					if name == "node-health-strategy" {
						prop["value"] = "none"
					}
					if name == "start-failure-is-fatal" {
						prop["value"] = "true"
					}
					if name == "default-resource-stickiness" { // not required in the current version
						prop["value"] = 0
					}
					if name == "stonith-action" {
						prop["value"] = "reboot"
					}
					if name == "placement-strategy" {
						prop["value"] = "default"
					}
					// pacemaker-controld
					if name == "dc-version" {
						prop["value"] = "none"
					}
					if name == "cluster-name" {
						prop["value"] = "(null)"
					}
					if name == "cluster-recheck-interval" {
						prop["value"] = "15min"
					}
					if name == "load-threshold" {
						prop["value"] = "80%"
					}
					if name == "node-action-limit" {
						prop["value"] = "0"
					}
					if name == "transition-delay" {
						prop["value"] = "0s"
					}
					// if name == "stonith-watchdog-timeout" {
					// 	prop["value"] = "(null)"
					// }
					if name == "stonith-max-attempts" {
						prop["value"] = "10"
					}
					// pacemaker-based
					if name == "enable-acl" {
						prop["value"] = "false"
					}
					if name == "cluster-ipc-limit" {
						prop["value"] = "500"
					}
					if name == "priority-fencing-delay" {
						prop["value"] = "0"
					}
					if name == "stop-all-resources" {
						prop["value"] = "false"
					}
				}
				propContent := make(map[string]interface{})
				propContent["default"] = prop["default"]
				propContent["type"] = prop["type"]

				if prop["type"] == "enum" {
					propContent["values"] = prop["enum"]
					delete(prop, "enum")
				}
				delete(prop, "default")
				delete(prop, "type")

				propContent["unit"] = ""
				propType := propContent["type"].(string)
				if utils.IsInSlice(propType, TypeToSplit) { // split value like 15min, 80%
					prop["value"], _ = utils.GetNumAndUnitFromStr(prop["value"].(string))
					propContent["default"], propContent["unit"] = utils.GetNumAndUnitFromStr(propContent["default"].(string))
				}
				prop["content"] = propContent
				prop["enabled"] = 1
				result[name] = prop
			}
		}
	}

	// special for getting resource-stickiness property
	result["resource-stickiness"] = map[string]interface{}{
		"name":    "resource-stickiness",
		"enabled": 1,
		"value":   strconv.Itoa(getResourceStickiness()),
		"content": map[string]string{
			"default": "0",
			"type":    "integer",
			"unit":    "",
		},
		"shortdesc": "",
		"longdesc":  "",
	}

ret:
	return result, nil
}

func getClusterProperties() (map[string]interface{}, error) {
	clusterProperties := map[string]interface{}{}
	var doc *etree.Document
	var nvParis []*etree.Element

	out, err := utils.RunCommand(utils.CmdQueryCrmConfig)
	if err != nil {
		logs.Error("get cluster properties failed", err)
		goto ret
	}

	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		logs.Error("parse xml config error", err)
		goto ret
	}

	nvParis = doc.FindElements("//nvpair")
	for _, v := range nvParis {
		clusterProperties[v.SelectAttr("name").Value] = v.SelectAttr("value").Value
	}
	return clusterProperties, nil

ret:
	return nil, err
}

func getClusterPropertyFromXml(e *etree.Element) map[string]interface{} {
	prop := map[string]interface{}{
		"name":      e.SelectAttrValue("name", ""),
		"shortdesc": "",
		"longdesc":  "",
	}
	for _, item := range []string{"shortdesc", "longdesc"} {
		if ele := e.SelectElement(item); ele != nil {
			prop[item] = ele.Text()
		}
	}

	content := e.SelectElement("content")
	if content != nil {
		prop["type"] = content.SelectAttrValue("type", "")
		prop["default"] = content.SelectAttrValue("default", "")
	} else {
		prop["type"] = ""
		prop["default"] = ""
	}

	if prop["type"] == "enum" {
		propEnums := []string{}
		if prop["longdesc"] != "" {
			values := strings.Split(prop["longdesc"].(string), "Allowed values:")
			if len(values) == 2 {
				propEnums = strings.Split(values[1], ", ")
				prop["longdesc"] = values[0]
			}
		}
		if !utils.IsInSlice(prop["default"].(string), propEnums) {
			propEnums = append(propEnums, prop["default"].(string))
		}

		prop["enum"] = propEnums
	}

	if prop["longdesc"] == prop["shortdesc"] {
		prop["longdesc"] = ""
	}

	return prop
}

func OperationClusterAction(action string) map[string]interface{} {
	result := map[string]interface{}{}
	if action == "start" {
		utils.RunCommand(utils.CmdStartCluster)
	}
	if action == "stop" {
		utils.RunCommand(utils.CmdStopClusterLocal)
	}
	if action == "restart" {
		utils.RunCommand(utils.CmdStopClusterLocal)
		utils.RunCommand(utils.CmdStartCluster)
	}
	if action == "" {
		result["action"] = false
		result["error"] = gettext.Gettext("Action on node failed")
		return result
	} else {
		result["action"] = true
		result["info"] = gettext.Gettext("Action on node success")
		return result
	}
}

func getResourceStickiness() int {
	cmdStr := utils.CmdDefaultResourceStickness
	out, err := utils.RunCommand(cmdStr)
	if err != nil {
		logs.Error("get resource-stickiness failed: ", err.Error())
		return 0
	}

	// resource-stickiness=100
	outStr := strings.Split(string(out), "\n")[0]
	valueStr := strings.Split(outStr, "=")[1]
	value, _ := strconv.Atoi(valueStr)

	return value
}
