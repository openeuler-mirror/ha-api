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
		result["action"] = []string{}
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

func GetResourceFailedMessage() {}

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
	// TODO:
	return nil
}

func GetAllConstraints() map[string]interface{} {
	return nil
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

	return result
}

func GetSubResources(rscId string) map[string]interface{} {
	// TODO:
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
