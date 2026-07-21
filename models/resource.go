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
	"log/slog"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

// safeResourceName 校验资源名称，仅允许字母、数字、下划线、连字符、冒号和点号
var safeResourceName = regexp.MustCompile(`^[a-zA-Z0-9_\-:.]+$`)

// safePeriod 校验迁移生命周期参数，仅允许字母和数字
var safePeriod = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// safeIdentifier 校验 pcs 标识符（class/provider/type/属性键/操作名），仅允许字母、数字、下划线和连字符
var safeIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

func validateResourceID(id string) error {
	if id == "" || !safeResourceName.MatchString(id) {
		return fmt.Errorf("invalid resource id: %q", id)
	}
	return nil
}

func validateIdentifier(value, fieldName string) error {
	if value == "" || !safeIdentifier.MatchString(value) {
		return fmt.Errorf("invalid %s: %q", fieldName, value)
	}
	return nil
}

// splitXMLOutput 安全地从 crm_resource 输出中提取 XML 部分
func splitXMLOutput(out []byte) (string, error) {
	parts := strings.SplitN(string(out), ":\n", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected command output format")
	}
	return parts[1], nil
}

func GetAllResourceStatusForNew() []map[string]interface{} {
	slog.Debug("Get all resource status")
	rscInfo := []map[string]interface{}{}
	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	if err != nil {
		return []map[string]interface{}{}
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return []map[string]interface{}{}
	}

	if len(doc.FindElements("/crm_mon/resources")) == 0 {
		return []map[string]interface{}{}
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
				subRes := []map[string]interface{}{}
				cloneInfo["status"] = "Running"
				for _, subRsc := range subRscs {
					info := map[string]interface{}{}
					info["status"] = GetResourceStatus(subRsc)
					roleVal := subRsc.SelectAttrValue("role", "")
					if roleVal == "Slave" || roleVal == "Master" {
						isMs = true
					}
					info["status_message"] = ""
					nodename := ""
					if node := subRsc.FindElement("node"); node != nil {
						nodename = node.SelectAttrValue("name", "")
					}
					id := subRsc.SelectAttrValue("id", "") + ":" + strconv.Itoa(index)
					index++
					info["running_node"] = []string{nodename}
					info["id"] = id
					info["type"] = "primitive"
					info["svc"] = GetResourceSvcFromXml(subRsc)

					cloneRunNodes = append(cloneRunNodes, nodename)

					if info["status"] != "Running" && info["status"] != "Running but failed" {
						cloneInfo["status"] = "Not Running"
					} else if info["status"] == "Running but failed" && cloneInfo["status"] != "Not Running" {
						cloneInfo["status"] = "Running but failed"
					}

					subRes = append(subRes, info)
				}
				cloneInfo["subrscs"] = subRes
				if isMs {
					cloneInfo["isMs"] = true
				} else {
					cloneInfo["isMs"] = false
				}
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = cloneRunNodes
				cloneInfo["type"] = "clone"
				cloneId := rsc.SelectAttrValue("id", "")
				cloneInfo["id"] = cloneId
				isManaged, _ := strconv.ParseBool(rsc.SelectAttrValue("managed", ""))
				if !isManaged {
					cloneInfo["status"] = "Unmanaged"
				}
				rscInfo = append(rscInfo, cloneInfo)

			}
			// subResources is gourp resources
			if subRscs := rsc.SelectElements("group"); len(subRscs) != 0 {
				cloneRunNodes := []string{}
				cloneInfo := map[string]interface{}{}
				subRes := []map[string]interface{}{}
				cloneInfo["status"] = "Running"
				for _, subRsc := range subRscs {
					subRscId := subRsc.SelectAttrValue("id", "")
					groupInfo := map[string]interface{}{}
					groupInfo["status"] = "Running"
					groupRunNodes := []string{}
					groupSubres := []map[string]interface{}{}
					if innerRscs := subRsc.SelectElements("resource"); len(innerRscs) != 0 {
						for _, innerRsc := range innerRscs {
							innerRscId := innerRsc.SelectAttrValue("id", "")
							info := map[string]interface{}{}
							info["status"] = GetResourceStatus(innerRsc)

							info["status_message"] = ""
							if node := innerRsc.FindElement("node"); node != nil {
								nodename := node.SelectAttrValue("name", "")
								info["running_node"] = []string{nodename}
								cloneRunNodes = append(cloneRunNodes, nodename)
								groupRunNodes = append(groupRunNodes, nodename)
							}
							parts := strings.SplitN(subRscId, ":", 2)
							fatherId := ""
							if len(parts) > 1 {
								fatherId = parts[1]
							}
							id := innerRscId + ":" + fatherId

							info["id"] = id
							info["type"] = "primitive"
							info["svc"] = GetResourceSvcFromXml(innerRsc)

							if info["status"] != "Running" && info["status"] != "Running but failed" {
								groupInfo["status"] = "Not Running"
							} else if info["status"] == "Running but failed" && groupInfo["status"] != "Not Running" {
								groupInfo["status"] = "Running but failed"
							}

							if groupInfo["status"] != "Running" && groupInfo["status"] != "Running but failed" {
								cloneInfo["status"] = "Not Running"
							} else if groupInfo["status"] == "Running but failed" && cloneInfo["status"] != "Not Running" {
								cloneInfo["status"] = "Running but failed"
							}

							groupSubres = append(groupSubres, info)
						}
					}
					isManaged, _ := strconv.ParseBool(subRsc.SelectAttrValue("managed", ""))
					if !isManaged {
						groupInfo["status"] = "Unmanaged"
					}
					groupInfo["running_node"] = utils.RemoveDupl(groupRunNodes)
					groupInfo["status_message"] = ""
					groupInfo["type"] = "group"
					groupId := subRscId
					groupInfo["id"] = groupId
					groupInfo["subrscs"] = groupSubres
					subRes = append(subRes, groupInfo)

				}
				cloneInfo["status_message"] = ""
				cloneInfo["running_node"] = utils.RemoveDupl(cloneRunNodes)
				cloneInfo["isMs"] = false
				cloneInfo["type"] = "clone"
				cloneInfo["subrscs"] = subRes
				cloneId := rsc.SelectAttrValue("id", "")
				//rscInfo[cloneId] = cloneInfo
				cloneInfo["id"] = cloneId
				isManaged, _ := strconv.ParseBool(rsc.SelectAttrValue("managed", ""))
				if !isManaged {
					cloneInfo["status"] = "Unmanaged"
				}
				rscInfo = append(rscInfo, cloneInfo)
			}
		}
	}

	if len(rscGroup) != 0 {
		for _, rsc := range rscGroup {
			subRscs := rsc.SelectElements("resource")
			groupRunNodes := []string{}
			groupInfo := map[string]interface{}{}
			groupInfo["status"] = "Running"
			groupSubres := []map[string]interface{}{}
			// several resources in each group

			for _, subRsc := range subRscs {
				info := map[string]interface{}{}
				info["status"] = GetResourceStatus(subRsc)

				info["status_message"] = ""
				nodename := ""
				if node := subRsc.FindElement("node"); node != nil {
					nodename = node.SelectAttrValue("name", "")
					info["running_node"] = []string{nodename}
					groupRunNodes = append(groupRunNodes, nodename)
					//info["status"] = "Running"
				}
				if info["status"] != "Running" && info["status"] != "Running but failed" {
					groupInfo["status"] = "Not Running"
				} else if info["status"] == "Running but failed" && groupInfo["status"] != "Not Running" {
					groupInfo["status"] = "Running but failed"
				}
				id := subRsc.SelectAttrValue("id", "")
				info["id"] = id
				info["type"] = "primitive"
				info["svc"] = GetResourceSvcFromXml(subRsc)

				groupSubres = append(groupSubres, info)

			}
			groupInfo["status_message"] = ""
			groupInfo["running_node"] = groupRunNodes
			groupInfo["subrscs"] = groupSubres

			isManaged, _ := strconv.ParseBool(rsc.SelectAttrValue("managed", ""))
			if !isManaged {
				groupInfo["status"] = "Unmanaged"
			}
			id := rsc.SelectAttrValue("id", "")
			groupInfo["id"] = id
			groupInfo["type"] = "group"
			location, err := GetResourceConstraints(id, "location")
			if err == nil {
				groupInfo["location"] = location["node_level"]
			}

			rscInfo = append(rscInfo, groupInfo)
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
					runningNode = append(runningNode, node.SelectAttrValue("name", ""))
				}
			}
			resourceInfo["running_node"] = runningNode
			resourceInfo["status_message"] = ""

			resourceInfo["type"] = "primitive"
			resourceInfo["svc"] = GetResourceSvcFromXml(rsc)
			id := rsc.SelectAttrValue("id", "")
			resourceInfo["id"] = id

			location, err := GetResourceConstraints(id, "location")
			if err == nil {
				resourceInfo["location"] = location["node_level"]
			}

			rscInfo = append(rscInfo, resourceInfo)
		}
	}

	return rscInfo
}

func GetResourceSvcFromXml(rscInfo *etree.Element) string {
	attr := rscInfo.SelectAttr("resource_agent")
	if attr == nil {
		return ""
	}

	parts := strings.Split(attr.Value, ":")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		return lastPart
	}
	return ""
}

// GetResourceInfo
var GetResourceInfo = func() map[string]interface{} {
	slog.Debug("Get Resource Info")
	result := make(map[string]interface{})
	clusterStatus := GetClusterStatus()
	if clusterStatus != 0 {
		result["action"] = false
		result["error"] = gettext.Gettext("The current node cluster status is incorrect")
		result["data"] = []string{}
		return result
	}

	// 在获取集群资源时为全局变量failInfo、clusterPro赋值，避免操作过程多次查询对响应时长造成影响
	setFailInfo(GetResourceFailedMessage())
	setGlobalClusterProperties(GetClusterPropertiesInfo())

	resStatus := GetAllResourceStatusForNew()
	result["action"] = true
	result["data"] = resStatus
	return result
}

func GetResourceCategory(rscID string) string {
	if err := validateResourceID(rscID); err != nil {
		slog.Error("GetResourceCategory: invalid resource id", "id", rscID, "err", err)
		return ""
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscID, "--query-xml")
	if err != nil {
		slog.Error("GetResourceCategory: command failed", "id", rscID, "err", err)
		return ""
	}
	xml, err := splitXMLOutput(out)
	if err != nil {
		slog.Error("GetResourceCategory: failed to split XML output", "id", rscID, "err", err)
		return ""
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		slog.Error("GetResourceCategory: failed to parse XML", "id", rscID, "err", err)
		return ""
	}
	return doc.Root().Tag
}

// 基于cib查询资源类型
func GetResourceType(rscID string) string {
	if err := validateResourceID(rscID); err != nil {
		slog.Error("GetResourceType: invalid resource id", "id", rscID, "err", err)
		return ""
	}
	out, err := utils.RunCommandWithArgs("cibadmin", "--query", "--scope", "resources")
	if err != nil {
		slog.Error("GetResourceType: command failed", "id", rscID, "err", err)
		return ""
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(out); err != nil {
		slog.Error("GetResourceType: failed to parse XML", "id", rscID, "err", err)
		return ""
	}
	for _, tag := range []string{"primitive", "group", "clone"} {
		for _, elem := range doc.FindElements("//" + tag) {
			if elem.SelectAttrValue("id", "") == rscID {
				return tag
			}
		}
	}
	return ""
}

type OrderList struct {
	RscName  string `json:"rsc_name"`
	Location string `json:"location"`
	Baction  string `json:"before_action"`
	Aaction  string `json:"after_action"`
}

// TODO needs to integrate to func GetResourceByConstraintAndId
// or func GetAllConstraints??
var GetResourceConstraints = func(rscID, relation string) (map[string]interface{}, error) {
	slog.Debug("GetResourceConstraints", "id", rscID, "relation", relation)
	if err := validateResourceID(rscID); err != nil {
		return nil, err
	}
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
			rsc := resourceLocation.SelectAttrValue("rsc", "")
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
		et := root.FindElements("./rsc_colocation")
		sameNodes, diffNodes := getRscColocation(et, rscID)
		retData["same_node"] = sameNodes
		retData["rsc_id"] = rscID
		retData["diff_node"] = diffNodes
	case "order":
		var retData1 []OrderList
		for _, order := range root.FindElements("rsc_order") {
			first := order.SelectAttrValue("first", "")
			then := order.SelectAttrValue("then", "")
			baction := order.SelectAttrValue("first-action", "")
			aaction := order.SelectAttrValue("then-action", "")
			var tmp OrderList
			if first == rscID {
				tmp.RscName = then
				tmp.Location = "after"
			} else if then == rscID {
				tmp.RscName = first
				tmp.Location = "before"
			} else {
				continue
			}
			tmp.Baction = baction
			tmp.Aaction = aaction
			retData1 = append(retData1, tmp)
		}
		retData["order"] = retData1
	}
	return retData, nil
}

func getRscColocation(et []*etree.Element, rscID string) ([]string, []string) {
	sameNode := make([]string, 0)
	diffNode := make([]string, 0)

	for _, item := range et {
		role := ""
		rsc := item.SelectAttrValue("rsc", "")
		rscWith := item.SelectAttrValue("with-rsc", "")

		// 处理with-rsc-role属性
		if hasAttribute(item, "with-rsc-role") {
			if roleVal := item.SelectAttrValue("with-rsc-role", ""); roleVal != "" {
				role = roleVal
				rscWith += "/" + role
			}
		}

		// 主资源匹配逻辑
		if rsc == rscID {
			switch score := item.SelectAttrValue("score", ""); score {
			case "INFINITY":
				sameNode = append(sameNode, rscWith)
			case "-INFINITY":
				diffNode = append(diffNode, rscWith)
			}
		}

		// 反向关联检查
		if rscWith == rscID {
			rscEntry := rsc
			if hasAttribute(item, "rsc-role") {
				if roleVal := item.SelectAttrValue("rsc-role", ""); roleVal != "" {
					rscEntry += "/" + roleVal
				}
			}
			switch score := item.SelectAttrValue("score", ""); score {
			case "INFINITY":
				sameNode = append(sameNode, rscEntry)
			case "-INFINITY":
				diffNode = append(diffNode, rscEntry)
			}
		}
	}
	return sameNode, diffNode
}

// 自定义属性检查函数
func hasAttribute(e *etree.Element, attrName string) bool {
	for _, attr := range e.Attr {
		if attr.Key == attrName {
			return true
		}
	}
	return false
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
	slog.Debug("GetResourceFailedMessage")
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
	}

	for _, failure := range failures {
		infoFail := map[string]string{}
		// 提取rscID
		rscId := extractRscID(failure.SelectAttrValue("op_key", ""))

		// 处理失败项信息
		node := failure.SelectAttrValue("node", "")
		exitreason := failure.SelectAttrValue("exitreason", "")
		infoFail["node"] = node
		infoFail["exitreason"] = exitreason
		failInfo[rscId] = infoFail
	}

	return failInfo
}

func GetResourceFailedList() []string {
	out, err := utils.RunCommand(utils.CmdClusterStatusAsXML)
	var failList []string
	if err != nil {
		return failList
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return failList
	}
	failures := doc.FindElements("/crm_mon/failures/failure")
	if len(failures) == 0 {
		return failList
	} else {

		seen := make(map[string]struct{}, len(failures))
		failList = make([]string, 0, len(failures))

		for _, failure := range failures {
			rscId := extractRscID(failure.SelectAttrValue("op_key", ""))
			if _, exists := seen[rscId]; !exists {
				seen[rscId] = struct{}{}
				failList = append(failList, rscId)
			}
		}
		return failList
	}
}

var rscIDRegex = regexp.MustCompile(`(_stop_|_start_|_monitor_|_demote_|_promote_)`)

func extractRscID(opKey string) string {
	if opKey == "" {
		return ""
	}
	loc := rscIDRegex.FindStringIndex(opKey)
	if loc != nil {
		start := loc[0]
		return opKey[0:start]
	}
	return opKey
}

func GetResourceMetaAttributes(category string) map[string]interface{} {
	slog.Debug("GetResourceMetaAttributes", "category", category)
	retjson := make(map[string](map[string]interface{}))

	retjson["target-role"] = make(map[string]interface{})
	retjson["target-role"]["content"] = make(map[string]interface{})
	retjson["target-role"]["name"] = "target-role"
	retjson["target-role"]["content"].(map[string]interface{})["values"] = []string{"Stopped", "Started"}
	retjson["target-role"]["content"].(map[string]interface{})["default"] = "Stopped"
	retjson["target-role"]["content"].(map[string]interface{})["type"] = "enum"
	retjson["target-role"]["content"].(map[string]interface{})["desc"] = "What state should the cluster attempt to keep this resource in?"

	retjson["priority"] = make(map[string]interface{})
	retjson["priority"]["content"] = make(map[string]interface{})
	retjson["priority"]["name"] = "priority"
	retjson["priority"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["priority"]["content"].(map[string]interface{})["desc"] = "If not all resources can be active, the cluster will stop lower priority resources in order to keep higher priority ones active."

	retjson["is-managed"] = make(map[string]interface{})
	retjson["is-managed"]["content"] = make(map[string]interface{})
	retjson["is-managed"]["name"] = "is-managed"
	retjson["is-managed"]["content"].(map[string]interface{})["type"] = "boolean"
	retjson["is-managed"]["content"].(map[string]interface{})["desc"] = "Is the cluster allowed to start and stop the resource?"

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
	retjson["resource-stickiness"]["content"].(map[string]interface{})["desc"] = "How much does the resource prefer to stay where it is? Defaults to the value of \"default-resource-stickiness\""

	retjson["migration-threshold"] = make(map[string]interface{})
	retjson["migration-threshold"]["content"] = make(map[string]interface{})
	retjson["migration-threshold"]["name"] = "migration-threshold"
	retjson["migration-threshold"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["migration-threshold"]["content"].(map[string]interface{})["desc"] = "How many failures should occur for this resource on a node before making the node ineligible to host this resource. Default: \"none\""

	retjson["multiple-active"] = make(map[string]interface{})
	retjson["multiple-active"]["content"] = make(map[string]interface{})
	retjson["multiple-active"]["name"] = "multiple-active"
	retjson["multiple-active"]["content"].(map[string]interface{})["values"] = []string{"stop_start", "stop_only", "block"}
	retjson["multiple-active"]["content"].(map[string]interface{})["type"] = "enum"
	retjson["multiple-active"]["content"].(map[string]interface{})["desc"] = "What should the cluster do if it ever finds the resource active on more than one node."

	retjson["failure-timeout"] = make(map[string]interface{})
	retjson["failure-timeout"]["content"] = make(map[string]interface{})
	retjson["failure-timeout"]["name"] = "failure-timeout"
	retjson["failure-timeout"]["content"].(map[string]interface{})["type"] = "integer"
	retjson["failure-timeout"]["content"].(map[string]interface{})["desc"] = "How many seconds to wait before acting as if the failure had not occurred (and potentially allowing the resource back to the node on which it failed. Default: \"never\""

	retjson["allow-migrate"] = make(map[string]interface{})
	retjson["allow-migrate"]["content"] = make(map[string]interface{})
	retjson["allow-migrate"]["name"] = "allow-migrate"
	retjson["allow-migrate"]["content"].(map[string]interface{})["type"] = "boolean"
	retjson["allow-migrate"]["content"].(map[string]interface{})["desc"] = "Allow resource migration for resources which support migrate_to/migrate_from actions."

	retjson["allow-unhealthy-nodes"] = make(map[string]interface{})
	retjson["allow-unhealthy-nodes"]["content"] = make(map[string]interface{})
	retjson["allow-unhealthy-nodes"]["name"] = "allow-unhealthy-nodes"
	retjson["allow-unhealthy-nodes"]["content"].(map[string]interface{})["type"] = "boolean"
	retjson["allow-unhealthy-nodes"]["content"].(map[string]interface{})["desc"] = "Whether the resource should be able to run on a node even if the node's health score would otherwise prevent it."

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
		retjson["interleave"]["content"].(map[string]interface{})["desc"] = "Changes the behavior of ordering constraints (between clones/masters) so that instances can start/stop as soon as their peer instance has (rather than waiting for every instance of the other clone has)."

		retjson["clone-max"] = make(map[string]interface{})
		retjson["clone-max"]["content"] = make(map[string]interface{})
		retjson["clone-max"]["name"] = "clone-max"
		retjson["clone-max"]["content"].(map[string]interface{})["type"] = "integer"
		retjson["clone-max"]["content"].(map[string]interface{})["desc"] = "How many copies of the resource to start. Defaults to the number of nodes in the cluster."

		retjson["promoted-max"] = make(map[string]interface{})
		retjson["promoted-max"]["content"] = make(map[string]interface{})
		retjson["promoted-max"]["name"] = "promoted-max"
		retjson["promoted-max"]["content"].(map[string]interface{})["type"] = "integer"
		retjson["promoted-max"]["content"].(map[string]interface{})["desc"] = "If promotable is true, the number of instances that can be promoted at one time across the entire cluster"

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

// 命令执行接口
type CommandRunner interface {
	RunCommand(cmd string) ([]byte, error)
	RunCommandWithArgs(binary string, args ...string) ([]byte, error)
}

// 默认命令执行器
type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) RunCommand(cmd string) ([]byte, error) {
	return utils.RunCommand(cmd)
}

func (r *DefaultCommandRunner) RunCommandWithArgs(binary string, args ...string) ([]byte, error) {
	return utils.RunCommandWithArgs(binary, args...)
}

// 全局命令执行器实例
var (
	cmdRunner   CommandRunner = &DefaultCommandRunner{}
	runnerMutex sync.Mutex
)

// 设置命令执行器（用于测试）
func SetCommandRunner(runner CommandRunner) {
	runnerMutex.Lock()
	defer runnerMutex.Unlock()
	cmdRunner = runner
}

// 获取当前命令执行器
func getCommandRunner() CommandRunner {
	runnerMutex.Lock()
	defer runnerMutex.Unlock()
	return cmdRunner
}

// ResourceRequest 定义创建资源的请求结构
type ResourceRequest struct {
	Category           string                 `json:"category"`
	ID                 string                 `json:"id"`
	MetaAttributes     map[string]interface{} `json:"meta_attributes,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Class              string                 `json:"class,omitempty"`
	Provider           string                 `json:"provider,omitempty"`
	InstanceAttributes map[string]interface{} `json:"instance_attributes,omitempty"`
	Actions            []Action               `json:"actions,omitempty"`
	Rscs               []string               `json:"rscs,omitempty"`
	RscID              string                 `json:"rsc_id,omitempty"`
	SelfFlag           bool                   `json:"selfFlag,omitempty"`
}

// Action 定义资源的操作属性
type Action struct {
	Name          string `json:"name"`
	Interval      string `json:"interval,omitempty"`
	StartDelay    string `json:"start-delay,omitempty"`
	Timeout       string `json:"timeout,omitempty"`
	Role          string `json:"role,omitempty"`
	Requires      string `json:"requires,omitempty"`
	OnFail        string `json:"on-fail,omitempty"`
	OCFCheckLevel string `json:"OCF_CHECK_LEVEL,omitempty"`
	Dep           string `json:"depth,omitempty"`
}

// CreateResource 创建资源
func CreateResource(dataBytes []byte) utils.GeneralResponse {
	if len(dataBytes) == 0 {
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("No input data"),
		}
	}
	var data ResourceRequest
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Cannot convert data to json"),
		}
	}
	if data.Category == "" || data.ID == "" {
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("category and id are required"),
		}
	}
	var res utils.GeneralResponse
	switch data.Category {
	case "primitive":
		res = createPrimitiveResource(data)
	case "group":
		res = createGroupResource(data)
	case "clone":
		res = createCloneResource(data)
	default:
		slog.Error("unsupported resource category", "category", data.Category)
		res = utils.GeneralResponse{
			Action: false,
			Error:  fmt.Sprintf(gettext.Gettext("unsupported resource category: %s"), data.Category),
		}
	}
	if data.SelfFlag {
		var emptySlice []byte
		startRes := ResourceAction(data.ID, "start", emptySlice)
		slog.Info("start resource", "rsc_id", data.ID, "res", startRes)
	}
	return res
}

// 创建primitive资源
func createPrimitiveResource(data ResourceRequest) utils.GeneralResponse {
	// 校验资源ID格式
	if err := validateResourceID(data.ID); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}

	// 校验 class / type / provider 格式
	if err := validateIdentifier(data.Class, "class"); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}
	if err := validateIdentifier(data.Type, "type"); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}
	if data.Provider != "" {
		if err := validateIdentifier(data.Provider, "provider"); err != nil {
			return utils.GeneralResponse{Action: false, Error: err.Error()}
		}
	}

	// 校验属性键名
	for key := range data.InstanceAttributes {
		if err := validateIdentifier(key, "instance attribute key"); err != nil {
			return utils.GeneralResponse{Action: false, Error: err.Error()}
		}
	}
	for key := range data.MetaAttributes {
		if err := validateIdentifier(key, "meta attribute key"); err != nil {
			return utils.GeneralResponse{Action: false, Error: err.Error()}
		}
	}

	// 校验操作名称
	for _, action := range data.Actions {
		if err := validateIdentifier(action.Name, "action name"); err != nil {
			return utils.GeneralResponse{Action: false, Error: err.Error()}
		}
	}

	// 检查资源ID是否已存在
	exists, err := resourceExists(data.ID)
	if err != nil {
		return utils.GeneralResponse{
			Action: false,
			Error:  fmt.Sprintf("failed to check resource existence: %v", err),
		}
	}
	if exists {
		return utils.GeneralResponse{
			Action: false,
			Error:  fmt.Sprintf(gettext.Gettext("%s:The resource name already exists"), data.ID),
		}
	}

	// Health类型的检查
	healthTypes := []string{"HealthCPU", "HealthMEM", "HealthIOWait", "HealthSMART"}
	if utils.Contains(healthTypes, data.Type) {
		typeExists, err := resourceTypeExists(data.Type)
		if err != nil {
			return utils.GeneralResponse{
				Action: false,
				Error:  fmt.Sprintf("failed to check resource type existence: %v", err),
			}
		}
		if typeExists {
			return utils.GeneralResponse{
				Action: false,
				Error:  fmt.Sprintf(gettext.Gettext("%s: The corresponding type of resource already exists"), data.Type),
			}
		}
	}

	// 构建类型规格: class:provider:type 或 class:type (stonith: provider:type 或 type)
	var typeSpec string
	if data.Class == "stonith" {
		if data.Provider != "" {
			typeSpec = data.Provider + ":" + data.Type
		} else {
			typeSpec = data.Type
		}
	} else {
		typeSpec = data.Class + ":"
		if data.Provider != "" {
			typeSpec += data.Provider + ":"
		}
		typeSpec += data.Type
	}

	// 构建命令参数
	var args []string
	if data.Class == "stonith" {
		args = []string{"stonith", "create", data.ID, typeSpec}
	} else {
		args = []string{"resource", "create", data.ID, typeSpec}
	}

	// 添加实例属性
	for key, value := range data.InstanceAttributes {
		if value == nil {
			continue
		}
		args = append(args, key+"="+fmt.Sprintf("%v", value))
	}

	// 添加元属性
	metaAdded := false
	for key, value := range data.MetaAttributes {
		if value == nil {
			continue
		}
		if !metaAdded {
			args = append(args, "meta")
			metaAdded = true
		}
		args = append(args, key+"="+fmt.Sprintf("%v", value))
	}

	// 如果没有设置target-role，则默认为Stopped
	if _, exists := data.MetaAttributes["target-role"]; !exists {
		if !metaAdded {
			args = append(args, "meta")
			metaAdded = true
		}
		args = append(args, "target-role=Stopped")
	}

	// 添加操作
	for _, action := range data.Actions {
		args = append(args, "op", action.Name)
		if action.Interval != "" {
			args = append(args, "interval="+action.Interval)
		}
		if action.StartDelay != "" {
			args = append(args, "start-delay="+action.StartDelay)
		}
		if action.Timeout != "" {
			args = append(args, "timeout="+action.Timeout)
		}
		if action.Role != "" {
			args = append(args, "role="+action.Role)
		}
		if action.Requires != "" {
			args = append(args, "requires="+action.Requires)
		}
		if action.OnFail != "" {
			args = append(args, "on-fail="+action.OnFail)
		}
		if action.OCFCheckLevel != "" {
			args = append(args, "OCF_CHECK_LEVEL="+action.OCFCheckLevel)
		}
	}

	args = append(args, "--no-default-ops", "--force")

	// 执行命令（不经 shell，消除注入风险）
	_, err = getCommandRunner().RunCommandWithArgs("pcs", args...)
	if err != nil {
		slog.Error("Create primitive resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add primitive resource failed"),
		}
	}

	// 如果没有提供monitor操作，则使用默认的monitor操作
	hasMonitor := false
	for _, action := range data.Actions {
		if action.Name == "monitor" {
			hasMonitor = true
			break
		}
	}

	if !hasMonitor {
		updateArgs := []string{"resource", "update", data.ID, "op", "monitor"}

		monitorAction := getDefaultMonitorAction(data.Class, data.Provider, data.Type)
		if monitorAction.Interval != "" {
			updateArgs = append(updateArgs, "interval="+monitorAction.Interval)
		}
		if monitorAction.Timeout != "" {
			updateArgs = append(updateArgs, "timeout="+monitorAction.Timeout)
		}
		if monitorAction.OCFCheckLevel != "" {
			updateArgs = append(updateArgs, "OCF_CHECK_LEVEL="+monitorAction.OCFCheckLevel)
		}

		_, err := getCommandRunner().RunCommandWithArgs("pcs", updateArgs...)
		if err != nil {
			slog.Error("Add monitor action of primitive resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
			return utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Add primitive resource failed"),
			}
		}
	}
	slog.Info("create primitive resource success", "rscId", data.ID)

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Add " + data.Category + " resource success"),
	}
}

// 创建group资源
func createGroupResource(data ResourceRequest) utils.GeneralResponse {
	// 校验资源ID格式
	if err := validateResourceID(data.ID); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}
	for _, rsc := range data.Rscs {
		if err := validateResourceID(rsc); err != nil {
			return utils.GeneralResponse{Action: false, Error: fmt.Sprintf("invalid sub-resource: %v", err)}
		}
	}

	// 检查资源ID是否已存在
	exists, err := resourceExists(data.ID)
	if err != nil {
		return utils.GeneralResponse{
			Action: false,
			Error:  fmt.Sprintf("failed to check resource existence: %v", err),
		}
	}
	if exists {
		return utils.GeneralResponse{
			Action: false,
			Error:  fmt.Sprintf(gettext.Gettext("%s:The resource name already exists"), data.ID),
		}
	}

	// 删除组内每个资源的属性
	for _, rsc := range data.Rscs {
		if err := DeletePriAttrib(rsc); err != nil {
			slog.Error("Delete attrib of prim resource in group failed", "id", data.ID, "category", data.Category, "rsc", rsc, "err", err.Error())
			return utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Add group resource failed"),
			}
		}
	}

	args := []string{"resource", "group", "add", data.ID}
	args = append(args, data.Rscs...)

	_, err = getCommandRunner().RunCommandWithArgs("pcs", args...)
	if err != nil {
		slog.Error("Create group resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add group resource failed"),
		}
	}

	// 更新组资源的属性
	if err := UpdateResourceAttributes(data.ID, data); err != nil {
		slog.Error("Update attributes of group resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add group resource failed"),
		}
	}

	// 检查是否有target-role，如果没有则禁用资源
	if _, exists := data.MetaAttributes["target-role"]; !exists {
		_, err := getCommandRunner().RunCommandWithArgs("pcs", "resource", "disable", data.ID)
		if err != nil {
			slog.Error("Update target-role of group resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
			return utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Add group resource failed"),
			}
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Add " + data.Category + " resource success"),
	}
}

// 创建clone资源
func createCloneResource(data ResourceRequest) utils.GeneralResponse {
	// 校验资源ID格式
	if err := validateResourceID(data.ID); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}
	if err := validateResourceID(data.RscID); err != nil {
		return utils.GeneralResponse{Action: false, Error: fmt.Sprintf("invalid rsc_id: %v", err)}
	}

	// 删除原始资源的位置约束
	locationIds, err := getResourceConstraintIDs(data.RscID, "location")
	if err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}

	for _, id := range locationIds {
		_, err := getCommandRunner().RunCommandWithArgs("pcs", "constraint", "location", "delete", id)
		if err != nil {
			slog.Error("Delete location constraint of clone resource failed", "id", data.ID, "category", data.Category, "location-id", id, "err", err.Error())
			return utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Add clone resource failed"),
			}
		}
	}

	// TODO:删除资源tag
	// 删除克隆属性
	if err := DeleteCloneAttrib(data.RscID); err != nil {
		slog.Error("Delete attribute of clone resource failed", "id", data.ID, "category", data.Category, "rscId", data.RscID, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add clone resource failed"),
		}
	}

	// 创建克隆资源
	if _, err := getCommandRunner().RunCommandWithArgs("pcs", "resource", "clone", data.RscID); err != nil {
		slog.Error("Create clone resource failed", "id", data.ID, "category", data.Category, "rscId", data.RscID, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add clone resource failed"),
		}
	}

	// 更新克隆资源的属性
	if err := UpdateResourceAttributes(data.ID, data); err != nil {
		slog.Error("Update attribute of clone resource failed", "id", data.ID, "category", data.Category, "rscId", data.RscID, "err", err.Error())
		return utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Add clone resource failed"),
		}
	}

	// 检查是否有target-role，如果没有则禁用资源
	if _, exists := data.MetaAttributes["target-role"]; !exists {
		_, err := getCommandRunner().RunCommandWithArgs("pcs", "resource", "disable", data.ID)
		if err != nil {
			slog.Error("Update target-role of clone resource failed", "id", data.ID, "category", data.Category, "err", err.Error())
			return utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Add clone resource failed"),
			}
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Add " + data.Category + " resource success"),
	}
}

// 检查资源是否存在
func resourceExists(id string) (bool, error) {
	_, err := getCommandRunner().RunCommandWithArgs("crm_resource", "--locate", "--resource", id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 检查资源类型是否已存在
func resourceTypeExists(resType string) (bool, error) {
	out, err := getCommandRunner().RunCommand(utils.CmdResourceStatus)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), ":"+resType), nil
}

// 获取资源默认的monitor操作属性
func getDefaultMonitorAction(class string, provider string, resType string) Action {
	metas := GetResourceMetas(class, resType, provider)
	// 类型断言
	data, ok := metas["data"].(map[string]interface{})
	if !ok {
		slog.Error("Failed to get resource metas: 'data' field is not a map", "class", class, "provider", provider, "resType", resType, "data", data)

		return Action{}
	}
	rsc_actions, ok := data["actions"].(map[string]interface{})
	if !ok {
		slog.Error("Failed to get resource actions: 'actions' field is not a map", "class", class, "provider", provider, "resType", resType, "data", data)
		return Action{}
	}
	var monitorAction map[string]interface{}
	for _, action := range rsc_actions {
		actionMap, ok := action.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := actionMap["name"].(string); ok && name == "monitor" {
			monitorAction = actionMap
			break
		}
	}

	// 解析
	var defaultActions Action
	if timeout, ok := monitorAction["timeout"].(string); ok {
		defaultActions.Timeout = timeout
	}

	if interval, ok := monitorAction["interval"].(string); ok {
		defaultActions.Interval = interval
		if class == "stonith" {
			defaultActions.Interval = "1800s"
		}
	}

	if check_level, ok := monitorAction["OCF_CHECK_LEVEL"].(string); ok {
		defaultActions.OCFCheckLevel = check_level
	}

	return defaultActions
}

// UpdateResourceAttributes 更新资源属性
func UpdateResourceAttributes(rscId string, data ResourceRequest) error {
	if len(data.MetaAttributes) == 0 && len(data.InstanceAttributes) == 0 && len(data.Actions) == 0 && len(data.Rscs) == 0 {
		return errors.New(gettext.Gettext("No input data"))
	}

	if err := validateResourceID(rscId); err != nil {
		return err
	}

	exists, err := resourceExists(rscId)
	if err != nil {
		return fmt.Errorf("failed to check resource existence: %v", err)
	}
	if !exists {
		return errors.New(gettext.Gettext("Resource not found"))
	}

	// 1.按需清除现有属性（仅当对应字段有更新时才清除）
	if data.MetaAttributes != nil {
		if err := clearExistingMetaAttributes(rscId); err != nil {
			return err
		}
	}
	if data.InstanceAttributes != nil {
		if err := clearExistingInstanceAttributes(rscId); err != nil {
			return err
		}
	}

	// 2.更新元属性
	if err := updateMetaAttributes(rscId, data); err != nil {
		return err
	}

	// 3.更新实例属性
	if err := updateInstanceAttributes(rscId, data.InstanceAttributes); err != nil {
		return err
	}

	// 4.更新操作属性
	if err := updateOperations(rscId, data.Actions); err != nil {
		return err
	}

	// 5.处理组资源
	if data.Category == "group" {
		return updateGroupResources(rscId, data.Rscs)
	}

	return nil
}

// 清除资源现有属性
func clearExistingMetaAttributes(rscId string) error {
	attrib := GetMetaAndInst(rscId)
	if metaAttrs, ok := attrib["meta_attributes"]; ok {
		for _, attr := range metaAttrs {
			if _, err := utils.RunCommandWithArgs("crm_resource", "-r", rscId, "-m", "--delete-parameter", attr); err != nil {
				return fmt.Errorf("failed to clear meta attribute: %v", err)
			}
		}
	}
	return nil
}

func clearExistingInstanceAttributes(rscId string) error {
	attrib := GetMetaAndInst(rscId)
	if instAttrs, ok := attrib["instance_attributes"]; ok {
		for _, attr := range instAttrs {
			if _, err := utils.RunCommandWithArgs("crm_resource", "-r", rscId, "--delete-parameter", attr); err != nil {
				return fmt.Errorf("failed to clear instance attribute: %v", err)
			}
		}
	}
	return nil
}

// 更新元属性
func updateMetaAttributes(rscId string, data ResourceRequest) error {
	if data.MetaAttributes == nil {
		return nil
	}

	for key, val := range data.MetaAttributes {
		if val == nil {
			continue
		}
		if err := validateIdentifier(key, "meta attribute key"); err != nil {
			return err
		}
		strVal := fmt.Sprintf("%v", val)
		if strVal == "" {
			slog.Warn(fmt.Sprintf("Skipping meta attribute %s with unparsable value: %v", key, val))
			continue
		}

		var args []string
		if data.Category == "group" {
			args = []string{"resource", "meta", rscId, key + "=" + strVal}
		} else {
			args = []string{"resource", "update", rscId, "meta", key + "=" + strVal, "--force"}
		}

		if _, err := utils.RunCommandWithArgs("pcs", args...); err != nil {
			return fmt.Errorf("failed to set meta attribute %s: %v", key, err)
		}
	}
	return nil
}

// 更新实例属性
func updateInstanceAttributes(rscId string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}

	args := []string{"resource", "update", rscId}
	for key, val := range attributes {
		if err := validateIdentifier(key, "instance attribute key"); err != nil {
			return err
		}
		strVal := fmt.Sprintf("%v", val)
		if val == nil || strVal == "" {
			continue
		}
		if strVal == "" {
			slog.Warn(fmt.Sprintf("Skipping instance attribute %s with unparsable value: %v", key, val))
			continue
		}
		args = append(args, key+"="+strVal)
	}
	if len(args) > 3 {
		args = append(args, "--force")
		if _, err := utils.RunCommandWithArgs("pcs", args...); err != nil {
			return fmt.Errorf("failed to set instance attributes: %v", err)
		}
	}
	return nil
}

// 更新操作属性
func updateOperations(rscId string, actions []Action) error {
	if len(actions) == 0 {
		return nil
	}

	// 删除现有操作属性
	if err := deleteAllOperations(rscId); err != nil {
		return err
	}

	// 添加新操作
	for _, action := range actions {
		if err := addActionParam(rscId, action); err != nil {
			return err
		}
	}
	return nil
}

func addActionParam(rscId string, action Action) error {
	if err := validateIdentifier(action.Name, "action name"); err != nil {
		return err
	}

	args := []string{"resource", "op", "add", rscId, action.Name}

	// 使用反射动态处理所有可能的参数
	v := reflect.ValueOf(action)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Name" {
			continue
		}

		value := v.Field(i).String()
		if value != "" {
			paramName := strings.ToLower(strings.ReplaceAll(field.Name, "OCFCheckLevel", "OCF_CHECK_LEVEL"))
			switch paramName {
			case "ocf_check_level":
				paramName = "OCF_CHECK_LEVEL"
			case "startdelay":
				paramName = "start-delay"
			case "onfail":
				paramName = "on-fail"
			case "dep":
				paramName = "depth"
			}
			args = append(args, paramName+"="+value)
		}
	}
	args = append(args, "--force")
	if _, err := utils.RunCommandWithArgs("pcs", args...); err != nil {
		return fmt.Errorf("failed to add operation %s: %v", action.Name, err)
	}
	return nil
}

// 删除所有操作属性
func deleteAllOperations(rscId string) error {
	opList := GetAllOps(rscId)
	if len(opList) == 0 {
		slog.Warn("resource op list is empty", "id", rscId)
		return nil
	}
	var errs []error
	for _, op := range opList {
		if _, err := utils.RunCommandWithArgs("pcs", "resource", "op", "delete", rscId, op); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete operation %s: %v", op, err))
		}
	}
	return errors.Join(errs...)
}

// 更新组资源
func updateGroupResources(groupId string, resources []string) error {
	if resources == nil {
		return nil
	}

	currentPrimResourcesInGroup, err := getGroupRscs(groupId)
	if err != nil {
		return fmt.Errorf("failed to get current group resources: %v", err)
	}

	// step1： 添加不在当前组中的新资源
	for _, rsc := range resources {
		if !utils.Contains(currentPrimResourcesInGroup, rsc) {
			if err := addResourceToGroup(groupId, rsc); err != nil {
				return err
			}
		}
	}

	// step2：移除不再在组中的资源
	for _, rsc := range currentPrimResourcesInGroup {
		if !utils.Contains(resources, rsc) {
			if err := removeResourceFromGroup(groupId, rsc); err != nil {
				return err
			}
		}
	}

	return nil
}

// 添加资源到组
func addResourceToGroup(groupId, resourceId string) error {
	if err := validateResourceID(groupId); err != nil {
		return fmt.Errorf("invalid group id: %v", err)
	}
	if err := validateResourceID(resourceId); err != nil {
		return fmt.Errorf("invalid resource id: %v", err)
	}
	// 清除资源属性
	if err := DeletePriAttrib(resourceId); err != nil {
		return fmt.Errorf("failed to clear attributes for %s: %v", resourceId, err)
	}

	if _, err := utils.RunCommandWithArgs("pcs", "resource", "group", "add", groupId, resourceId); err != nil {
		return fmt.Errorf("failed to add %s to group %s: %v", resourceId, groupId, err)
	}
	return nil
}

// 从组中移除资源
func removeResourceFromGroup(groupId, resourceId string) error {
	if err := validateResourceID(groupId); err != nil {
		return fmt.Errorf("invalid group id: %v", err)
	}
	if err := validateResourceID(resourceId); err != nil {
		return fmt.Errorf("invalid resource id: %v", err)
	}
	if _, err := utils.RunCommandWithArgs("pcs", "resource", "group", "remove", groupId, resourceId); err != nil {
		return fmt.Errorf("failed to remove %s from group %s: %v", resourceId, groupId, err)
	}
	return nil
}

// 获取资源已有的操作属性列表(仅属性名，如start)
var GetAllOps = func(rscId string) []string {
	opList := []string{}
	if err := validateResourceID(rscId); err != nil {
		return opList
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscId, "--query-xml")
	if err != nil {
		return opList
	}
	xml, err2 := splitXMLOutput(out)
	if err2 != nil {
		return opList
	}
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

// 删除普通资源的属性
func DeletePriAttrib(rscId string) error {
	if err := validateResourceID(rscId); err != nil {
		return err
	}
	// delete attribute
	attrib := GetMetaAndInst(rscId)
	if metaAttri, ok := attrib["meta_attributes"]; ok {
		metaArr := metaAttri
		for _, v := range metaArr {
			if v == "is-managed" || v == "priority" || v == "target-role" {
				_, err := utils.RunCommandWithArgs("crm_resource", "-r", rscId, "-m", "--delete-parameter", v)
				if err != nil {
					return err
				}
			}
		}
	}
	// delete constraint
	// colocation
	targetId, err := getResourceConstraintIDs(rscId, "colocation")
	if err != nil {
		return err
	}
	err = DeleteColocationByIdAndAction(rscId, targetId)
	if err != nil {
		return err
	}
	// location
	ids, err := getResourceConstraintIDs(rscId, "location")
	if err != nil {
		return err
	}
	for _, item := range ids {
		_, err := utils.RunCommandWithArgs("pcs", "constraint", "location", "delete", item)
		if err != nil {
			return err
		}
	}
	// order
	hasOrder, err := findOrder(rscId)
	if err != nil {
		return err
	}
	if hasOrder {
		if _, err := utils.RunCommandWithArgs("pcs", "constraint", "order", "delete", rscId); err != nil {
			return err
		}
	}
	return nil
}

// 删除克隆资源的属性
func DeleteCloneAttrib(rscId string) error {
	if err := validateResourceID(rscId); err != nil {
		return err
	}
	// delete attribute
	attrib := GetMetaAndInst(rscId)
	if metaAttri, ok := attrib["meta_attributes"]; ok {
		metaArr := metaAttri
		for _, v := range metaArr {
			_, err := utils.RunCommandWithArgs("crm_resource", "-r", rscId, "-m", "--delete-parameter", v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var GetMetaAndInst = func(rscId string) map[string][]string {
	if err := validateResourceID(rscId); err != nil {
		return map[string][]string{}
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscId, "--query-xml")
	if err != nil {
		slog.Warn("Get meta and instance attributes failed", "id", rscId)
		return map[string][]string{}
	}
	xml, err2 := splitXMLOutput(out)
	if err2 != nil {
		return map[string][]string{}
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		return map[string][]string{}
	}
	data := map[string][]string{
		"meta_attributes":     []string{},
		"instance_attributes": []string{},
	}

	// 动态匹配三种资源类型
	/*
		xpathExpr := fmt.Sprintf("//*[(self::group or self::primitive or self::clone) and @id='%s']", rscId)
		resource := doc.FindElement(xpathExpr)
		if resource == nil {
			return data
		}


		// 遍历直接子节点
		for _, child := range resource.ChildElements() {
			switch child.Tag {
			case "meta_attributes":
				data["meta_attributes"] = collectNVPairs(child)
			case "instance_attributes":
				data["instance_attributes"] = collectNVPairs(child)
			}
		}*/
	resource := doc.ChildElements()
	for _, child := range resource {
		if child.Tag == "primitive" || child.Tag == "group" || child.Tag == "clone" {
			for _, subChild := range child.ChildElements() {
				switch subChild.Tag {
				case "meta_attributes":
					data["meta_attributes"] = collectNVPairs(subChild)
				case "instance_attributes":
					data["instance_attributes"] = collectNVPairs(subChild)
				}
			}
		}
	}
	return data

}

// 通用属性收集函数
func collectNVPairs(parent *etree.Element) []string {
	names := make([]string, 0)
	for _, nv := range parent.SelectElements("nvpair") {
		if nameAttr := nv.SelectAttr("name"); nameAttr != nil {
			names = append(names, nameAttr.Value)
		}
	}
	return names
}

func GetAllMigrateResources() []string {
	result := make([]string, 0)

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
		id := resourceLocation.SelectAttrValue("id", "")
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

var (
	failInfo          map[string]map[string]string
	clusterProperties map[string]interface{}
	globalStateMu     sync.RWMutex
)

func getFailInfo() map[string]map[string]string {
	globalStateMu.RLock()
	defer globalStateMu.RUnlock()
	return failInfo
}

func setFailInfo(fi map[string]map[string]string) {
	globalStateMu.Lock()
	defer globalStateMu.Unlock()
	failInfo = fi
}

func getGlobalClusterProperties() map[string]interface{} {
	globalStateMu.RLock()
	defer globalStateMu.RUnlock()
	return clusterProperties
}

func setGlobalClusterProperties(cp map[string]interface{}) {
	globalStateMu.Lock()
	defer globalStateMu.Unlock()
	clusterProperties = cp
}

var GetResourceStatus = func(rscInfo *etree.Element) string {
	rscId := rscInfo.SelectAttrValue("id", "")

	if cp := getGlobalClusterProperties(); cp != nil {
		if data, _ := cp["data"].(map[string]interface{}); data != nil {
			if params, _ := data["parameters"].(map[string]interface{}); params != nil {
				if mode, _ := params["maintenance-mode"].(map[string]interface{}); mode != nil {
					if value, ok := mode["value"].(string); ok && strings.EqualFold(value, "true") {
						return "Unmanaged"
					}
				}
			}
		}
	}
	if fi := getFailInfo(); fi != nil {
		if _, ok := fi[rscId]; ok {
			return GetRscStatusWithFailedInfo(rscInfo)
		}
	}

	if rscInfo.SelectAttrValue("managed", "") == "false" {
		return "Unmanaged"
	}
	if rscInfo.SelectAttrValue("failed", "") == "true" {
		return getRscStatusWithFailed(rscInfo)
	}
	if role := rscInfo.SelectAttr("role"); role != nil {
		if role.Value == "Started" {
			return "Running"
		}
		if role.Value == "Stopped" {
			return "Not Running"
		}

	}
	if disabled := rscInfo.SelectAttr("disabled"); disabled != nil {
		if disabled.Value == "true" {
			return "Not Running"
		}
	}
	return "Running"
}

func GetRscStatusWithFailedInfo(rscInfo *etree.Element) string {
	if rscInfo.SelectAttrValue("role", "") == "true" && rscInfo.SelectAttrValue("failed", "") == "true" {
		return "Failed"
	}
	if role := rscInfo.SelectAttr("role"); role != nil {
		switch strings.TrimSpace(role.Value) {
		case "Started":
			return "Running but failed"
		case "Stopped":
			return "Failed"
		}
	}
	return "Failed"
}

func getRscStatusWithFailed(rscInfo *etree.Element) string {
	if rscInfo.SelectAttrValue("blocked", "") == "true" {
		return "Failed"
	}
	if role := rscInfo.SelectAttr("role"); role != nil {
		if strings.TrimSpace(role.Value) == "Started" {
			return "Running but failed"
		}
	}
	return "Failed"
}

type ResourceParams struct {
	ID       string
	Managed  string
	Failed   string
	Role     string
	Disabled string
}

func getCloneSubrscPriStatus(rscInfo *etree.Element) string {
	// failInfo := GetResourceFailedMessage()
	params := ResourceParams{
		ID:       rscInfo.SelectAttrValue("id", ""),
		Managed:  rscInfo.SelectAttrValue("managed", ""),
		Failed:   rscInfo.SelectAttrValue("failed", ""),
		Role:     rscInfo.SelectAttrValue("role", ""),
		Disabled: rscInfo.SelectAttrValue("disabled", ""),
	}
	cachedFailInfo := getFailInfo()
	if _, ok := cachedFailInfo[params.ID]; ok {
		return GetCloneSubrscStatusWithFailedInfo(rscInfo, params.ID, cachedFailInfo)
	}
	if params.Managed == "false" {
		return "Unmanaged"
	}
	if params.Failed == "true" {
		return getRscStatusWithFailed(rscInfo)
	}
	switch params.Role {
	case "Started":
		return "Running"
	case "Stopped":
		return "Not Running"
	}
	if params.Disabled == "true" {
		return "Not Running"
	}
	return "Running"

}

func GetCloneSubrscStatusWithFailedInfo(rscInfo *etree.Element, rscId string, failInfo map[string]map[string]string) string {
	// failedNode = failInfo["rscId"]["node"]
	if failedRscId, ok := failInfo[rscId]; ok {
		// if !ok{
		// 	return "Failed"
		// }
		if failedNode, ok := failedRscId["node"]; ok {
			if role := rscInfo.SelectAttr("role"); role != nil {
				if rscInfo.SelectAttrValue("blocked", "") == "true" && rscInfo.SelectAttrValue("failed", "") == "true" {
					return "Failed"
				}
				if role.Value == "Started" {
					node := rscInfo.SelectAttr("node")
					rscRunningNode := node.Element().SelectAttrValue("name", "")
					if failedNode != rscRunningNode {
						return "Running"
					}
					return "Running but failed"
				}
			}
		}
	}
	return "Failed"
}

func GetResourceSvc(rscId string) string {
	if err := validateResourceID(rscId); err != nil {
		return ""
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscId, "--query-xml")
	if err != nil {
		slog.Error("GetResourceSvc: command failed", "id", rscId, "err", err)
		return ""
	}
	// Provide compatibility with different versions of Corosync
	xmlIndex := strings.Index(string(out), "XML:")
	if xmlIndex == -1 {
		xmlIndex = strings.Index(string(out), "xml:")
	}
	if xmlIndex == -1 {
		slog.Error("GetResourceSvc: XML marker not found in output", "id", rscId)
		return ""
	}
	xmlStr := string(out)[xmlIndex+len("XML:"):]

	doc := etree.NewDocument()
	if err = doc.ReadFromString(xmlStr); err != nil {
		slog.Error("GetResourceSvc: failed to parse XML", "id", rscId, "err", err)
		return ""
	}
	elem := doc.FindElement("primitive")
	if elem == nil {
		slog.Error("GetResourceSvc: primitive element not found", "id", rscId)
		return ""
	}
	return elem.SelectAttrValue("type", "")
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
		result = append(result, element.SelectAttrValue("id", ""))
	}
	elements = doc.FindElements("/resources/group")
	for _, element := range elements {
		result = append(result, element.SelectAttrValue("id", ""))
	}
	elements = doc.FindElements("/resources/primitive")
	for _, element := range elements {
		result = append(result, element.SelectAttrValue("id", ""))
	}

	return result
}

// ResourceAction 执行资源操作
func ResourceAction(rscID, action string, data []byte) error {
	if err := validateResourceID(rscID); err != nil {
		return err
	}
	// 处理资源ID中的冒号（防止clone资源名称问题）
	rscID = strings.Split(rscID, ":")[0]

	switch action {
	case "start", "stop":
		return handleStartStopAction(rscID, action)
	case "delete":
		return handleDeleteAction(rscID, data)
	case "cleanup":
		return handleCleanupAction(rscID)
	case "unclone":
		return handleUncloneAction(rscID)
	case "ungroup":
		return handleUngroupAction(rscID)
	case "migrate":
		return handleMigrateAction(rscID, data)
	case "unmigrate":
		return handleUnmigrateAction(rscID)
	case "location":
		return handleLocationAction(rscID, data)
	case "colocation":
		return handleColocationAction(rscID, data)
	case "order":
		return handleOrderAction(rscID, data)
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

// 处理启动/停止操作
func handleStartStopAction(rscID, action string) error {
	if action == "start" {
		_, err := utils.RunCommandWithArgs("pcs", "resource", "enable", rscID)
		return err
	}
	_, err := utils.RunCommandWithArgs("pcs", "resource", "disable", rscID)
	return err
}

// 处理删除操作
func handleDeleteAction(rscID string, data []byte) error {
	category := GetResourceCategory(rscID)

	switch category {
	case "group":
		_, err := utils.RunCommandWithArgs("pcs", "resource", "delete", rscID, "--force")
		return err
	case "clone":
		_, err := utils.RunCommandWithArgs("pcs", "resource", "delete", strings.TrimSuffix(rscID, "-clone"), "--force")
		return err
	default:
		// 处理guest资源
		if data != nil {
			var req struct {
				ResFlag string `json:"res_flag"`
			}
			if err := json.Unmarshal(data, &req); err == nil && req.ResFlag == "guest" {
				if _, err := utils.RunCommandWithArgs("pcs", "cluster", "node", "delete-guest", rscID); err != nil {
					return err
				}
				_, err := utils.RunCommandWithArgs("pcs", "resource", "delete", rscID, "--force")
				return err
			}
		}
		_, err := utils.RunCommandWithArgs("pcs", "resource", "delete", rscID, "--force")
		return err
	}
}

// 处理清理操作
func handleCleanupAction(rscID string) error {
	_, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscID, "--cleanup")
	return err
}

// 处理取消克隆操作
func handleUncloneAction(rscID string) error {
	_, err := utils.RunCommandWithArgs("pcs", "resource", "unclone", rscID)
	return err
}

// 处理解组资源操作
func handleUngroupAction(rscID string) error {
	_, err := utils.RunCommandWithArgs("pcs", "resource", "ungroup", rscID)
	return err
}

// @desperated
// 处理迁移操作
func handleMigrateAction(rscID string, data []byte) error {
	var req struct {
		IsForce bool   `json:"is_force"`
		ToNode  string `json:"to_node"`
		Period  string `json:"period"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("invalid migrate data: %v", err)
	}

	if req.ToNode == "" || !safeResourceName.MatchString(req.ToNode) {
		return fmt.Errorf("invalid target node name: %q", req.ToNode)
	}

	args := []string{"--resource", rscID, "--move", "-N", req.ToNode}

	if req.Period != "" {
		if !safePeriod.MatchString(req.Period) {
			return fmt.Errorf("invalid period value: %q", req.Period)
		}
		args = append(args, "--lifetime="+req.Period)
	}

	if req.IsForce {
		args = append(args, "--force")
	}

	out, err := utils.RunCommandWithArgs("crm_resource", args...)
	if err != nil {
		if strings.Contains(string(out), "Situation already as requested") {
			return fmt.Errorf("the resource %s is already running on node %s", rscID, req.ToNode)
		}
		return fmt.Errorf("migrate failed: %v", err)
	}

	return nil
}

// @desperated
// 处理取消迁移操作
func handleUnmigrateAction(rscID string) error {
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return fmt.Errorf("failed to query constraints: %v", err)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(out); err != nil {
		return fmt.Errorf("failed to parse constraints XML: %v", err)
	}

	found := false
	for _, elem := range doc.FindElements("//rsc_location") {
		if elem.SelectAttrValue("rsc", "") == rscID {
			locationID := elem.SelectAttrValue("id", "")
			if locationID != "" {
				if _, err := utils.RunCommandWithArgs("pcs", "constraint", "location", "delete", locationID); err != nil {
					return fmt.Errorf("failed to delete location constraint: %v", err)
				}
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("no migration constraints found for resource %s", rscID)
	}

	return nil
}

// 处理位置约束操作
func handleLocationAction(rscID string, data []byte) error {
	// 删除所有旧的位置约束
	ids, err := getResourceConstraintIDs(rscID, "location")
	if err != nil {
		return err
	}
	var delErrs []error
	for _, id := range ids {
		if _, err := utils.RunCommandWithArgs("pcs", "constraint", "location", "delete", id); err != nil {
			delErrs = append(delErrs, fmt.Errorf("failed to delete location constraint %s: %v", id, err))
		}
	}
	if err := errors.Join(delErrs...); err != nil {
		return err
	}

	// 解析新约束
	var req struct {
		NodeLevel []struct {
			Node  string `json:"node"`
			Level string `json:"level"`
		} `json:"node_level"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("invalid location data: %v", err)
	}

	// 创建新约束
	// {"node_level": [{"node": "ns187", "level": "Master Node"},
	for _, item := range req.NodeLevel {
		if item.Node == "" || !safeResourceName.MatchString(item.Node) {
			return fmt.Errorf("invalid node name: %q", item.Node)
		}

		score := getScoreForLevel(item.Level)
		if score == 0 {
			return fmt.Errorf("invalid node level: %s", item.Level)
		}

		if _, err := utils.RunCommandWithArgs("pcs", "constraint", "location", rscID, "prefers", item.Node+"="+strconv.Itoa(score)); err != nil {
			return fmt.Errorf("failed to set location constraint: %v", err)
		}
	}

	return nil
}

// 根据节点级别获取分值
func getScoreForLevel(level string) int {
	switch level {
	case "Master Node":
		return 20000
	case "Slave 1":
		return 16000
	case "Slave 2":
		return 15000
	case "Slave 3":
		return 14000
	case "Slave 4":
		return 13000
	default:
		return 0
	}
}

// 处理协同约束操作
func handleColocationAction(rscID string, data []byte) error {
	// 删除旧的协同约束
	ids, err := getResourceConstraintIDs(rscID, "colocation")
	if err != nil {
		return err
	}
	if err := DeleteColocationByIdAndAction(rscID, ids); err != nil {
		return fmt.Errorf("failed to delete old colocation constraints: %v", err)
	}

	// 解析新约束
	var req struct {
		SameNode []string `json:"same_node"`
		DiffNode []string `json:"diff_node"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("invalid colocation data: %v", err)
	}

	// 添加相同节点约束
	for _, item := range req.SameNode {
		if err := addColocationConstraint(rscID, item, true); err != nil {
			return err
		}
	}

	// 添加不同节点约束
	for _, item := range req.DiffNode {
		if err := addColocationConstraint(rscID, item, false); err != nil {
			return err
		}
	}

	return nil
}

// 添加协同约束
func addColocationConstraint(rscID, target string, sameNode bool) error {
	if err := validateResourceID(rscID); err != nil {
		return err
	}

	score := "INFINITY"
	if !sameNode {
		score = "-INFINITY"
	}

	var args []string
	// 处理克隆资源（with-rsc-role）
	if parts := strings.Split(target, "/"); len(parts) == 2 {
		if err := validateResourceID(parts[0]); err != nil {
			return fmt.Errorf("invalid target resource: %v", err)
		}
		if err := validateIdentifier(parts[1], "role"); err != nil {
			return err
		}
		args = []string{"constraint", "colocation", "add", rscID, "with", parts[0], score, "with-rsc-role=" + parts[1]}
	} else {
		if err := validateResourceID(target); err != nil {
			return fmt.Errorf("invalid target: %v", err)
		}
		args = []string{"constraint", "colocation", "add", rscID, "with", target, score}
	}

	if _, err := utils.RunCommandWithArgs("pcs", args...); err != nil {
		return fmt.Errorf("failed to add colocation constraint: %v", err)
	}
	return nil
}

// 处理顺序约束操作
func handleOrderAction(rscID string, data []byte) error {
	// 删除旧的顺序约束
	hasOrder, err := findOrder(rscID)
	if err != nil {
		return err
	}
	if hasOrder {
		if _, err := utils.RunCommandWithArgs("pcs", "constraint", "order", "delete", rscID); err != nil {
			return fmt.Errorf("failed to delete old order constraints: %v", err)
		}
	}

	//添加新约束
	var jsonData []interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	for _, d := range jsonData {
		m, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		beforeAction, ok := m["before_action"].(string)
		if !ok || beforeAction == "" {
			return fmt.Errorf("missing or invalid before_action")
		}
		afterAction, ok := m["after_action"].(string)
		if !ok || afterAction == "" {
			return fmt.Errorf("missing or invalid after_action")
		}
		rscName, ok := m["rsc_name"].(string)
		if !ok || rscName == "" {
			return fmt.Errorf("missing or invalid rsc_name")
		}

		if err := validateIdentifier(beforeAction, "before_action"); err != nil {
			return err
		}
		if err := validateIdentifier(afterAction, "after_action"); err != nil {
			return err
		}
		if err := validateResourceID(rscName); err != nil {
			return err
		}

		var args []string
		switch m["location"] {
		case "before":
			args = []string{"constraint", "order", beforeAction, rscName, "then", afterAction, rscID}
		case "after":
			args = []string{"constraint", "order", beforeAction, rscID, "then", afterAction, rscName}
		default:
			continue
		}
		if _, err := utils.RunCommandWithArgs("pcs", args...); err != nil {
			return err
		}
	}
	return nil
}
func getResourceConstraintIDs(rscID, action string) ([]string, error) {
	ids := []string{}
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return nil, fmt.Errorf("failed to query constraints: %v", err)
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return nil, fmt.Errorf("failed to parse constraints XML: %v", err)
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
		return ids, nil
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
		return ids, nil
	}
	return ids, nil
}

// 删除sourceID有关的所有的colocation关系
func DeleteColocationByIdAndAction(sourceID string, targetIds []string) error {
	if err := validateResourceID(sourceID); err != nil {
		return err
	}
	for _, item := range targetIds {
		if err := validateResourceID(item); err != nil {
			return fmt.Errorf("invalid target id: %v", err)
		}
		if _, err := utils.RunCommandWithArgs("pcs", "constraint", "colocation", "delete", sourceID, item); err != nil {
			return err
		}
	}
	return nil
}

func findOrder(rscID string) (bool, error) {
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return false, fmt.Errorf("failed to query constraints: %v", err)
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return false, fmt.Errorf("failed to parse constraints XML: %v", err)
	}
	et := doc.FindElements("/constraints/rsc_order")
	for _, item := range et {
		first := item.SelectAttrValue("first", "")
		then := item.SelectAttrValue("then", "")
		if first == rscID || then == rscID {
			return true, nil
		}
	}
	return false, nil
}

func GetResourceInfoByrscID(rscID string) (interface{}, error) {
	if err := validateResourceID(rscID); err != nil {
		return nil, err
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", rscID, "--query-xml")
	if err != nil {
		return nil, err
	}

	xml, err := splitXMLOutput(out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command output: %w", err)
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xml); err != nil {
		return nil, err
	}
	root := doc.Root()

	ct := root.Tag
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
	slog.Debug("Get resource info by id")
	doc := etree.NewDocument()
	doc.ReadFromString(xmlData)
	data := map[string]interface{}{}
	if err := doc.ReadFromString(xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Format data to map here
	switch ct {
	case "primitive":
		d, err := getResourceInfoFromXml("primitive", doc.Root())
		if err != nil {
			return nil, err
		}
		info, ok := d.(PrimitiveResource)
		if !ok {
			return nil, fmt.Errorf("unexpected type for primitive resource: %T", d)
		}
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
		info, ok := d.(GroupResource)
		if !ok {
			return nil, fmt.Errorf("unexpected type for group resource: %T", d)
		}
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

		var info CloneResource
		jsonData, _ := json.Marshal(d)
		if err := json.Unmarshal(jsonData, &info); err != nil {
			return nil, fmt.Errorf("unmarshal error: %v", err)
		}

		data["id"] = info.ID

		// TODO: check if only one Primitive resource or list
		rscs := []string{}
		for _, p := range info.Primitives {
			rscs = append(rscs, p.ID)
		}

		if len(rscs) == 1 {
			data["rsc_id"] = rscs[0]
		} else {
			data["rsc_id"] = rscs
		}
	}

	// For meta_attributes
	e := doc.FindElement("/" + ct + "/meta_attributes")
	if e != nil {
		prop, _ := getResourceInfoFromXml("meta", e)
		if m, ok := prop.(map[string]string); ok && len(m) > 0 {
			data["meta_attributes"] = m
		}
	}

	//For instance_attributes
	e = doc.FindElement("/" + ct + "/instance_attributes")
	if e != nil {
		prop, _ := getResourceInfoFromXml("inst", e)
		if m, ok := prop.(map[string]string); ok && len(m) > 0 {
			data["instance_attributes"] = m
		}
	}

	//For actions
	e = doc.FindElement("/" + ct + "/operations")
	if e != nil {
		prop, _ := getResourceInfoFromXml("operations", e)
		if ops, ok := prop.([]map[string]string); ok && len(ops) > 0 {
			data["actions"] = ops
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
	ID         string              `json:"id"`
	Primitives []PrimitiveResource `json:"Primitives"`
}

type CloneResource struct {
	ID         string              `json:"id"`
	Primitives []PrimitiveResource `json:"Primitives"`
}

type BundleResource struct {
	ID          string             `json:"id"`
	DockerImage string             `json:"docker_image"`
	Replicas    string             `json:"replicas"`
	Primitive   *PrimitiveResource `json:"primitive,omitempty"`
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

		// 针对普通资源
		els := et.FindElements("primitive")

		// 针对组资源
		if len(els) == 0 {
			groupElem := et.FindElement("group")
			if groupElem == nil {
				return rsc, nil
			}
			return getResourceInfoFromXml("group", groupElem)
		}

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
			name := item.SelectAttrValue("name", "")
			value := item.SelectAttrValue("value", "")
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
			other := item.FindElement(".//nvpair")
			if other != nil {
				i[other.SelectAttrValue("name", "")] = other.SelectAttrValue("value", "")
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
	if err := validateResourceID(groupId); err != nil {
		return nil, err
	}
	out, err := utils.RunCommandWithArgs("crm_resource", "--resource", groupId, "--query-xml")

	if err != nil {
		// result := map[string]interface{}{}
		// result["action"]=
		// return result
		return nil, err
	}
	xml, err2 := splitXMLOutput(out)
	if err2 != nil {
		return nil, err2
	}
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
	if nodeNum <= 1 {
		return []string{}
	}
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
