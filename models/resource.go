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

// GerResourceInfo
func GerResourceInfo() map[string]interface{} {
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

func GetResourceConstraints() {

}

func GetResourceFailedMessage() {}

func GetResourceMetaAttributes() {}

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
}

func findOrder(rscID string) bool {
	// TODO:
	return false
}
