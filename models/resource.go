package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"openkylin.com/ha-api/utils"
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
			constraint["allow_unmigrat"] = true
		} else {
			constraint["allow_unmigrat"] = false
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
	// TODO:
	return ""
}

func GetResourceType(rscID string) string {
	cmd := "cibadmin --query --scope resources|grep 'id=\"" + rscID + "\"'"
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return ""
	}

	typeStr := strings.TrimSpace(string(out))
	rscType := strings.Replace(strings.Split(typeStr, " ")[0], "<", "", -1)

	return rscType
}

//TODO needs to integrate to func GetResourceByConstraintAndId
// or func GetAllConstraints??
func GetResourceConstraints(rscID, relation string) (map[string]interface{}, error) {
	retData := make(map[string]interface{})

	var cmd string
	cmd = "cibadmin --query --scope constraints"
	out, err := utils.RunCommand(cmd)
	logs.Debug(string(out))

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
		var resourceLocations []map[string]string
		for _, resourceLocation := range root.FindElements("./rsc_location") {
			rsc := resourceLocation.SelectAttr("rsc").Value
			if rsc == rscID {
				var rscConstraint map[string]string
				score := resourceLocation.SelectAttr("score").Value
				if score == "-INFINITY" || score == "-infinity" {
					continue
				}
				if score == "INFINITY" || score == "infinity" {
					continue
				}
				rscConstraint["node"] = resourceLocation.SelectAttr("id").Value
				if score == "20000" {
					rscConstraint["level"] = "Master Node"
				} else if score == "16000" {
					//TODO implements func turnScoreToLevel
					rscConstraint["level"] = "Slave 1"
				}
				resourceLocations = append(resourceLocations, rscConstraint)
			}

		}
		retData["node_level"] = resourceLocations
		retData["rsc_id"] = rscID
		break
	case "colocation":
		var sameNodes, diffNodes []string
		for _, colocation := range root.FindElements("./rsc_colocation") {
			rsc := colocation.SelectAttr("rsc").Value
			rscWith := colocation.SelectAttr("with-rsc").Value

			if rsc == rscID {
				score := colocation.SelectAttr("score").Value
				switch score {
				case "INFINITY":
					sameNodes = append(sameNodes, rscWith)
					break
				case "-INFINITY":
					diffNodes = append(diffNodes, rscWith)
					break
				}
			}

			//TODO find better way to solve the rsc and with-rsc
			if rscWith == rscID {
				score := colocation.SelectAttr("score").Value
				switch score {
				case "INFINITY":
					sameNodes = append(sameNodes, rsc)
					break
				case "-INFINITY":
					diffNodes = append(diffNodes, rsc)
					break
				}
			}
		}
		retData["same_node"] = sameNodes
		retData["rsc_id"] = rscID
		retData["diff_node"] = diffNodes
		break
	case "order":
		var before, after []string

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
		break
	}
	return retData, nil
}

func GetResourceFailedMessage() map[string]map[string]string {
	out, err := utils.RunCommand("crm_mon -1 --as-xml")
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
			node := failure.SelectAttr("node").Value
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
	if data == nil || len(data) == 0 {
		return map[string]interface{}{"action": false, "error": "No input data"}
	}
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return map[string]interface{}{"action": false, "error": "Cannot convert data to json map"}
	}
	jsonMap := jsonData.(map[string]interface{})

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
		role := "pcs resource disable "
		out, err := utils.RunCommand("cibadmin -Q")
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(out); err != nil {
			return map[string]interface{}{"action": false, "error": "xml parse failed"}
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
		out, err = utils.RunCommand(cmd)
		if err != nil {
			return map[string]interface{}{"action": false, "error": out}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if "target-role" == s {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			out, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": out}
			}
		}
		if _, ok := jsonMap["provider"]; !ok {
			instance := "crm_resource --resource "
			class := jsonMap["class"]
			if class == "stonith" {
				cmdStr := instance + rscId + " --set-parameter pcmk_host_check --parameter-value static-list"
				out, err := utils.RunCommand(cmdStr)
				if err != nil {
					return map[string]interface{}{"action": false, "error": out}
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
		// rscs := jsonMap["rscs"].([]string)
		rscsArr := jsonMap["rscs"].([]interface{})
		rscs := make([]string, len(rscsArr))
		for ix, v := range rscsArr {
			rscs[ix] = v.(string)
		}
		role := "pcs resource disable "
		for _, rsc := range rscs {
			DeletePriAttrib(rsc)
		}
		rscId := jsonMap["id"].(string)
		cmdStr := "pcs resource group add " + rscId
		for _, r := range rscs {
			cmdStr = cmdStr + " " + r

		}
		out, err := utils.RunCommand("cibadmin -Q")
		if err != nil {
			return map[string]interface{}{"action": false, "error": err.Error()}
		}
		doc := etree.NewDocument()
		if err = doc.ReadFromBytes(out); err != nil {
			return map[string]interface{}{"action": false, "error": "xml parse failed"}
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
		out, err = utils.RunCommand(cmdStr)
		if err != nil {
			return map[string]interface{}{"action": false, "error": out}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if "target-role" == s {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			out, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": out}
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
			cmd := "pcs constraint location delete " + item
			out, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": out}
			}
		}
		role := "pcs resource disable "
		cmdStr := "pcs resource clone " + oriId
		DeleteCloneAttrib(oriId)
		out, err := utils.RunCommand(cmdStr)
		if err != nil {
			return map[string]interface{}{"action": false, "error": out}
		}
		flag := 0
		UpdateResourceAttributes(rscId, jsonMap)
		attrib := GetMetaAndInst(rscId)
		for _, arr := range attrib {
			for _, s := range arr {
				if "target-role" == s {
					flag = 1
					break
				}
			}
		}
		if flag == 0 {
			cmd := role + rscId
			out, err := utils.RunCommand(cmd)
			if err != nil {
				return map[string]interface{}{"action": false, "error": out}
			}
		}
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
	if data == nil || len(data) == 0 {
		return errors.New("No input data")
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
			cmd := "crm_resource -r " + rscId + " -m --delete-parameter " + v
			_, err := utils.RunCommand(cmd)
			if err != nil {
				return err
			}
		}
	}
	if data["category"] == "group" {
		if _, ok := data["meta_attributes"]; ok {
			metaAttri := data["meta_attributes"].(map[string]string)
			for k, v := range metaAttri {
				cmd := "pcs resource meta " + rscId + " " + k + "=" + v
				_, err := utils.RunCommand(cmd)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if _, ok := data["meta_attributes"]; ok {
			metaAttri := data["meta_attributes"].(map[string]string)
			for k, v := range metaAttri {
				cmd := "pcs resource update " + rscId + " meta " + k + "=" + v + " --force"
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
		instAttri := data["instance_attributes"].(map[string]string)
		for k, v := range instAttri {
			instStr = instStr + k + "=" + v + " "
		}
	}
	_, err = utils.RunCommand("pcs resource update " + rscId + " " + instStr + " --force")
	if err != nil {
		return err
	}

	// change operation
	if _, ok := data["actions"]; ok {
		// delete all the attribute
		opList := GetAllOps(rscId)
		if len(opList) != 0 {
			cmdDelHead := "pcs resource op delete " + rscId
			for _, op := range opList {
				cmdDel := cmdDelHead + " " + op
				_, err = utils.RunCommand(cmdDel)
				if err != nil {
					return err
				}
			}
		}
		action := data["actions"].([]map[string]string)
		// overwrite
		cmdIn := "pcs resource update " + rscId + " op"
		for _, ops := range action {
			name := ops["name"]
			cmdIn = cmdIn + " " + name
			if v, ok := ops["interval"]; ok {
				cmdIn = cmdIn + " " + "interval=" + v
			}
			if v, ok := ops["start-delay"]; ok {
				cmdIn = cmdIn + " " + "start-delay=" + v
			}
			if v, ok := ops["timeout"]; ok {
				cmdIn = cmdIn + " " + "timeout=" + v
			}
			if v, ok := ops["role"]; ok {
				cmdIn = cmdIn + " " + "role=" + v
			}
			if v, ok := ops["requires"]; ok {
				cmdIn = cmdIn + " " + "requires=" + v
			}
			if v, ok := ops["on-fail"]; ok {
				cmdIn = cmdIn + " " + "on-fail=" + v
			}
		}
		_, err = utils.RunCommand(cmdIn)
		if err != nil {
			return err
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
		rscs := data["rscs"].([]string)
		rscsExist, _ := getGroupRscs(rscId)
		//  单独对新列表第一项进行判断操作
		//  如果第一项已经存在
		rsc := rscs[0]
		isFound := false
		for i, v := range rscsExist {
			if v == rsc {
				rscsExist = append(rscsExist[:i], rscsExist[i+1:]...)
				isFound = true
			}
		}
		// 在rscs_exist中移出该项
		if isFound == false {
			// 如果第一项不存在，则添加第一项后进行删除和添加操作
			DeletePriAttrib(rsc)

			cmd := "pcs resource group add " + rscId + " " + rsc
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		rscs = rscs[1:]
		// 	针对rscs_exist的删除操作和rscs的添加操作
		//   将rscs_exist中的资源全部删除
		for _, i := range rscsExist {
			cmdRmv := "pcs resource group remove " + rscId + " " + i
			if _, err := utils.RunCommand(cmdRmv); err != nil {

				return err
			}
		}
		cmdAdd := ""
		for _, i := range rscs {
			DeletePriAttrib(i)
			cmdAdd = "pcs resource group add " + rscId + " " + i
			if _, err := utils.RunCommand(cmdAdd); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetAllOps(rscId string) []string {
	opList := []string{}
	cmd := "crm_resource --resource " + rscId + " --query-xml"
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return opList
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return opList
	}
	e := doc.FindElement("operations")
	if e != nil {
		op := e.SelectElements("./op")
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
		cmd := "pcs constraint location delete " + item
		_, err := utils.RunCommand(cmd)
		if err != nil {
			return err
		}
	}
	// order
	if findOrder(rscId) {
		cmd := "pcs constraint order delete " + rscId
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
	cmdStr := "crm_resource --resource " + rscId + " --query-xml"
	out, err := utils.RunCommand(cmdStr)
	if err != nil {
		// return map[string]interface{}{"action": false, "error": out}
		return map[string][]string{}
	}
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		// return map[string]interface{}{"action": false, "error": err}
		return map[string][]string{}
	}
	data := map[string][]string{}
	eMeta := doc.FindElement("meta_attributes")
	if eMeta != nil {
		prop := []string{}
		items := eMeta.SelectElements("./nvpair")
		for _, item := range items {
			prop = append(prop, item.SelectAttr("name").Value)
		}
		data["meta_attributes"] = prop
	}
	eInst := doc.FindElement("instance_attributes")
	if eInst != nil {
		prop := []string{}
		items := eInst.SelectElements("./nvpair")
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

	out, err := utils.RunCommand("cibadmin -Q")
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

	constraints := doc.FindElement("constraints")
	if constraints != nil {
		//location
		if locations := constraints.FindElements("rsc_location"); locations != nil {
			for _, location := range locations {
				if strings.HasPrefix(location.SelectAttr("id").Value, "cli-prefer-") {
					continue
				}
				node := location.SelectAttr("node").Value
				rscId := location.SelectAttr("rsc").Value
				score := location.SelectAttr("score").Value
				locationSingle := make(map[string]string)
				locationSingle["node"] = node
				locationSingle["level"] = ScoreToLevel(score)
				locationArr := []map[string]string{}
				for key := range data {
					if rscId == key {
						locationArr = append(locationArr, locationSingle)
					}
				}
				data[rscId]["location"] = locationArr
			}
		}

		//order
		if orders := constraints.FindElements("rsc_order"); orders != nil {
			for _, order := range orders {
				first := order.SelectAttr("first").Value
				then := order.SelectAttr("then").Value

				//try except
				score := order.SelectAttr("score").Value
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
				score := colocation.SelectAttr("score").Value
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

	cmd := "cibadmin -Q"
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
	for _, resourceLocation := range doc.FindElements("rsc_location") {
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
	out, err := utils.RunCommand("crm_mon -1 --as-xml")
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
		if len(rscClone) != 1 {
			for _, rsc := range rscClone {
				// subResources is common resources
				if subRscs := rsc.SelectElements("resource"); len(subRscs) != 0 {
					index := 0
					cloneRunNodes := []string{}
					cloneInfo := map[string]interface{}{}
					for _, subRsc := range subRscs {
						info := map[string]interface{}{}
						info["status"] = GetResourceStatus(subRsc)
						info["status_message"] = ""
						nodename := ""
						if node := subRsc.FindElement("node"); node != nil {
							nodename = node.SelectAttr("name").Value
						}
						id := subRsc.SelectAttr("id").Value + ":" + string(index)
						index++
						info["running_node"] = []string{nodename}
						rscInfo[id] = info
						cloneRunNodes = append(cloneRunNodes, nodename)
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
							if len(innerRscs) == 1 {
								innerRsc := innerRscs[0]
								innerRscId := innerRsc.SelectAttr("id").Value
								info := map[string]interface{}{}
								info["status"] = GetResourceStatus(subRsc)
								info["status_message"] = ""
								if node := innerRsc.FindElement("node"); node != nil {
									nodename := ""
									if node := subRsc.FindElement("node"); node != nil {
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
									info["status"] = GetResourceStatus(subRsc)
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
					cloneId := rsc.SelectAttr("id").Value
					rscInfo[cloneId] = cloneInfo
				}
			}
		} else { // single clone
			// clone is group resource
			if rscCloneGroups := rscClone[0].SelectElements("group"); len(rscCloneGroups) != 0 {
				cloneRunNodes := []string{}
				cloneInfo := map[string]interface{}{}
				cloneInfo["status"] = "Not Running"
				for _, rscCloneGroup := range rscCloneGroups {
					cloneGroupRunNodes := []string{}
					cloneGroupInfo := map[string]interface{}{}
					subGroupId := rscCloneGroup.SelectAttr("id").Value
					index := strings.Split(subGroupId, ":")[1]
					if subRscs := rscCloneGroup.SelectElements("resource"); len(subRscs) != 0 {
						cloneGroupInfo["status"] = "Not Running"

						for _, subRsc := range subRscs {
							info := map[string]interface{}{}
							info["status"] = GetResourceStatus(subRsc)
							info["status_message"] = ""
							id := subRsc.SelectAttr("id").Value + ":" + index
							if node := subRsc.FindElement("node"); node != nil {
								nodename := node.SelectAttr("name").Value
								info["running_node"] = []string{nodename}
								cloneGroupInfo["status"] = "Running"
								cloneInfo["status"] = "Running"
								cloneGroupRunNodes = append(cloneGroupRunNodes, nodename)
								cloneRunNodes = append(cloneRunNodes, nodename)
							}
							rscInfo[id] = info
						}
					}
					cloneGroupInfo["running_node"] = utils.RemoveDupl(cloneGroupRunNodes)
					cloneGroupInfo["status_message"] = ""
					rscInfo[subGroupId] = cloneGroupInfo
				}
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = utils.RemoveDupl(cloneRunNodes)
				id := rscClone[0].SelectAttr("id").Value
				rscInfo[id] = cloneInfo
			}
			// clone is common resource
			if subRscs := rscClone[0].SelectElements("resource"); len(subRscs) != 0 {
				index := 0
				cloneRunNodes := []string{}
				cloneInfo := map[string]interface{}{}
				for _, subRsc := range subRscs {
					info := map[string]interface{}{}
					info["status"] = GetResourceStatus(subRsc)
					info["status_message"] = ""
					nodename := subRsc.FindElement("node").SelectAttr("name").Value
					id := string(subRsc.SelectAttr("id").Value) + ":" + string(index)
					index++
					info["running_node"] = []string{nodename}
					rscInfo[id] = info
					cloneRunNodes = append(cloneRunNodes, nodename)
				}
				cloneInfo["status"] = GetResourceStatus(rscClone[0])
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = cloneRunNodes
				id := rscClone[0].SelectAttr("id").Value
				rscInfo[id] = cloneInfo
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

	out, err := utils.RunCommand("cibadmin --query --scope resources")
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
						primitiveInfo["running_node"] = priRscStatus["running_node"]
					} else {
						primitiveInfo["status"] = ""
						primitiveInfo["running_node"] = []string{}
					}
					// judgment failure message
					if priFail, ok := failInfo[primitiveId]; ok {
						primitiveInfo["status_message"] = priFail["exitreason"] + " on " + priFail["node"]
						primitiveInfo["status"] = "Failed"
						if priRscStatus, ok := rscStatus[primitiveId]; ok {
							primitiveInfo["running_node"] = priRscStatus["running_node"]
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
	cmd := "crm_resource --resource " + rscId + " --query-xml"
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return ""
	}

	xmlStr := strings.Split(string(out), "xml:")[1]
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xmlStr); err != nil {
		return ""
	}
	rscType := doc.FindElement("primitive").SelectAttrValue("type", "")

	return rscType
}

func GetTopResource() []string {
	result := []string{}

	out, err := utils.RunCommand("cibadmin --query --scope resources")
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
		cmd := "pcs resource enable " + rscID
		_, err := utils.RunCommand(cmd)
		return err
	case "stop":
		cmd := "pcs resource disable " + rscID
		_, err := utils.RunCommand(cmd)
		return err
	case "delete":
		var cmd string
		category := GetResourceCategory(rscID)
		if category == "clone" {
			cmd = "pcs resource delete " + rscID[:len(rscID)-6] + " --force"
		} else {
			// not clone
			cmd = "pcs resource delete " + rscID + " --force"
		}
		_, err := utils.RunCommand(cmd)
		return err
	case "cleanup":
		cmd := "crm_resource --resource " + rscID + " --cleanup"
		_, err := utils.RunCommand(cmd)
		return err
	case "migrate":
		d := struct {
			IsForce bool   `json:"is_force"`
			ToNode  string `json:"to_node"`
			Period  string `json:"period"`
		}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return errors.New("invalid json data")
		}

		cmd := "crm_resource --resource " + rscID + " --move -N " + d.ToNode
		out, err := utils.RunCommand(cmd)
		if err != nil {
			if string(out) == "Error performing operation: Situation already as requested" {
				return errors.New("The resource " + rscID + " is running on node " + d.ToNode + " already!")
			}
		}
		return err
	case "unmigrate":
		cmd := "cibadmin --query --scope constraints"
		out, err := utils.RunCommand(cmd)
		if err != nil {
			return err
		}
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(out); err != nil {
			return err
		}
		rscNames := doc.FindElements("./rsc_location")
		for _, item := range rscNames {
			rsc := item.SelectAttrValue("rsc", "")
			if rscID == rsc {
				locationID := item.SelectAttrValue("id", "")
				cmd2 := "pcs constraint location delete " + locationID
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
			cmd := "pcs constraint location delete " + item
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}

		d := map[string]interface{}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return err
		}
		for _, mapItem := range d["node_level"].([]map[string]string) {
			var score int
			if mapItem["level"] == "Master Node" {
				score = 20000
			} else if mapItem["level"] == "Slave 1" {
				score = 16000
			}
			node := mapItem["node"]
			cmd := "pcs constraint location " + rscID + " prefers " + node + "=" + strconv.Itoa(score)
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
		if err := json.Unmarshal(data, d); err != nil {
			return err
		}
		for _, item := range d.SameNode {
			cmd := "pcs constraint colocation add " + rscID + " with " + item + " INFINITY"
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		for _, item := range d.DiffNode {
			cmd := "pcs constraint colocation add " + rscID + " with " + item + " -INFINITY"
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
	case "order":
		if findOrder(rscID) {
			cmd := "pcs constraint order delete " + rscID
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		d := struct {
			BeforeRscs []string `json:"before_rscs"`
			AfterRscs  []string `json:"after_rscs"`
		}{}
		if err := json.Unmarshal(data, d); err != nil {
			return err
		}
		for _, item := range d.BeforeRscs {
			cmd := "pcs constraint order start " + item + " then " + rscID
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
		for _, item := range d.AfterRscs {
			cmd := "pcs constraint order start " + rscID + " then " + item
			if _, err := utils.RunCommand(cmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func getResourceConstraintIDs(rscID, action string) []string {
	ids := []string{}
	out, err := utils.RunCommand("cibadmin --query --scope constraints")
	if err != nil {
		return ids
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return ids
	}

	if action == "colocation" {
		et := doc.SelectElements("./rsc_colocation")
		for _, item := range et {
			rsc := item.SelectAttrValue("rsc", "")
			rscWith := item.SelectAttrValue("with-rsc", "")
			if rsc == rscID {
				ids = append(ids, rscWith)
			}
			if rscWith == rscID {
				ids = append(ids, rsc)
			}
			return ids
		}
	} else if action == "location" {
		et := doc.SelectElements("./rsc_location")
		for _, item := range et {
			rsc := item.SelectAttrValue("rsc", "")
			if rsc == rscID {
				if item.SelectAttr("score") != nil && item.SelectAttrValue("score", "") == "-INFINITY" {
					continue
				}
				ids = append(ids, item.SelectAttrValue("id", ""))
			}
			return ids
		}
	}
	return ids
}

func DeleteColocationByIdAndAction(rscID string, targetIds []string) error {
	for _, item := range targetIds {
		cmd := "pcs constraint colocation delete " + rscID + " " + item
		if _, err := utils.RunCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func findOrder(rscID string) bool {
	out, err := utils.RunCommand("cibadmin --query --scope constraints")
	if err != nil {
		return false
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return false
	}
	et := doc.SelectElements("./rsc_order")
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
	var result map[string]interface{}

	cmd := "crm_resource --resource " + rscID + " --query-xml"
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
	result, err = GetResourceInfoID(ct, xml)
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
	var data map[string]interface{}

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
	// var prop map[string]interface{}
	e := doc.FindElement("meta_attributes")
	prop, _ := getResourceInfoFromXml("meta", e)
	if len(prop.(map[string]string)) > 0 {
		data["meta_attributes"] = prop
	}

	//For instance_attributes
	e = doc.FindElement("instance_attributes")
	prop, _ = getResourceInfoFromXml("inst", e)
	if len(prop.(map[string]string)) > 0 {
		data["instance_attributes"] = prop
	}

	//For actions
	e = doc.FindElement("operations")
	prop, _ = getResourceInfoFromXml("operations", e)
	if len(prop.([]map[string]string)) > 0 {
		data["action"] = prop
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
			for _, v := range item.Attr {
				i := map[string]string{}
				i[v.Key] = v.Value
				result = append(result, i)
			}
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
		op := Op{}
		op.ID = v.SelectAttrValue("id", "")
		op.Interval = v.SelectAttrValue("interval", "")
		op.Name = v.SelectAttrValue("name", "")
		op.Timeout = v.SelectAttrValue("timeout", "")
		result.Operations = append(result.Operations, op)
	}

	return result
}

func getGroupRscs(groupId string) ([]string, error) {
	cmd := "crm_resource --resource " + groupId + " --query-xml"
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
	et := doc.FindElements("./primitive")
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
	if isIn == false {
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
