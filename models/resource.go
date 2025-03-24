/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

// GetResourceInfo
func GetResourceInfo() map[string]interface{} {
	result := make(map[string]interface{})
	clusterStatus := GetClusterStatus()
	if clusterStatus != 0 {
		result["action"] = true
		result["data"] = []string{}
		return result
	}

	constraints := GetAllConstraints()
	if _, ok := constraints["action"]; !ok {
		return constraints
	}

	// reAllRscStatus := map[string]interface{}{}
	for _, constraint := range constraints["data"].([]map[string]interface{}) {
		// TODO : check constraint really modified
		rscId := constraint["id"].(string)
		migrateResources := GetAllMigrateResources()
		if utils.IsInSlice(rscId, migrateResources) {
			constraint["allow_unmigrate"] = true
		} else {
			constraint["allow_unmigrate"] = false
		}

		t := GetResourceType(rscId)
		subRscs := GetSubResources(rscId)
		if t == "group" || t == "clone" {
			constraint["subrscs"] = subRscs["subrscs"]
		} else {
			constraint["svc"] = GetResourceSvc(rscId)
		}
		constraint["type"] = t
	}

	return constraints
}

func GetResourceCategory(rscID string) string {
	ct := ""
	cmd_str := "crm_resource --resource " + rscID + " --query-xml"
	out, err := utils.RunCommand(cmd_str)
	if err != nil {
		return ""
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return ""
	}
	ct = doc.Root().Tag
	return ct
}

func GetResourceType(rscID string) string {
	cmd := fmt.Sprintf(utils.CmdQueryResourcesById, rscID)
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return ""
	}

	typeStr := strings.TrimSpace(string(out))
	rscType := strings.Replace(strings.Split(typeStr, " ")[0], "<", "", -1)

	return rscType
}

// TODO needs to integrate to func GetResourceByConstraintAndId
// or func GetAllConstraints??
func GetResourceConstraints(rscID, relation string) (map[string]interface{}, error) {
	retData := make(map[string]interface{})

	cmd := utils.CmdQueryConstraints
	out, err := utils.RunCommand(cmd)

	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return nil, err
	}
	root := doc.SelectElement("constraints")
	switch relation {
	case "location":
		resourceLocations := []map[string]string{}
		for _, resourceLocation := range root.FindElements("./rsc_location") {
			rsc := resourceLocation.SelectAttr("rsc").Value
			if rsc == rscID {
				rscConstraint := map[string]string{}
				score := resourceLocation.SelectAttrValue("score", "")
				if score == "-INFINITY" || score == "-infinity" || score == "INFINITY" || score == "infinity" {
					continue
				}
				rscConstraint["node"] = resourceLocation.SelectAttrValue("node", "")
				rscConstraint["level"] = getLevelFromScore(score)
				resourceLocations = append(resourceLocations, rscConstraint)
			}

		}
		retData["node_level"] = resourceLocations
		retData["rsc_id"] = rscID
	case "colocation":
		sameNodes := []string{}
		diffNodes := []string{}
		for _, colocation := range root.FindElements("./rsc_colocation") {
			rsc := colocation.SelectAttr("rsc").Value
			rscWith := colocation.SelectAttr("with-rsc").Value
			score := colocation.SelectAttrValue("score", "")

			if (rsc == rscID && score == "INFINITY") || (rscWith == rscID && score == "INFINITY") {
				sameNodes = append(sameNodes, getOtherRsc(rsc, rscWith))
			} else if (rsc == rscID && score == "-INFINITY") || (rscWith == rscID && score == "-INFINITY") {
				diffNodes = append(diffNodes, getOtherRsc(rsc, rscWith))
			}
		}
		retData["same_node"] = sameNodes
		retData["rsc_id"] = rscID
		retData["diff_node"] = diffNodes
	case "order":
		before := []string{}
		after := []string{}

		for _, order := range root.FindElements("rsc_order") {
			first := order.SelectAttrValue("first", "")
			then := order.SelectAttr("then").Value

			if first == rscID {
				after = append(after, then)
			} else if then == rscID {
				before = append(before, first)
			}
		}
		logs.Debug(before)
		logs.Debug(after)
		retData["before_rscs"] = before
		retData["rsc_id"] = rscID
		retData["after_rscs"] = after
	}
	return retData, nil
}

func getOtherRsc(rsc, rscWith string) string {
	if rsc == "" {
		return rscWith
	}
	return rsc
}

func getLevelFromScore(score string) string {
	switch score {
	case "20000":
		return "Master Node"
	case "16000":
		return "Slave 1"
	case "15000":
		return "Slave 2"
	case "14000":
		return "Slave 3"
	case "13000":
		return "Slave 4"
	default:
		return ""
	}
}

func GetResourceFailedMessage() map[string]map[string]string {
	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	failInfo := map[string]map[string]string{}
	if err != nil {
		return failInfo
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return failInfo
	}
	failures := doc.FindElements("/crm_mon/failures/failure")
	if len(failures) == 0 {
		return failInfo
	} else {
		for _, failure := range failures {
			infoFail := map[string]string{}
			rscIdf := strings.Split(failure.SelectAttr("op_key").Value, "_stop_")[0]
			rscIdm := strings.Split(rscIdf, "_start_")[0]
			rscId := strings.Split(rscIdm, "_start_")[0]
			node := failure.SelectAttrValue("node", "")
			exitreason := failure.SelectAttr("exitreason").Value
			infoFail["node"] = node
			infoFail["exitreason"] = exitreason
			failInfo[rscId] = infoFail
		}
	}
	return failInfo
}

func GetResourceMetaAttributes(category string) map[string]interface{} {
	retjson := make(map[string](map[string]interface{}))

	retjson["target-role"] = make(map[string]interface{})
	retjson["target-role"]["content"] = make(map[string]interface{})
	retjson["target-role"]["name"] = "target-role"
	retjson["target-role"]["content"].(map[string]interface{})["values"] = []string{"Stopped", "Started"}
	retjson["target-role"]["content"].(map[string]interface{})["default"] = "Stopped"
	retjson["target-role"]["content"].(map[string]interface{})["type"] = "enum"
	retjson["target-role"]["content"].(map[string]interface{})["dec"] = "What state should the cluster attempt to keep this resource in?"

	retjson["priority"] = make(map[string]interface{})
	retjson["priority"]["content"] = make(map[string]interface{})
	retjson["priority"]["name"] = "priority"
	retjson["priority"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["priority"]["content"].(map[string]interface{})["dec"] = "If not all resources can be active, the cluster will stop lower priority resources in order to keep higher priority ones active."

	retjson["is-managed"] = make(map[string]interface{})
	retjson["is-managed"]["content"] = make(map[string]interface{})
	retjson["is-managed"]["name"] = "is-managed"
	retjson["is-managed"]["content"].(map[string]interface{})["type"] = "boolean"
	retjson["is-managed"]["content"].(map[string]interface{})["dec"] = "Is the cluster allowed to start and stop the resource?"

	if category == "group" {
		return map[string]interface{}{
			"action": true,
			"data":   retjson,
		}
	}

	retjson["resource-stickiness"] = make(map[string]interface{})
	retjson["resource-stickiness"]["content"] = make(map[string]interface{})
	retjson["resource-stickiness"]["name"] = "resource-stickiness"
	retjson["resource-stickiness"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["resource-stickiness"]["content"].(map[string]interface{})["dec"] = "How much does the resource prefer to stay where it is? Defaults to the value of \"default-resource-stickiness\""

	retjson["migration-threshold"] = make(map[string]interface{})
	retjson["migration-threshold"]["content"] = make(map[string]interface{})
	retjson["migration-threshold"]["name"] = "migration-threshold"
	retjson["migration-threshold"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["migration-threshold"]["content"].(map[string]interface{})["dec"] = "How many failures should occur for this resource on a node before making the node ineligible to host this resource. Default: \"none\""

	retjson["multiple-active"] = make(map[string]interface{})
	retjson["multiple-active"]["content"] = make(map[string]interface{})
	retjson["multiple-active"]["name"] = "multiple-active"
	retjson["multiple-active"]["content"].(map[string]interface{})["values"] = []string{"stop_start", "stop_only", "block"}
	retjson["multiple-active"]["content"].(map[string]interface{})["type"] = "enum"
	retjson["multiple-active"]["content"].(map[string]interface{})["dec"] = "What should the cluster do if it ever finds the resource active on more than one node."

	retjson["failure-timeout"] = make(map[string]interface{})
	retjson["failure-timeout"]["content"] = make(map[string]interface{})
	retjson["failure-timeout"]["name"] = "failure-timeout"
	retjson["failure-timeout"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["failure-timeout"]["content"].(map[string]interface{})["dec"] = "How many seconds to wait before acting as if the failure had not occurred (and potentially allowing the resource back to the node on which it failed. Default: \"never\""

	retjson["allow-migrate"] = make(map[string]interface{})
	retjson["allow-migrate"]["content"] = make(map[string]interface{})
	retjson["allow-migrate"]["name"] = "allow-migrate"
	retjson["allow-migrate"]["content"].(map[string]interface{})["type"] = "boolean"
	retjson["allow-migrate"]["content"].(map[string]interface{})["dec"] = "Allow resource migration for resources which support migrate_to/migrate_from actions."

	if category == "primitive" {
		return map[string]interface{}{
			"action": true,
			"data":   retjson,
		}
	}

	if category == "clone" {
		retjson["interleave"] = make(map[string]interface{})
		retjson["interleave"]["content"] = make(map[string]interface{})
		retjson["interleave"]["name"] = "interleave"
		retjson["interleave"]["content"].(map[string]interface{})["type"] = "boolean"
		retjson["interleave"]["content"].(map[string]interface{})["dec"] = "Changes the behavior of ordering constraints (between clones/masters) so that instances can start/stop as soon as their peer instance has (rather than waiting for every instance of the other clone has)."

		retjson["clone-max"] = make(map[string]interface{})
		retjson["clone-max"]["content"] = make(map[string]interface{})
		retjson["clone-max"]["name"] = "clone-max"
		retjson["clone-max"]["content"].(map[string]interface{})["type"] = "integer"
		retjson["clone-max"]["content"].(map[string]interface{})["dec"] = "How many copies of the resource to start. Defaults to the number of nodes in the cluster."

		retjson["promoted-max"] = make(map[string]interface{})
		retjson["promoted-max"]["content"] = make(map[string]interface{})
		retjson["promoted-max"]["name"] = "promoted-max"
		retjson["promoted-max"]["content"].(map[string]interface{})["type"] = "integer"
		retjson["promoted-max"]["content"].(map[string]interface{})["dec"] = "If promotable is true, the number of instances that can be promoted at one time across the entire cluster"

		retjson["promotable"] = make(map[string]interface{})
		retjson["promotable"]["content"] = make(map[string]interface{})
		retjson["promotable"]["name"] = "promotable"
		retjson["promotable"]["content"].(map[string]interface{})["type"] = "boolean"
		retjson["promotable"]["content"].(map[string]interface{})["desc"] = "If true, clone instances can perform a special role that Pacemaker will manage via the resource agent's promote and demote actions. The resource agent must support these actions"

		retjson["notify"] = make(map[string]interface{})
		retjson["notify"]["content"] = make(map[string]interface{})
		retjson["notify"]["name"] = "notify"
		retjson["notify"]["content"].(map[string]interface{})["type"] = "boolean"
		retjson["notify"]["content"].(map[string]interface{})["desc"] = "Call the resource agent's notify action for all active instances, before and after starting or stopping any clone instance"

		return map[string]interface{}{
			"action": true,
			"data":   retjson,
		}

	}
	return map[string]interface{}{
		"action": true,
		"data":   retjson,
	}
}

func GetResourceByConstraintAndId() {

}

func CreateResource(data []byte) map[string]interface{} {
	if len(data) == 0 {
		return map[string]interface{}{"action": false, "error": gettext.Gettext("No input data")}
	}
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return map[string]interface{}{"action": false, "error": gettext.Gettext("Cannot convert data to json map")}
	}
	jsonMap := jsonData

	var rscId string
	if v, ok := jsonMap["id"].(string); ok {
		rscId = v
	} else if v, ok := jsonMap["id"].(int); ok {
		rscId = strconv.Itoa(v)
	}
	cate := jsonMap["category"].(string)
	if cate == "primitive" {
		rscIdStr := " id=\"" + rscId + "\""
		rscClass := " class=\"" + jsonMap["class"].(string) + "\""
		rscType := " type=\"" + jsonMap["type"].(string) + "\""
		cmd := "cibadmin --create -o resources --xml-text '<"
		role := utils.CmdResourceStop
		out, err := utils.RunCommand(utils.CmdQueryCIB)
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(out); err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		primitives := doc.FindElements("primitive")
		for _, primitive := range primitives {
			id := primitive.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		groups := doc.FindElements("group")
		for _, group := range groups {
			id := group.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		clones := doc.FindElements("clone")
		for _, clone := range clones {
			id := clone.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		if _, ok := jsonMap["provider"]; ok {
			provider := " provider=\"" + jsonMap["provider"].(string) + "\""
			cmd = cmd + cate + rscIdStr + rscClass + rscType + provider + ">'"
		} else {
			cmd = cmd + cate + rscIdStr + rscClass + rscType + ">'"
		}
		_, err = utils.RunCommand(cmd)
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if s == "target-role" {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": err.Error()}
			}
		}
		if _, ok := jsonMap["provider"]; !ok {
			instance := fmt.Sprintf(utils.CmdCrmResource, rscId)
			class := jsonMap["class"]
			if class == "stonith" {
				cmdStr := instance + " --set-parameter pcmk_host_check --parameter-value static-list"
				_, err := utils.RunCommand(cmdStr)
				if err != nil {
					return map[string]interface{}{"action": false, "error": err.Error()}
				}
			}
		}
	} else if cate == "group" {
		/*
			{
				"category": "group",
				"id":"tomcat_group",
				"rscs":[
						"tomcat6",
						"tomcat7"
				],
							"meta_attributes":{
					"target-role":"Stopped"
				}
			}
		*/
		rscsArr := jsonMap["rscs"].([]interface{})
		rscs := make([]string, len(rscsArr))
		for ix, v := range rscsArr {
			rscs[ix] = v.(string)
		}
		role := utils.CmdResourceStop
		for _, rsc := range rscs {
			DeletePriAttrib(rsc)
		}
		rscId := jsonMap["id"].(string)
		cmdStr := fmt.Sprintf(utils.CmdResourceGroupAdd, rscId)
		for _, r := range rscs {
			cmdStr = cmdStr + " " + r

		}
		out, err := utils.RunCommand(utils.CmdQueryCIB)
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(out); err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		primitives := doc.FindElements("primitive")
		for _, primitive := range primitives {
			id := primitive.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		groups := doc.FindElements("group")
		for _, group := range groups {
			id := group.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		clones := doc.FindElements("clone")
		for _, clone := range clones {
			id := clone.SelectAttr("id").Value
			resourceId := map[string]string{"id": id}
			if rscId == resourceId["id"] {
				return map[string]interface{}{"action": false, "error": rscId + " is exist"}
			}
		}
		_, err = utils.RunCommand(cmdStr)
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if s == "target-role" {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": err.Error()}
			}
		}
	} else if cate == "clone" {
		/*
			{
				"category": "clone",
				"id":"test5",
				"rsc_id":"test4",
				"meta_attributes":{
					"target-role":"Stopped"
				}
			}
			id is unused
		*/
		oriId := jsonMap["rsc_id"].(string)
		ids := getResourceConstraintIDs(oriId, "location")
		for _, item := range ids {
			cmd := fmt.Sprintf(utils.CmdLocationDelete, item)
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": err.Error()}
			}
		}
		role := utils.CmdResourceStop
		cmdStr := fmt.Sprintf(utils.CmdResourceClone, oriId)
		DeleteCloneAttrib(oriId)
		_, err := utils.RunCommand(cmdStr)
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if s == "target-role" {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": err.Error()}
			}
		}
	}
	if _, exists := jsonMap["selfFlag"].(string); exists {
		var emptySlice []byte
		ResourceAction(rscId, "start", emptySlice)
	}
	return map[string]interface{}{"action": true, "info": "Add " + cate + " resource success"}

}

// UpdateResourceAttributes updates meta_attributes and instance_attributes
func UpdateResourceAttributes(rscId string, data map[string]interface{}) error {
	/*
		Example data:
		{
			"category": "primitive",
			"actions":[
					{
						"interval":"100",
						"name":"start"
					}
				],
			"meta_attributes":{
				"resource-stickiness":"104",
				"is-managed":"true",
				"target-role":"Started"
			},
			"type":"Filesystem",
			"id":"iscisi",
			"provider":"heartbeat",
			"instance_attributes":{
				"device":"/dev/sda1",
				"directory":"/var/lib/mysql",
				"fstype":"ext4"
			},
			"class":"ocf"
		}
	*/
	if len(data) == 0 {
		return errors.New(gettext.Gettext("No input data"))
	}
	// delete all the attribute
	attrib := GetMetaAndInst(rscId)
	if _, ok := attrib["meta_attributes"]; ok {
		metaAttri := attrib["meta_attributes"]
		for _, v := range metaAttri {
			cmd := "crm_resource -r " + rscId + " -m --delete-parameter " + v
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return err
			}
		}
	}
	if _, ok := attrib["instance_attributes"]; ok {
		metaAttri := attrib["instance_attributes"]
		for _, v := range metaAttri {
			cmd := "crm_resource -r " + rscId + " --delete-parameter " + v
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return err
			}
		}
	}
	if data["category"] == "group" {
		if _, ok := data["meta_attributes"]; ok {
			metaAttri := data["meta_attributes"].(map[string]interface{})
			for k, v := range metaAttri {
				value := ""
				if t, ok := v.(string); ok {
					value = t
				} else if t, ok := v.(bool); ok {
					if t {
						value = "true"
					} else {
						value = "false"
					}
				} else if t, ok := v.(float64); ok {
					m := int(t)
					value = strconv.Itoa(m)
				} else {
					logs.Error("unparsed value: ", v)
				}
				cmd := fmt.Sprintf(utils.CmdResourceMetaAdd, rscId, k, value)
				_, err := utils.RunCommand(cmd)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if _, ok := data["meta_attributes"]; ok {
			metaAttri := data["meta_attributes"].(map[string]interface{})
			for k, v := range metaAttri {
				value := ""
				if t, ok := v.(string); ok {
					value = t
				} else if t, ok := v.(bool); ok {
					if t {
						value = "true"
					} else {
						value = "false"
					}
				} else if t, ok := v.(float64); ok {
					m := int(t)
					value = strconv.Itoa(m)
				} else {
					logs.Error("unparsed value: ", v)
				}
				cmd := fmt.Sprintf(utils.CmdResourceUpdateMetaForce, rscId, k, value)
				_, err := utils.RunCommand(cmd)
				if err != nil {
					return err
				}
			}
		}
	}

	_, err := utils.RunCommand("sleep 1")
	if err != nil {
		return err
	}

	instStr := ""
	if _, ok := data["instance_attributes"]; ok {
		instAttri := data["instance_attributes"].(map[string]interface{})
		for k, v := range instAttri {
			instStr = instStr + k + "=" + v.(string) + " "
		}

		if instStr != "" {
			_, err = utils.RunCommand(fmt.Sprintf(utils.CmdResourceUpdateForce, rscId, instStr))
			if err != nil {
				return err
			}
		}
	}

	// change operation
	if _, ok := data["actions"]; ok {
		// delete all the attribute
		opList := GetAllOps(rscId)
		if len(opList) != 0 {
			cmdDelHead := fmt.Sprintf(utils.CmdResourceOpDelete, rscId)
			for _, op := range opList {
				cmdDel := cmdDelHead + " " + op
				_, err = utils.RunCommand(cmdDel)
				if err != nil {
					return err
				}
			}
		}
		action := data["actions"].([]interface{})
		// overwrite
		if len(action) > 0 {
			cmdIn := "pcs resource update " + rscId + " op"
			for _, b := range action {
				ops := b.(map[string]interface{})
				name := ops["name"].(string)
				cmdIn = cmdIn + " " + name
				if v, ok := ops["interval"]; ok {
					cmdIn = cmdIn + " " + "interval=" + v.(string)
				}
				if v, ok := ops["start-delay"]; ok {
					cmdIn = cmdIn + " " + "start-delay=" + v.(string)
				}
				if v, ok := ops["timeout"]; ok {
					cmdIn = cmdIn + " " + "timeout=" + v.(string)
				}
				if v, ok := ops["role"]; ok {
					cmdIn = cmdIn + " " + "role=" + v.(string)
				}
				if v, ok := ops["requires"]; ok {
					cmdIn = cmdIn + " " + "requires=" + v.(string)
				}
				if v, ok := ops["on-fail"]; ok {
					cmdIn = cmdIn + " " + "on-fail=" + v.(string)
				}
			}
			_, err = utils.RunCommand(cmdIn)
			if err != nil {
				return err
			}
		}
	}

	// add group resource
	if data["category"] == "group" {
		/*
			{
				"id":"group1",
				"category":"group",
				"rscs":["iscisi", "test1" ],
				"meta_attributes":{
					"target-role":"Started"
				}
			}
		*/
		t := data["rscs"].([]interface{})
		rscs := []string{}
		for _, item := range t {
			rscs = append(rscs, item.(string))
		}
		rscsExist, _ := getGroupRscs(rscId)

		allRscs := append(rscs, rscsExist...)
		for _, v := range allRscs {
			if utils.IsInSlice(v, rscs) {
				if !utils.IsInSlice(v, rscsExist) {
					// 在新配置中但是不在原配置中，需要新增
					DeletePriAttrib(v)
					cmdAdd := fmt.Sprintf(utils.CmdResourceGroupAdd, rscId+" "+v)
					if _, err := utils.RunCommand(cmdAdd); err != nil {
						return err
					}
				}
			} else {
				if utils.IsInSlice(v, rscsExist) {
					// 不在新配置中但是在原配置中，需要删除
					cmdRmv := fmt.Sprintf(utils.CmdResourceGroupRemove, rscId+" "+v)
					if _, err := utils.RunCommand(cmdRmv); err != nil {

						return err
					}
				}
			}
		}
	}
	return nil
}

func GetAllOps(rscId string) []string {
	opList := []string{}
	cmd := fmt.Sprintf(utils.CmdQueryResourceAsXml, rscId)
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return opList
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return opList
	}
	e := doc.FindElement("//operations")
	if e != nil {
		op := e.SelectElements("op")
		for _, item := range op {
			opList = append(opList, item.SelectAttrValue("name", ""))
		}
	}
	return opList
}

func DeletePriAttrib(rscId string) error {
	// delete attribute
	attrib := GetMetaAndInst(rscId)
	if metaAttri, ok := attrib["meta_attributes"]; ok {
		metaArr := metaAttri
		for _, v := range metaArr {
			cmd := "crm_resource -r " + rscId + " -m --delete-parameter "
			if v == "is-managed" || v == "priority" || v == "target-role" {
				cmd += v
				_, err := utils.RunCommand(cmd)
				if err != nil {
					return err
				}
			}
		}
	}
	// delete constraint
	// colocation
	targetId := getResourceConstraintIDs(rscId, "colocation")
	err := DeleteColocationByIdAndAction(rscId, targetId)
	if err != nil {
		return err
	}
	// location
	ids := getResourceConstraintIDs(rscId, "location")
	for _, item := range ids {
		cmd := fmt.Sprintf(utils.CmdLocationDelete, item)
		_, err := utils.RunCommand(cmd)
		if err != nil {
			return err
		}
	}
	// order
	if findOrder(rscId) {
		cmd := fmt.Sprintf(utils.CmdOrderDelete, rscId)
		_, err := utils.RunCommand(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteCloneAttrib(rscId string) error {
	// delete attribute
	attrib := GetMetaAndInst(rscId)
	if metaAttri, ok := attrib["meta_attributes"]; ok {
		metaArr := metaAttri
		cmd := "crm_resource -r " + rscId + " -m --delete-parameter "
		for _, v := range metaArr {
			cmdStr := cmd + v
			_, err := utils.RunCommand(cmdStr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetMetaAndInst(rscId string) map[string][]string {
	cmdStr := fmt.Sprintf(utils.CmdQueryResourceAsXml, rscId)
	out, err := utils.RunCommand(cmdStr)
	if err != nil {
		// return map[string]interface{}{"action": false, "error": err}
		return map[string][]string{}
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		// return map[string]interface{}{"action": false, "error": err}
		return map[string][]string{}
	}
	data := map[string][]string{}
	eMeta := doc.FindElement("//meta_attributes")
	if eMeta != nil {
		prop := []string{}
		items := eMeta.SelectElements("nvpair")
		for _, item := range items {
			prop = append(prop, item.SelectAttr("name").Value)
		}
		data["meta_attributes"] = prop
	}
	eInst := doc.FindElement("//instance_attributes")
	if eInst != nil {
		prop := []string{}
		items := eInst.SelectElements("nvpair")
		for _, item := range items {
			prop = append(prop, item.SelectAttr("name").Value)
		}
		data["instance_attributes"] = prop
	}
	return data
}

func GetAllConstraints() map[string]interface{} {
	rscStatus := GetAllResourceStatus()
	data := map[string](map[string]interface{}){}

	topRsc := GetTopResource()
	for _, rscId := range topRsc {
		rsc := strings.Split(rscId, ":")[0]
		data[rsc] = map[string]interface{}{}
		if _, ok := rscStatus[rsc]; !ok {
			data[rsc]["status"] = "Running"
			data[rsc]["status_message"] = ""
			data[rsc]["running_node"] = []string{}
		} else {
			rscInfo := rscStatus[rsc]
			if rscInfo["isMs"] != nil {
				if rscInfo["isMs"] == true {
					data[rsc]["isMs"] = true
				} else {
					data[rsc]["isMs"] = false
				}
			}
			data[rsc]["status"] = "Running"
			data[rsc]["status_message"] = ""
			data[rsc]["running_node"] = rscStatus[rsc]["running_node"]
		}
		data[rsc]["before_rscs"] = []map[string]string{}
		data[rsc]["after_rscs"] = []map[string]string{}
		data[rsc]["same_node"] = []map[string]string{}
		data[rsc]["diff_node"] = []map[string]string{}
		data[rsc]["location"] = []map[string]string{}
	}

	out, err := utils.RunCommand(utils.CmdQueryCIB)
	if err != nil {
		result := map[string]interface{}{}
		result["action"] = false
		result["data"] = data
		return result
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		result := map[string]interface{}{}
		result["action"] = false
		result["data"] = data
		return result
	}

	constraints := doc.FindElement("/cib/configuration/constraints")
	if constraints != nil {
		//location
		if locations := constraints.FindElements("rsc_location"); locations != nil {
			for _, location := range locations {
				if strings.HasPrefix(location.SelectAttr("id").Value, "cli-prefer-") {
					continue
				}
				node := location.SelectAttrValue("node", "")
				rscId := location.SelectAttrValue("rsc", "")
				score := location.SelectAttrValue("score", "")
				locationSingle := make(map[string]string)
				locationSingle["node"] = node
				locationSingle["level"] = ScoreToLevel(score)
				if data[rscId]["location"] != nil {
					// may resource is not in top resource but constraint rules remains
					data[rscId]["location"] = append(data[rscId]["location"].([]map[string]string), locationSingle)
				}
			}
		}

		//order
		if orders := constraints.FindElements("rsc_order"); orders != nil {
			for _, order := range orders {
				first := order.SelectAttr("first").Value
				then := order.SelectAttr("then").Value

				//try except
				score := order.SelectAttrValue("score", "")
				if score == "" || len(score) == 0 {
					score = "infinity"
				}
				if score != "INFINITY" && score != "+INFINITY" && score != "infinity" && score != "+infinity" {
					continue
				}

				afterRscsArr := []map[string]string{}
				if _, ok := data[first]; !ok {
					afterRscsArr = append(afterRscsArr, map[string]string{"id": then})
				}
				data[first]["after_rscs"] = afterRscsArr
				beforeRscsArr := []map[string]string{}
				if _, ok := data[then]; !ok {
					beforeRscsArr = append(beforeRscsArr, map[string]string{"id": then})
				}
				data[then]["before_rscs"] = beforeRscsArr
			}
		}

		//colocation
		if colocations := constraints.FindElements("rsc_colocation"); colocations != nil {
			for _, colocation := range colocations {
				first := colocation.SelectAttr("rsc").Value
				with := colocation.SelectAttr("with-rsc").Value

				//try except
				score := colocation.SelectAttrValue("score", "")
				if score == "INFINITY" || score == "+INFINITY" || score == "infinity" || score == "+infinity" {
					rsc := map[string]string{}
					rsc["rsc"] = first
					rsc["with_rsc"] = with
					data[first]["same_node"] = rsc
					data[with]["same_node"] = rsc
				} else if score == "-INFINITY" || score == "-infinity" {
					rsc := map[string]string{}
					rsc["rsc"] = first
					rsc["with_rsc"] = with
					data[first]["diff_node"] = rsc
					data[with]["diff_node"] = rsc
				}
			}
		}
	}

	failureInfo := GetResourceFailedMessage()

	constraintMaps := []map[string]interface{}{}
	for rscId := range data {
		constraint := map[string]interface{}{}
		constraint["id"] = rscId
		rscIDFirst := strings.Split(rscId, ":")[0]
		if _, ok := rscStatus[rscId]; !ok {
			if _, ok := failureInfo[rscIDFirst]; !ok {
				constraint["status"] = "Stopped"
				constraint["status_message"] = ""
				constraint["running_node"] = []string{}
			} else if strings.HasSuffix(rscId, "-clone") {
				constraint["status"] = "Failed"
				constraint["status_message"] = failureInfo[rscIDFirst]["exitreason"] + " on " + failureInfo[rscIDFirst]["node"]
				constraint["running_node"] = []string{}
			}
		} else {
			rscInfo := rscStatus[rscId]
			if rscInfo["isMs"] != nil {
				if rscInfo["isMs"] == true {
					constraint["isMs"] = true
				} else {
					constraint["isMs"] = false
				}
			}
			constraint["status_message"] = ""
			constraint["status"] = rscInfo["status"]
			constraint["running_node"] = rscInfo["running_node"]
		}

		colocation := map[string]interface{}{}
		colocation["same_node"] = data[rscId]["same_node"]
		colocation["diff_node"] = data[rscId]["diff_node"]
		if tempArray, ok := colocation["same_node"].([]map[string]string); ok {
			colocation["same_node_num"] = len(tempArray)
		}
		if tempArray, ok := colocation["diff_node"].([]map[string]string); ok {
			colocation["diff_node_num"] = len(tempArray)
		}
		order := map[string]interface{}{}
		order["before_rscs"] = data[rscId]["before_rscs"]
		order["after_rscs"] = data[rscId]["after_rscs"]
		if tempArray, ok := colocation["before_rscs"].([]map[string]string); ok {
			colocation["before_rscs_num"] = len(tempArray)
		}
		if tempArray, ok := colocation["after_rscs"].([]map[string]string); ok {
			colocation["after_rscs_num"] = len(tempArray)
		}
		constraint["location"] = data[rscId]["location"]
		constraint["colocation"] = colocation
		constraint["order"] = order
		constraintMaps = append(constraintMaps, constraint)
	}

	result := map[string]interface{}{}
	result["action"] = true
	result["data"] = constraintMaps
	return result
}

func GetAllMigrateResources() []string {
	result := make([]string, 1)

	cmd := utils.CmdQueryCIB
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return result
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return result
	}

	resourceLocations := make(map[string]interface{})
	// TODO: check real xml document here
	for _, resourceLocation := range doc.FindElements("/cib/configuration/constraints/rsc_location") {
		id := resourceLocation.SelectAttr("id").Value
		resourceLocations[id] = resourceLocation
	}

	migrateIds := map[string]interface{}{}
	for id := range resourceLocations {
		// prefixs := []string{"cli-prefer-", "cli-standby-"}
		for _, prefix := range []string{"cli-prefer-", "cli-standby-"} {
			if strings.HasPrefix(id, prefix) {
				splitId := strings.Split(id, prefix)
				if len(splitId) > 1 {
					rscId := splitId[1]
					if _, ok := migrateIds[rscId]; !ok {
						migrateIds[rscId] = []string{}
					}
					migrateIds[rscId] = append(migrateIds[rscId].([]string), id)
				}
			}
		}
	}
	rscList := []string{}
	if len(migrateIds) != 0 {
		for key := range migrateIds {
			rscList = append(rscList, key)
		}
	}

	return rscList
}

func GetAllResourceStatus() map[string]map[string]interface{} {
	/*
		infos = {
			0:_('Running'),
			1:_('Not Running'),
			2:_('Unmanaged'),
			3:_('Failed'),
			4:_('Stop Failed'),
			5:_('running (Master)'),
			6:_('running (Slave)')}
			rsc_info = {
				"hj1": {"status": 0 , "status_message": "test", running_node: []}
				"hj2": {"status": 0 , "status_message": "test", running_node: []}
		}
	*/
	rscInfo := map[string]map[string]interface{}{}
	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	if err != nil {
		return map[string]map[string]interface{}{}
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return map[string]map[string]interface{}{}
	}

	if len(doc.FindElements("/crm_mon/resources")) == 0 {
		return map[string]map[string]interface{}{}
	}

	rscClone := doc.FindElements("/crm_mon/resources/clone")
	rscGroup := doc.FindElements("/crm_mon/resources/group")
	rscResource := doc.FindElements("/crm_mon/resources/resource")
	if len(rscClone) != 0 {

		// several clone
		for _, rsc := range rscClone {
			// subResources is common resources
			if subRscs := rsc.SelectElements("resource"); len(subRscs) != 0 {
				index := 0
				isMs := false
				cloneRunNodes := []string{}
				cloneInfo := map[string]interface{}{}
				for _, subRsc := range subRscs {
					info := map[string]interface{}{}
					info["status"] = GetResourceStatus(subRsc)
					if subRsc.SelectAttr("role").Value == "Slave" || subRsc.SelectAttr("role").Value == "Master" {
						isMs = true
					}
					info["status_message"] = ""
					nodename := ""
					if node := subRsc.FindElement("node"); node != nil {
						nodename = node.SelectAttr("name").Value
					}
					id := subRsc.SelectAttr("id").Value + ":" + strconv.Itoa(index)
					index++
					info["running_node"] = []string{nodename}
					rscInfo[id] = info
					cloneRunNodes = append(cloneRunNodes, nodename)
				}
				if isMs {
					cloneInfo["isMs"] = true
				} else {
					cloneInfo["isMs"] = false
				}
				cloneInfo["status"] = GetResourceStatus(rsc)
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = cloneRunNodes
				cloneId := rsc.SelectAttr("id").Value
				rscInfo[cloneId] = cloneInfo
			}
			// subResources is gourp resources
			if subRscs := rsc.SelectElements("group"); len(subRscs) != 0 {
				cloneRunNodes := []string{}
				cloneInfo := map[string]interface{}{}
				for _, subRsc := range subRscs {
					subRscId := subRsc.SelectAttr("id").Value
					groupInfo := map[string]interface{}{}
					groupRunNodes := []string{}
					if innerRscs := subRsc.SelectElements("resource"); len(innerRscs) != 0 {
						groupInfo["status"] = "Not Running"
						if false {
							innerRsc := innerRscs[0]
							innerRscId := innerRsc.SelectAttr("id").Value
							info := map[string]interface{}{}
							info["status"] = GetResourceStatus(innerRsc)
							info["status_message"] = ""
							if node := innerRsc.FindElement("node"); node != nil {
								nodename := ""
								if node := innerRsc.FindElement("node"); node != nil {
									nodename = node.SelectAttr("name").Value
								}
								info["running_node"] = []string{nodename}
								groupRunNodes = append(groupRunNodes, nodename)
								cloneRunNodes = append(cloneRunNodes, nodename)
								groupInfo["status"] = "Running"
							}
							fatherId := strings.Split(subRscId, ":")[1]
							id := innerRscId + ":" + fatherId
							rscInfo[id] = info
						} else {
							for _, innerRsc := range innerRscs {
								innerRscId := innerRsc.SelectAttr("id").Value
								info := map[string]interface{}{}
								groupInfo["status"] = "Stopped"
								info["status"] = GetResourceStatus(innerRsc)
								if info["status"] == "Running" {
									groupInfo["status"] = "Running"
								}
								info["status_message"] = ""
								if node := innerRsc.FindElement("node"); node != nil {
									nodename := node.SelectAttr("name").Value
									info["running_node"] = []string{nodename}
									cloneRunNodes = append(cloneRunNodes, nodename)
									groupRunNodes = append(groupRunNodes, nodename)
								}
								fatherId := strings.Split(subRscId, ":")[1]
								id := innerRscId + ":" + fatherId
								rscInfo[id] = info
							}
						}
					}
					groupInfo["running_node"] = utils.RemoveDupl(groupRunNodes)
					groupInfo["status_message"] = ""
					groupId := subRscId
					rscInfo[groupId] = groupInfo
				}
				cloneInfo["status"] = GetResourceStatus(rsc)
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = utils.RemoveDupl(cloneRunNodes)
				cloneInfo["isMs"] = false
				cloneId := rsc.SelectAttr("id").Value
				rscInfo[cloneId] = cloneInfo
			}
		}
	}

	if len(rscGroup) != 0 {
		// several group
		if len(rscGroup) != 1 {
			for _, rsc := range rscGroup {
				subRscs := rsc.SelectElements("resource")
				// several resources in each group
				if len(subRscs) > 1 {
					groupRunNodes := []string{}
					groupInfo := map[string]interface{}{}
					groupInfo["status"] = "Not Running"
					for _, subRsc := range subRscs {
						info := map[string]interface{}{}
						info["status"] = GetResourceStatus(subRsc)
						info["status_message"] = ""
						nodename := ""
						if node := subRsc.FindElement("node"); node != nil {
							nodename = node.SelectAttr("name").Value
							info["running_node"] = []string{nodename}
							groupRunNodes = append(groupRunNodes, nodename)
							groupInfo["status"] = "Running"
						}
						id := subRsc.SelectAttr("id").Value
						rscInfo[id] = info
					}
					groupInfo["status_message"] = ""
					groupInfo["running_node"] = groupRunNodes
					id := rsc.SelectAttr("id").Value
					rscInfo[id] = groupInfo
				} else {
					subRsc := subRscs[0]
					groupRunNodes := []string{}
					groupInfo := map[string]interface{}{}
					groupInfo["status"] = "Not Running"
					info := map[string]interface{}{}
					info["status"] = GetResourceStatus(subRsc)
					info["status_message"] = ""
					nodename := ""
					if node := subRsc.FindElement("node"); node != nil {
						nodename = node.SelectAttr("name").Value
						info["running_node"] = []string{nodename}
						groupRunNodes = append(groupRunNodes, nodename)
						groupInfo["status"] = "Running"
					}
					id := subRsc.SelectAttr("id").Value
					rscInfo[id] = info
					groupInfo["status_message"] = ""
					groupInfo["running_node"] = utils.RemoveDupl(groupRunNodes)
					gid := rsc.SelectAttr("id").Value
					rscInfo[gid] = groupInfo
				}
			}
		} else { // single group
			rscGroupSin := rscGroup[0]
			subRscs := rscGroupSin.SelectElements("resource")
			groupInfo := map[string]interface{}{}
			flag := 0
			if len(subRscs) > 1 {
				groupInfo["status"] = "Not Running"
				for _, subRsc := range subRscs {
					info := map[string]interface{}{}
					info["status"] = GetResourceStatus(subRsc)
					info["status_message"] = ""
					id := subRsc.SelectAttr("id").Value
					if nodes := subRsc.SelectElements("node"); len(nodes) != 0 {
						infoRunNode := []string{}
						groupRunNodes := []string{}
						for _, node := range nodes {
							nodename := node.SelectAttr("name").Value
							infoRunNode = append(infoRunNode, nodename)
							groupRunNodes = append(groupRunNodes, nodename)
							groupInfo["status"] = "Running"
						}
						info["running_node"] = utils.RemoveDupl(infoRunNode)
						groupInfo["running_node"] = utils.RemoveDupl(groupRunNodes)
					}
					rscInfo[id] = info
				}
				groupInfo["status_message"] = ""
				groupId := rscGroupSin.SelectAttr("id").Value
				rscInfo[groupId] = groupInfo
			} else {
				subRsc := subRscs[0]
				info := map[string]interface{}{}
				info["status"] = GetResourceStatus(subRsc)
				info["status_message"] = ""
				nodename := ""
				if node := subRsc.FindElement("node"); node != nil {
					nodename = node.SelectAttr("name").Value
				} else {
					flag = 1
				}
				id := subRsc.SelectAttr("id").Value
				info["running_node"] = []string{nodename}
				rscInfo[id] = info
				groupInfo["status_message"] = ""
				groupInfo["running_node"] = []string{}
				if flag == 1 {
					groupInfo["status"] = "Not Running"
				} else {
					groupInfo["status"] = "Running"
					groupInfo["running_node"] = []string{nodename}
				}
				groupId := rscGroupSin.SelectAttr("id").Value
				rscInfo[groupId] = groupInfo
			}
		}
	}
	if len(rscResource) != 0 {
		// several common resource
		for _, rsc := range rscResource {
			resourceInfo := map[string]interface{}{}
			resourceInfo["status"] = GetResourceStatus(rsc)
			runningNode := []string{}
			if nodes := rsc.SelectElements("node"); len(nodes) != 0 {
				for _, node := range nodes {
					runningNode = append(runningNode, node.SelectAttr("name").Value)
				}
			}
			resourceInfo["running_node"] = runningNode
			resourceInfo["status_message"] = ""
			id := rsc.SelectAttr("id").Value
			rscInfo[id] = resourceInfo
		}
	}

	return rscInfo
}

func GetResourceStatus(rscInfo *etree.Element) string {
	rscId := rscInfo.SelectAttr("id").Value
	failInfo := GetResourceFailedMessage()
	if _, ok := failInfo[rscId]; ok {
		return "Failed"
	}

	if rscInfo.SelectAttr("managed").Value == "false" {
		return "Unmanaged"
	}
	if rscInfo.SelectAttr("failed").Value == "true" {
		return "Failed"
	}
	if role := rscInfo.SelectAttr("role"); role != nil {
		if role.Value == "Started" {
			return "Running"
		}
		if role.Value == "Stopped" {
			return "Not Running"
		}

	}
	return "Running"
}

func GetSubResources(rscId string) map[string]interface{} {
	rscStatus := GetAllResourceStatus()
	failInfo := GetResourceFailedMessage() // failure run information
	rscInfo := map[string]interface{}{}

	out, err := utils.RunCommand(utils.CmdQueryResources)
	if err != nil {
		return rscInfo
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return rscInfo
	}
	resJson := doc.FindElement("resources")
	rscType := GetResourceType(rscId)
	rscInfo["id"] = rscId
	subRscs := []map[string]interface{}{}

	if rscType == "clone" {
		nodeInfo, _ := GetNodesInfo()
		nodeNum := len(nodeInfo)
		clone := resJson.FindElements("clone")
		var cloneAim *etree.Element
		// find the clone resource's location
		if len(clone) < 2 {
			if clone[0].SelectAttr("id").Value == rscId {
				cloneAim = clone[0]
			}
		} else {
			for _, subClone := range clone {
				if subClone.SelectAttr("id").Value == rscId {
					cloneAim = subClone
				}
			}
		}
		// clone resource objects is group or primitive
		// Parse the types of resources and assemble them into strings
		if cloneGroup := cloneAim.FindElement("group"); cloneGroup != nil {
			subRscId := cloneGroup.SelectAttr("id").Value
			if rscPrimitive := cloneGroup.FindElements("primitive"); len(rscPrimitive) != 0 {
				for i := 0; i < nodeNum; i++ {
					subRsc := map[string]interface{}{}
					subRsc["status"] = "Not Running"
					subRsc["running_node"] = []string{}
					subRsc["status_message"] = ""
					subId := subRscId + ":" + strconv.Itoa(i)
					subRsc["id"] = subId
					if _, ok := rscStatus[subId]; ok {
						subRsc["status"] = rscStatus[subId]["status"]
						subRsc["running_node"] = rscStatus[subId]["running_node"]
					}
					subRsc["type"] = "group"
					subSubRsc := []map[string]interface{}{}
					for _, primitive := range rscPrimitive {
						pid := primitive.SelectAttr("id").Value + ":" + strconv.Itoa(i)
						primitiveInfo := map[string]interface{}{}
						primitiveInfo["id"] = pid
						primitiveInfo["status"] = "Not Running"
						primitiveInfo["status_message"] = ""
						primitiveInfo["running_node"] = []string{}
						if _, ok := rscStatus[pid]; ok {
							primitiveInfo["status"] = rscStatus[pid]["status"]
							primitiveInfo["running_node"] = rscStatus[pid]["running_node"]
						}
						primitiveInfo["svc"] = primitive.SelectAttr("type").Value
						primitiveInfo["type"] = "primitive"
						subSubRsc = append(subSubRsc, primitiveInfo)
					}
					subRsc["subrscs"] = subSubRsc
					subRscs = append(subRscs, subRsc)
				}
			}
		}
		if clonePrimitive := cloneAim.FindElement("primitive"); clonePrimitive != nil {
			subRscId := clonePrimitive.SelectAttr("id").Value
			for i := 0; i < nodeNum; i++ {
				subRsc := map[string]interface{}{}
				subRsc["status"] = "Not Running"
				subRsc["running_node"] = []string{}
				subRsc["status_message"] = ""
				subId := subRscId + ":" + strconv.Itoa(i)
				subRsc["id"] = subId
				if _, ok := rscStatus[subId]; ok {
					subRsc["status"] = rscStatus[subId]["status"]
					subRsc["running_node"] = rscStatus[subId]["running_node"]
				}
				// judgment failure message
				if subFailInfo, ok := failInfo[subRscId]; ok {
					subRsc["status_message"] = subFailInfo["exitreason"] + " on " + subFailInfo["node"]
					subRsc["status"] = "Failed"
					if subRscStatus, ok := rscStatus[subRscId]; ok {
						subRsc["running_node"] = subRscStatus["running_node"]
					} else {
						subRsc["running_node"] = []string{}
					}
				}
				subRsc["type"] = "primitive"
				subRsc["svc"] = clonePrimitive.SelectAttr("type").Value
				subRscs = append(subRscs, subRsc)
			}
		}
	}

	if rscType == "group" {
		groups := resJson.FindElements("group")
		for _, group := range groups {
			if rscId == group.SelectAttr("id").Value {
				rscPrimitive := group.FindElements("primitive")
				for _, primitive := range rscPrimitive {
					primitiveInfo := map[string]interface{}{}
					primitiveInfo["status_message"] = ""
					primitiveId := primitive.SelectAttr("id").Value
					primitiveInfo["id"] = primitiveId
					if priRscStatus, ok := rscStatus[primitiveId]; ok {
						primitiveInfo["status"] = priRscStatus["status"]
						if _, ok := priRscStatus["running_node"]; ok {
							if _, ok := priRscStatus["running_node"]; ok {
								primitiveInfo["running_node"] = priRscStatus["running_node"]
							} else {
								primitiveInfo["running_node"] = []string{}
							}
						} else {
							primitiveInfo["running_node"] = []string{}
						}
					} else {
						primitiveInfo["status"] = ""
						primitiveInfo["running_node"] = []string{}
					}
					// judgment failure message
					if priFail, ok := failInfo[primitiveId]; ok {
						primitiveInfo["status_message"] = priFail["exitreason"] + " on " + priFail["node"]
						primitiveInfo["status"] = "Failed"
						if priRscStatus, ok := rscStatus[primitiveId]; ok {
							if _, ok := priRscStatus["running_node"]; ok {
								primitiveInfo["running_node"] = priRscStatus["running_node"]
							} else {
								primitiveInfo["running_node"] = []string{}
							}
						} else {
							primitiveInfo["running_node"] = []string{}
						}
					}
					primitiveInfo["type"] = "primitive"
					primitiveInfo["svc"] = primitive.SelectAttr("type").Value
					subRscs = append(subRscs, primitiveInfo)
				}
			}
		}
	}
	rscInfo["subrscs"] = subRscs
	return rscInfo
}

func GetResourceSvc(rscId string) string {
	cmd := fmt.Sprintf(utils.CmdQueryResourceAsXml, rscId)
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return ""
	}
	// Provide compatibility with different versions of Corosync
	xmlIndex := strings.Index(string(out), "XML:")
	if xmlIndex == -1 {
		xmlIndex = strings.Index(string(out), "xml:")
	}
	if xmlIndex == -1 {
		return ""
	}
	xmlStr := string(out)[xmlIndex+len("XML:"):]

	doc := etree.NewDocument()
	if err = doc.ReadFromString(xmlStr); err != nil {
		return ""
	}
	rscType := doc.FindElement("primitive").SelectAttrValue("type", "")

	return rscType
}

func GetTopResource() []string {
	result := []string{}

	out, err := utils.RunCommand(utils.CmdQueryResources)
	if err != nil {
		return result
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return result
	}

	elements := doc.FindElements("/resources/clone")
	for _, element := range elements {
		result = append(result, element.SelectAttr("id").Value)
	}
	elements = doc.FindElements("/resources/group")
	for _, element := range elements {
		result = append(result, element.SelectAttr("id").Value)
	}
	elements = doc.FindElements("/resources/primitive")
	for _, element := range elements {
		result = append(result, element.SelectAttr("id").Value)
	}

	return result
}

func ResourceAction(rscID, action string, data []byte) error {
	// in case ":" within the resource name
	rscID = strings.Split(rscID, ":")[0]
	// cmd := "crm_resource --resource "
	switch action {
	case "start":
		cmd := utils.CmdResourceStart + rscID
		_, err := utils.RunCommand(cmd)
		return err
	case "stop":
		cmd := utils.CmdResourceStop + rscID
		_, err := utils.RunCommand(cmd)
		return err
	case "delete":
		var cmd string
		category := GetResourceCategory(rscID)
		if category == "clone" {
			cmd = fmt.Sprintf(utils.CmdResourceDeleteForce, rscID[:len(rscID)-6])
		} else {
			// not clone
			cmd = fmt.Sprintf(utils.CmdResourceDeleteForce, rscID)
		}
		_, err := utils.RunCommand(cmd)
		return err
	case "cleanup":
		cmd := fmt.Sprintf(utils.CmdCrmResource, rscID) + " --cleanup"
		_, err := utils.RunCommand(cmd)
		return err
	case "unclone":
		cmd := utils.CmdResourceUnclone + rscID
		_, err := utils.RunCommand(cmd)
		return err
	case "ungroup":
		cmd := utils.CmdResourceUngroup + rscID
		_, err := utils.RunCommand(cmd)
		return err
	// desperated
	case "migrate":
		d := struct {
			IsForce bool   `json:"is_force"`
			ToNode  string `json:"to_node"`
			Period  string `json:"period"`
		}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return errors.New("invalid json data")
		}

		cmd := fmt.Sprintf(utils.CmdCrmResource, rscID) + " --move -N " + d.ToNode
		out, err := utils.RunCommand(cmd)
		if err != nil {
			if string(out) == "Error performing operation: Situation already as requested" {
				return errors.New("The resource " + rscID + " is running on node " + d.ToNode + " already!")
			}
		}
		return err
	// desperated
	case "unmigrate":
		cmd := utils.CmdQueryConstraints
		out, err := utils.RunCommand(cmd)
		if err != nil {
			return err
		}
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(out); err != nil {
			return err
		}
		rscNames := doc.FindElements("/constraints/rsc_location")
		for _, item := range rscNames {
			rsc := item.SelectAttrValue("rsc", "")
			if rscID == rsc {
				locationID := item.SelectAttrValue("id", "")
				cmd2 := fmt.Sprintf(utils.CmdLocationDelete, locationID)
				if _, err := utils.RunCommand(cmd2); err != nil {
					return err
				}
			}
		}
		logs.Info("Unmigrate resource not found")
		return nil
	case "location":
		// location:
		// {"node_level": [{"node": "ns187", "level": "Master Node"},
		// {"node": "ns188", "level": "Slave 1"}]}
		ids := getResourceConstraintIDs(rscID, action)
		for _, item := range ids {
			cmd := fmt.Sprintf(utils.CmdLocationDelete, item)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}

		d := map[string]interface{}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return err
		}
		for _, item := range d["node_level"].([]interface{}) {
			var score int
			mapItem := item.(map[string]interface{})
			if mapItem["level"] == "Master Node" {
				score = 20000
			} else if mapItem["level"] == "Slave 1" {
				score = 16000
			} else if mapItem["level"] == "Slave 2" {
				score = 15000
			} else if mapItem["level"] == "Slave 3" {
				score = 14000
			} else if mapItem["level"] == "Slave 4" {
				score = 13000
			}
			node := mapItem["node"].(string)
			cmd := fmt.Sprintf(utils.CmdLocationAdd, rscID, node, strconv.Itoa(score))
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
	case "colocation":
		// 	colocation:
		// {"same_node": ["test1234"],"diff_node": ["group_tomcat"]}
		ids := getResourceConstraintIDs(rscID, action)
		if err := DeleteColocationByIdAndAction(rscID, ids); err != nil {
			return err
		}

		d := struct {
			SameNode []string `json:"same_node"`
			DiffNode []string `json:"diff_node"`
		}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return err
		}
		for _, item := range d.SameNode {
			cmd := fmt.Sprintf(utils.CmdColationAdd, rscID, item)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		for _, item := range d.DiffNode {
			cmd := fmt.Sprintf(utils.CmdColationAdd, rscID, item)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
	case "order":
		if findOrder(rscID) {
			cmd := fmt.Sprintf(utils.CmdOrderDelete, rscID)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		d := struct {
			BeforeRscs []string `json:"before_rscs"`
			AfterRscs  []string `json:"after_rscs"`
		}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return err
		}
		for _, item := range d.BeforeRscs {
			cmd := fmt.Sprintf(utils.CmdOrderAdd, item, rscID)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		for _, item := range d.AfterRscs {
			cmd := fmt.Sprintf(utils.CmdOrderAdd, rscID, item)
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func getResourceConstraintIDs(rscID, action string) []string {
	ids := []string{}
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return ids
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return ids
	}

	if action == "colocation" {
		et := doc.FindElements("/constraints/rsc_colocation")
		for _, item := range et {
			rsc := item.SelectAttrValue("rsc", "")
			rscWith := item.SelectAttrValue("with-rsc", "")
			if rsc == rscID {
				ids = append(ids, rscWith)
			}
			if rscWith == rscID {
				ids = append(ids, rsc)
			}
		}
		return ids
	} else if action == "location" {
		et := doc.FindElements("/constraints/rsc_location")
		for _, item := range et {
			rsc := item.SelectAttrValue("rsc", "")
			if rsc == rscID {
				if item.SelectAttr("score") != nil && item.SelectAttrValue("score", "") == "-INFINITY" {
					continue
				}
				ids = append(ids, item.SelectAttrValue("id", ""))
			}
		}
		return ids
	}
	return ids
}

func DeleteColocationByIdAndAction(rscID string, targetIds []string) error {
	for _, item := range targetIds {
		cmd := fmt.Sprintf(utils.CmdColationDelete, rscID, item)
		if _, err := utils.RunCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func findOrder(rscID string) bool {
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return false
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return false
	}
	et := doc.FindElements("/constraints/rsc_order")
	for _, item := range et {
		first := item.SelectAttrValue("first", "")
		then := item.SelectAttrValue("then", "")
		if first == rscID || then == rscID {
			return true
		}
	}
	return false
}

func GetResourceInfoByrscID(rscID string) (interface{}, error) {
	cmd := fmt.Sprintf(utils.CmdQueryResourceAsXml, rscID)
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return nil, err
	}

	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		return nil, err
	}

	ct := doc.Root().Tag
	result, err := GetResourceInfoID(ct, xml)
	if err != nil {
		return nil, err
	}
	result["id"] = string(rscID)
	result["category"] = string(ct)

	if _, ok := result["provider"]; ok {
		if result["provider"] == "" {
			delete(result, "provider")
		}
	}

	return result, nil
}

func GetResourceInfoID(ct, xmlData string) (map[string]interface{}, error) {
	doc := etree.NewDocument()
	doc.ReadFromString(xmlData)
	data := map[string]interface{}{}

	// Format data to map here
	switch ct {
	case "primitive":
		d, err := getResourceInfoFromXml("primitive", doc.Root())
		if err != nil {
			return nil, err
		}
		info := d.(PrimitiveResource)
		data["id"] = info.ID
		data["class"] = info.Class
		data["type"] = info.Type
		data["provider"] = info.Provider

		actions := []map[string]string{}
		for _, ac := range info.Operations {
			m := map[string]string{}
			m["name"] = ac.Name
			m["interval"] = ac.Interval
			m["timeout"] = ac.Timeout
			actions = append(actions, m)
		}
		data["actions"] = actions
	case "group":
		d, err := getResourceInfoFromXml("group", doc.Root())
		if err != nil {
			return nil, err
		}
		info := d.(GroupResource)
		data["id"] = info.ID

		rscs := []string{}
		for _, p := range info.Primitives {
			rscs = append(rscs, p.ID)
		}
		data["rscs"] = rscs
	case "clone":
		d, err := getResourceInfoFromXml("clone", doc.Root())
		if err != nil {
			return nil, err
		}
		info := d.(CloneResource)
		data["id"] = info.ID

		// TODO: check if only one Primitive resource or list
		rscs := []string{}
		for _, p := range info.Primitives {
			rscs = append(rscs, p.ID)
		}
		data["rsc_id"] = rscs
	}

	// For meta_attributes
	e := doc.FindElement("/" + ct + "/meta_attributes")
	if e != nil {
		prop, _ := getResourceInfoFromXml("meta", e)
		if len(prop.(map[string]string)) > 0 {
			data["meta_attributes"] = prop
		}
	}

	//For instance_attributes
	e = doc.FindElement("/" + ct + "/instance_attributes")
	if e != nil {
		prop, _ := getResourceInfoFromXml("inst", e)
		if len(prop.(map[string]string)) > 0 {
			data["instance_attributes"] = prop
		}
	}

	//For actions
	e = doc.FindElement("/" + ct + "/operations")
	if e != nil {
		prop, _ := getResourceInfoFromXml("operations", e)
		if len(prop.([]map[string]string)) > 0 {
			data["actions"] = prop
		}
	}

	return data, nil
}

type NvPair struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Op struct {
	ID       string `json:"id"`
	Interval string `json:"interval"`
	Name     string `json:"name"`
	Timeout  string `json:"timeout"`
}

type PrimitiveResource struct {
	Class      string `json:"class"`
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	Type       string `json:"type"`
	Operations []Op   `json:"operations"`
}

type GroupResource struct {
	ID         string `json:"id"`
	Primitives []PrimitiveResource
}

type CloneResource struct {
	ID         string `json:"id"`
	Primitives []PrimitiveResource
}

// getResourceInfoFromXml returns resource information parsed from xml.
// for meta and inst, returns a map;
// for operations, returns a map slice.
func getResourceInfoFromXml(cl string, et *etree.Element) (interface{}, error) {
	// var prop map[string]interface{}
	if cl == "group" {
		rsc := GroupResource{}
		rsc.ID = et.SelectAttrValue("id", "")

		rsc.Primitives = []PrimitiveResource{}
		els := et.FindElements("primitive")
		for _, e := range els {
			prsc := getPrimitiveResourceInfo(e)
			rsc.Primitives = append(rsc.Primitives, prsc)
		}
		return rsc, nil
	} else if cl == "clone" {
		rsc := CloneResource{}
		rsc.ID = et.SelectAttrValue("id", "")

		rsc.Primitives = []PrimitiveResource{}
		els := et.FindElements("primitive")
		for _, e := range els {
			prsc := getPrimitiveResourceInfo(e)
			rsc.Primitives = append(rsc.Primitives, prsc)
		}
		return rsc, nil
	} else if cl == "primitive" {
		rsc := getPrimitiveResourceInfo(et)
		return rsc, nil
	} else if cl == "meta" || cl == "inst" {
		result := map[string]string{}
		op := et.FindElements("./nvpair")
		for _, item := range op {
			name := item.SelectAttr("name").Value
			value := item.SelectAttr("value").Value
			if value == "True" {
				value = "true"
			}
			if value == "False" {
				value = "false"
			}
			result[name] = value
		}
		return result, nil
	} else if cl == "operations" {
		// var prop = []map[string]string{}
		result := []map[string]string{}
		op := et.FindElements("./op")
		for _, item := range op {
			i := map[string]string{}
			for _, v := range item.Attr {
				i[v.Key] = v.Value
			}
			result = append(result, i)
		}
		return result, nil
	}

	return nil, errors.New("invalid resource type")
}

func getPrimitiveResourceInfo(ele *etree.Element) PrimitiveResource {
	result := PrimitiveResource{}

	result.Class = ele.SelectAttrValue("class", "")
	result.ID = ele.SelectAttrValue("id", "")
	result.Provider = ele.SelectAttrValue("provider", "")
	result.Type = ele.SelectAttrValue("type", "")

	result.Operations = []Op{}
	for _, v := range ele.SelectElements("operations") {
		if r := v.SelectAttr("id"); r != nil {
			op := Op{}
			op.ID = v.SelectAttrValue("id", "")
			op.Interval = v.SelectAttrValue("interval", "")
			op.Name = v.SelectAttrValue("name", "")
			op.Timeout = v.SelectAttrValue("timeout", "")
			result.Operations = append(result.Operations, op)
		}
	}

	return result
}

func getGroupRscs(groupId string) ([]string, error) {
	cmd := fmt.Sprintf(utils.CmdQueryResourceAsXml, groupId)
	out, err := utils.RunCommand(cmd)

	if err != nil {
		// result := map[string]interface{}{}
		// result["action"]=
		// return result
		return nil, err
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return nil, err
	}
	et := doc.FindElements("//primitive")
	rscs := []string{}

	for _, pri := range et {
		rscs = append(rscs, pri.SelectAttrValue("id", ""))
	}

	return rscs, nil

}

func levelInit() []string {
	nodeinfo, _ := GetNodesInfo()
	nodeNum := len(nodeinfo)
	max := nodeNum - 1
	levelScoreArr := make([]string, max)
	for i := 0; i < max; i++ {
		levelScoreArr[i] = strconv.Itoa(16000 - 1000*i)
	}
	return levelScoreArr
}

func ScoreToLevel(score string) string {
	levelScoreArr := levelInit()
	if score == "20000" {
		return "Master Node"
	}
	if score == "-INFINITY" || score == "-infinity" {
		return "No Run Node"
	}

	isIn := false
	for _, v := range levelScoreArr {
		if score == v {
			isIn = true
			break
		}
	}
	if !isIn {
		return score
	}

	level := 1
	for _, s := range levelScoreArr {
		if s == score {
			return "Slave " + strconv.Itoa(level)
		}
		level = level + 1
	}
	return ""
}
