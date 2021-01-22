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
	var result map[string]interface{}

	clusterStatus := GetClusterStatus()
	if clusterStatus != 0 {
		result["action"] = true
		result["data"] = []string{}
		return result
	}

	constraints := GetAllConstraints() /////////////////
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
		subRscs := GetSubResources(rscId) ///////////////////
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
	cmd := "cibadmin --query --scope resources|grep 'id=\"" + rscID + "\""
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
	failures := doc.SelectElements("crm_mon/failures/failure")
	if len(failures) == 0 {
		return failInfo
	} else {
		// TODO
		// for _, failure := range failures {
		// 	infoFail := map[string]string{}
		// 	rscIdf := strings.Split(failure.SelectAttr("op_key").Value, "_stop_")[0]
		// 	rscIdm := strings.Split(rscIdf, "_start_")[0]
		// 	rscId := strings.Split(rscIdm, "_start_")[0]
		// 	node := failure.SelectAttr("node")
		// 	exitreason := failure.SelectAttr("exitreason")

		// }
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
	return nil
}

func GetAllConstraints() map[string]interface{} {
	rscStatus := GetAllResourceStatus() ///////////////
	var data map[string](map[string]interface{})

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
		var result map[string]interface{}
		result["action"] = false
		result["data"] = data
		return result
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		var result map[string]interface{}
		result["action"] = false
		result["data"] = data
		return result
	}
	constraints := doc.FindElement("constraints")

	//location
	for _, location := range constraints.FindElements("rsc_location") {
		if strings.HasPrefix(location.SelectAttr("id").Value, "cli-prefer-") {
			continue
		}
		node := location.SelectAttr("node").Value
		rscId := location.SelectAttr("rsc").Value
		score := location.SelectAttr("score").Value
		locationSingle := make(map[string]string)
		locationSingle["node"] = node
		locationSingle["level"] = utils.ScoreToLevel(score) //////////
		locationArr := []map[string]string{}
		for key := range data {
			if rscId == key {
				locationArr = append(locationArr, locationSingle)
			}
		}
		data[rscId]["location"] = locationArr
	}

	//order
	for _, order := range constraints.FindElements("rsc_order") {
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

	//colocation
	for _, colocation := range constraints.FindElements("rsc_colocation") {
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

	failureInfo := GetResourceFailedMessage() ///////////////

	constraintMaps := []map[string]interface{}{}
	for rscId := range data {
		var constraint map[string]interface{}
		constraint["id"] = rscId
		rscIDFirst := strings.Split(rscId, ":")[0]
		if _, ok := rscStatus[rscId]; !ok {
			if _, ok := failureInfo[rscIDFirst]; !ok {
				constraint["status"] = "Stopped"
				constraint["status_message"] = ""
				constraint["running_node"] = []string{}
			}
		} else if strings.HasSuffix(rscId, "-clone") {
			constraint["status"] = "Failed"
			////////////////////
			constraint["status_message"] = failureInfo[rscIDFirst]["exitreason"] + " on " + failureInfo[rscIDFirst]["node"]
			constraint["running_node"] = []string{}
		} else {
			rscInfo := rscStatus[rscId]
			constraint["status"] = rscInfo["status"]
			constraint["status_message"] = ""
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

	var result map[string]interface{}
	result["action"] = true
	result["data"] = constraints
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

func GetAllResourceStatus() map[string]map[string]string {
	// rscInfo:=
	out, err := utils.RunCommand("crm_mon -1 --as-xml")
	if err != nil {
		return map[string]map[string]string{}
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return map[string]map[string]string{}
	}

	if len(doc.SelectElements("crm_mon/resources")) == 0 {
		return map[string]map[string]string{}
	}
	// allRscRun:=[]string{} // doesn't used
	rscClone := doc.SelectElements("crm_mon/resources/clone")
	rscGroup := doc.SelectElements("crm_mon/resources/group")
	rscResource := doc.SelectElements("crm_mon/resources/resource")

	if len(rscClone) != 0 {
		if len(rscClone) != 1 {
			// TODO
			// for _, rsc := range rscClone {
			// 	subRscs := rsc.FindElement("resource")
			// 	index := 0
			//
			// 	for _, subRsc := range subRscs {
			//
			// 	}
			// }
		}
	}
	if len(rscGroup) != 0 {

	}
	if len(rscResource) != 0 {

	}

	return map[string]map[string]string{}
}

func GetSubResources(rscId string) map[string]interface{} {
	// TODO
	// rscStatus := GetAllResourceStatus()
	// failInfo := GetResourceFailedMessage()

	return nil
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
	rscType := doc.SelectElement("primitive").SelectAttrValue("type", "")

	return rscType
}

func GetTopResource() []string {
	result := []string{}

	out, err := utils.RunCommand("cibadmin --query --scope resources")
	if err != nil {
		return result
	}

	// TODO: check logic
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return result
	}

	elements := doc.SelectElements("resources/clone")
	for _, element := range elements {
		result = append(result, element.SelectAttr("id").Value)
	}
	elements = doc.SelectElements("resources/group")
	for _, element := range elements {
		result = append(result, element.SelectAttr("id").Value)
	}
	elements = doc.SelectElements("resources/primitive")
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
	return nil
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
	// TODO:
	return false
}

func GetResourceInfoByrscID(rscID string) interface{} {
	cmd := "crm_resource --resource " + rscID + " --query-xml"
	out, err := utils.RunCommand(cmd)
	return err

	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		return ""
	}

	ct := doc.Root().Tag
	d, err := GetResourceInfoID(ct, xml)
	data := d.(map[string]string)
	data["id"] = string(rscID)
	data["category"] = string(ct)
	var result map[string]interface{}
	if len(data) != 0 {
		result["data"] = data
		result["action"] = true
	} else {
		result["error"] = data
		result["action"] = false
	}
	if _, ok := result["data"]; ok {
		dataRes := result["data"].(map[string]string)
		if _, ok := dataRes["provider"]; ok {
			provider := dataRes["provider"]
			if len(provider) == 0 {
				delete(result, "data")
				delete(result, "provider")
			}
		}
	}
	return result
}

func GetResourceInfoID(ct, xmlData string) (interface{}, error) {
	doc := etree.NewDocument()
	doc.ReadFromString(xmlData)
	var data map[string]interface{}
	switch ct {
	case "primitive":
		// TODO:
		data, _ = getResourceInfoFromXml("primitive", doc.Root())
	case "group":
		data["rscs"], _ = getResourceInfoFromXml("group", doc.Root())
	case "clone":
		data["rsc_id"], _ = getResourceInfoFromXml("clone", doc.Root())
	}

	// For meta_attributes
	var prop map[string]interface{}
	e := doc.FindElement("meta_attributes")
	prop, _ = getResourceInfoFromXml("meta", e)
	data["meta_attributes"] = prop

	//For instance_attributes
	e = doc.FindElement("instance_attributes")
	prop, _ = getResourceInfoFromXml("inst", e)
	data["instance_attributes"] = prop

	//For actions
	e = doc.FindElement("operations")
	prop, _ = getResourceInfoFromXml("operations", e)
	data["action"] = prop

	return data, nil
}

func getResourceInfoFromXml(cl string, et *etree.Element) (map[string]interface{}, error) {
	var prop map[string]interface{}
	if cl == "group" {
		els := et.FindElements("primitive")
		for _, e := range els {
			for _, attr := range e.Attr {
				prop[attr.Key] = attr.Value
			}
		}
	} else if cl == "primitive" {
		//TODO
	} else if cl == "clone" {
		// et.FindElements("group"){
		// 	return
	} else if cl == "meta" || cl == "inst" {
		// prop := map[string]string{}
		// for _,item:=range et.FindElements("./nvpair")
		//TODO

	} else if cl == "operations" {
		// var prop = []string{}
		// result := map[string]string{}
		// op := et.FindElements("./op")
		// for _, item := range op {

		// for k, v := range item.items() {
		// 	result[string(k)] = string(v)
		// 	prop = append(prop, result)
		// }
		// }
		//TODO
	}

	return prop, nil
}
