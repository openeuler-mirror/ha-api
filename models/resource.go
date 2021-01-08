package models

import (
	"strings"

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

func GetResourceCategory() {

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

func DeleteColocationByIdAndAction() {

}

func CreateResource(data []byte) map[string]interface{} {

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
