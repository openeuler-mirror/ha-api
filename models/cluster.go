package models

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"

	"openkylin.com/ha-api/settings"
	"openkylin.com/ha-api/utils"
)

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

func UpdateClusterProperties(newProp map[string]string) map[string]interface{} {
	result := map[string]interface{}{}

	logs.Debug(newProp)
	if len(newProp) == 0 {
		result["action"] = false
		result["error"] = "No input data"
		return result
	}

	for k, v := range newProp {
		cmdStr := "crm_attribute -t crm_config -n " + k + " -v " + v
		out, err := utils.RunCommand(cmdStr)
		if err != nil {
			result["action"] = true
			result["error"] = string(out)
			return result
		}
	}

	result["action"] = true
	result["info"] = "Update cluster properties Success"
	return result
}

// GetClusterStatus returns crm_mon running status, 0 if normal, -1 if any error
func GetClusterStatus() int {
	_, err := utils.RunCommand("crm_mon - 1")
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

	enableList := []string{"node-health-green", "stonith-enabled", "symmetric-cluster",
		"maintenance-mode", "node-health-yellow", "node-health-yellow",
		"no-quorum-policy", "node-health-red", "node-health-strategy",
		"default-resource-stickiness", "start-failure-is-fatal"}
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

		for _, e := range doc.FindElements("./parameters/parameter") {
			prop := getClusterPropertyFromXml(e)
			logs.Debug(prop)
			name := prop["name"].(string)
			if utils.IsInSlice(name, enableList) {
				if _, ok := clusterProperties[name]; ok {
					prop["value"] = clusterProperties[name]
				} else {
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
					if name == "default-resource-stickiness" {
						prop["value"] = 0
					}
				}
				propContent := make(map[string]interface{})
				propContent["default"] = prop["default"]
				propContent["type"] = prop["type"]
				if prop["type"] == "enum" {
					propContent["values"] = prop["enum"]
				}
				prop["content"] = propContent
				prop["enabled"] = 1
				result[name] = prop
			}
		}
	}

ret:
	return result, nil
}

func getClusterProperties() (map[string]interface{}, error) {
	var clusterProperties map[string]interface{}
	var doc *etree.Document
	var nvParis []*etree.Element

	out, err := utils.RunCommand("cibadmin --query --scope crm_config")
	if err == nil {
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
		var propEnums []string
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
