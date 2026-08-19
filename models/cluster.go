/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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

func UpdateClusterProperties(newProp map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

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

	enableList := []string{
		"node-health-green", "stonith-enabled", "symmetric-cluster",
		"maintenance-mode", "node-health-yellow", 
		"no-quorum-policy", "node-health-red", "node-health-strategy",
		"default-resource-stickiness", "start-failure-is-fatal",
		"stop-all-resources", "priority-fencing-delay", "placement-strategy"}
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
			slog.Error(fmt.Sprintf("run command %s failed: %s", cmd, err))
			goto ret
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(out); err != nil {
			slog.Error(fmt.Sprintf("parse xml failed: %s", err.Error()))
			goto ret
		}

		for _, e := range doc.FindElements("//parameters/parameter") {
			prop := getClusterPropertyFromXml(e)
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

				if prop["type"] == "select" {
					propContent["values"] = prop["select"]
					delete(prop, "select")
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

	// // special for getting resource-stickiness property
	// result["resource-stickiness"] = map[string]interface{}{
	// 	"name":    "resource-stickiness",
	// 	"enabled": 1,
	// 	"value":   strconv.Itoa(getResourceStickiness()),
	// 	"content": map[string]string{
	// 		"default": "0",
	// 		"type":    "integer",
	// 		"unit":    "",
	// 	},
	// 	"shortdesc": "",
	// 	"longdesc":  "",
	// }

ret:
	return result, nil
}

var getClusterProperties = func() (map[string]interface{}, error) {
	clusterProperties := map[string]interface{}{}
	var doc *etree.Document
	var nvParis []*etree.Element

	out, err := utils.RunCommand(utils.CmdQueryCrmConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("get cluster properties failed: %s", err))
		goto ret
	}

	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		slog.Error(fmt.Sprintf("parse xml config error: %s", err))
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

	if prop["type"] == "select" {
		propEnums := []string{}
		if prop["longdesc"] != "" {
			values := strings.Split(prop["longdesc"].(string), "Allowed values:")
			if len(values) == 2 {
				propEnums = strings.Split(values[1], ",")
				// select中的值由于空格、换行符导致propEnums识别参数的默认值没有在select中，导致重复添加
				for i, s := range propEnums {
					propEnums[i] = strings.Join(strings.Fields(s), "")
				}
				prop["longdesc"] = values[0]
			}
		}
		if !utils.IsInSlice(prop["default"].(string), propEnums) {
			propEnums = append(propEnums, prop["default"].(string))
		}

		prop["select"] = propEnums
	}

	if prop["longdesc"] == prop["shortdesc"] {
		prop["longdesc"] = ""
	}

	return prop
}

func OperationClusterAction(action string) map[string]interface{} {
	result := map[string]interface{}{}
	if action == "" {
		result["action"] = false
		result["error"] = gettext.Gettext("Action on node failed")
		return result
	}

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

	result["action"] = true
	result["info"] = gettext.Gettext("Action on node success")
	return result
}

var getResourceStickiness = func() int {
	cmdStr := utils.CmdDefaultResourceStickness
	out, err := utils.RunCommand(cmdStr)
	if err != nil {
		slog.Error(fmt.Sprintf("get resource-stickiness failed: %s", err.Error()))
		return 0
	}

	// resource-stickiness=100
	outStr := strings.Split(string(out), "\n")[0]
	valueStr := strings.Split(outStr, "=")[1]
	value, _ := strconv.Atoi(valueStr)

	return value
}
