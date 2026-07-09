/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Fri Jul 04 15:54:28 2025 +0800
 */

package models

import (
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

const ATTRD_UPDATER = "attrd_updater -n "

type HealthGetResponse struct {
	Action bool       `json:"action"`
	Data   HealthData `json:"data"`
}

type HealthData struct {
	CpuFlag    bool                     `json:"cpu_flag"`
	CpuYellow  string                   `json:"cpu_yellow"`
	CpuRed     string                   `json:"cpu_red"`
	MemFlag    bool                     `json:"mem_flag"`
	MemYellow  string                   `json:"mem_yellow"`
	MemRed     string                   `json:"mem_red"`
	DiskFlag   bool                     `json:"disk_flag"`
	DiskYellow string                   `json:"disk_yellow"`
	DiskRed    string                   `json:"disk_red"`
	Disks      string                   `json:"disks"`
	DisksList  []map[string]interface{} `json:"diskslist"`
}

// 获取属性值返给前端展示
func HealthGet() HealthGetResponse {
	var healthDataResponse HealthGetResponse
	var healthData HealthData
	healthDataResponse.Action = true

	disksListData := []string{}
	healthData.CpuYellow = getIndexValue("health_cpu-clone", "yellow_limit")
	healthData.CpuRed = getIndexValue("health_cpu-clone", "red_limit")
	healthData.MemYellow = getIndexValue("health_mem-clone", "yellow_mem")
	healthData.MemRed = getIndexValue("health_mem-clone", "red_mem")
	healthData.DiskYellow = getIndexValue("sysinfo-clone", "yellow_disk_free")
	healthData.DiskRed = getIndexValue("sysinfo-clone", "min_disk_free")
	healthData.Disks = getIndexValue("sysinfo-clone", "disks")

	cmdGetDisksList := "df --block-size=1GB --output=source,avail | awk -F \" \" '{print $1,$2}'|grep ^/dev"
	disksListDataTemp, err := utils.RunCommand(cmdGetDisksList)
	if err != nil {
		slog.Error("get disk list failed", "cmd", cmdGetDisksList, "error", err)
	} else if len(disksListDataTemp) != 0 {
		disksListData = strings.Split(strings.TrimSuffix(string(disksListDataTemp), "\n"), "\n")
	}
	var diskname string
	var diskcap int
	diskslist := make([]map[string]interface{}, 0)
	for _, value := range disksListData {
		disklist := make(map[string]interface{})
		diskname = strings.Split(value, " ")[0]
		diskcap, err = strconv.Atoi(strings.Split(value, " ")[1])
		if err != nil {
			slog.Warn("parse disk capacity failed", "value", value, "error", err)
		}
		disklist["key"] = diskname
		disklist["value"] = diskcap
		diskslist = append(diskslist, disklist)
	}

	if healthData.CpuYellow != "" && healthData.CpuRed != "" {
		healthData.CpuFlag = true
	} else {
		healthData.CpuFlag = false
	}

	if healthData.MemYellow != "" && healthData.MemRed != "" {
		healthData.MemFlag = true
	} else {
		healthData.MemFlag = false
	}

	if healthData.DiskYellow != "" && healthData.DiskRed != "" {
		healthData.DiskFlag = true
		healthData.DisksList = diskslist
	} else {
		healthData.DiskFlag = false
		healthData.DisksList = diskslist
	}
	healthDataResponse.Data = healthData
	return healthDataResponse
}

// 获取cib文件中智能迁移资源对应的属性的值
func getIndexValue(resName, indexName string) string {
	getIndex := "cibadmin --query --xpath " + "\"//clone[@id='" + resName + "']//primitive//instance_attributes//nvpair[@name='" + indexName + "']\""
	_, err := utils.RunCommand(getIndex)
	if err == nil {
		value := getIndex + " | awk -F 'value=\"|\"/>' '{print $2}'"
		res, _ := utils.RunCommand(value)
		// 不加strings.TrimSpace在go的返回有\n
		return strings.TrimSpace(string(res))
	} else {
		return ""
	}
}

func HealthSet(data []byte) utils.GeneralResponse {
	var result utils.GeneralResponse

	if len(data) == 0 {
		result.Action = false
		result.Error = gettext.Gettext("No input data")
		return result
	}

	healthData := HealthData{}
	err := json.Unmarshal(data, &healthData)
	if err != nil {
		result.Action = false
		result.Error = gettext.Gettext("Cannot convert data to json map")
		return result
	}

	healthDeleteList, _ := GetResource()
	res := healthDelete(healthDeleteList)
	if res != "" {
		result.Action = false
		result.Error = gettext.Gettext("Failed to delete " + res + " resource!")
		return result
	}

	healthDeleteList, _ = GetResource()
	if len(healthDeleteList) == 0 {
		createResF, createResT := []string{}, []string{} // 创建资源列表
		updateResF, updateResT := []string{}, []string{} // 更新资源列表
		deleteResF, deleteResT := []string{}, []string{} // 删除资源列表

		if healthData.CpuFlag == true || healthData.MemFlag == true || healthData.DiskFlag == true {
			//设置健康策略
			getProperty := "cibadmin --query --scope crm_config | grep node-health-strategy |awk -F 'value=|/>' '{print $2}'"
			getPropertyData, _ := utils.RunCommand(getProperty)
			if string(getPropertyData) != "migrate-on-red" && string(getPropertyData) != "only-green" {
				setProperty := "pcs property set node-health-strategy=migrate-on-red"
				_, _ = utils.RunCommand(setProperty)
			}
		}

		if healthData.CpuFlag == true {
			//创建或更新资源
			crmResourceCpu := "crm_resource --resource health_cpu-clone --query-xml "
			_, errCpu := utils.RunCommand(crmResourceCpu)
			if errCpu != nil {
				//创建资源
				createCpu := "pcs resource create health_cpu ocf:pacemaker:HealthCPU" + " " + "yellow_limit=" + healthData.CpuYellow + " " + "red_limit=" + healthData.CpuRed + " " + " clone" + " " + "&& pcs resource meta health_cpu-clone allow-unhealthy-nodes=true"
				_, createCpuErr := utils.RunCommand(createCpu)

				if createCpuErr == nil {
					createResT = append(createResT, "cpu")
				} else {
					createResF = append(createResF, "cpu")
				}
			} else {
				//更新资源
				updateCpu := "pcs resource update health_cpu ocf:pacemaker:HealthCPU" + " " + "yellow_limit=" + healthData.CpuYellow + " " + "red_limit=" + healthData.CpuRed + " " + " clone" + " " + "&& pcs resource meta health_cpu-clone allow-unhealthy-nodes=true"
				_, updateCpuErr := utils.RunCommand(updateCpu)
				if updateCpuErr == nil {
					updateResT = append(updateResT, "cpu")
					attrdUpdateGreen("#health-cpu")
				} else {
					updateResF = append(updateResF, "cpu")
				}
			}
		} else {
			// healthData.CpuFlag == "False"
			// 删除资源
			crmResourceCpu := "crm_resource --resource health_cpu-clone --query-xml "
			_, errCpu := utils.RunCommand(crmResourceCpu)
			if errCpu == nil {
				deleteCpu := "pcs resource delete health_cpu-clone "
				_, deleteCpuErr := utils.RunCommand(deleteCpu)
				if deleteCpuErr == nil {
					deleteResT = append(deleteResT, "cpu")
					attrdUpdateGreen("#health-cpu")
				} else {
					deleteResF = append(deleteResF, "cpu")
				}
			}
		}

		if healthData.MemFlag == true {
			//处理内存：创建mem资源
			crmResourceMem := "crm_resource --resource health_mem-clone --query-xml "
			_, errMem := utils.RunCommand(crmResourceMem)
			if errMem != nil {
				createMem := "pcs resource create health_mem ocf:pacemaker:HealthMEM" + " " + "yellow_mem=" + healthData.MemYellow + " " + "red_mem=" + healthData.MemRed + " " + "clone" + " " + "&& pcs resource meta health_mem-clone allow-unhealthy-nodes=true"
				_, createMemErr := utils.RunCommand(createMem)
				if createMemErr == nil {
					createResT = append(createResT, "mem")
				} else {
					createResF = append(createResF, "mem")
				}
			} else {
				updateMem := "pcs resource update health_mem ocf:pacemaker:HealthMEM" + " " + "yellow_mem=" + healthData.MemYellow + " " + "red_mem=" + healthData.MemRed + " " + "clone" + " " + "&& pcs resource meta health_mem-clone allow-unhealthy-nodes=true"
				_, updateMemErr := utils.RunCommand(updateMem)
				if updateMemErr == nil {
					updateResT = append(updateResT, "mem")
					attrdUpdateGreen("#health-mem")
				} else {
					updateResF = append(updateResF, "mem")
				}
			}
		} else {
			// healthData.MemFlag == "False"
			crmResourceMem := "crm_resource --resource health_mem-clone --query-xml "
			_, errMem := utils.RunCommand(crmResourceMem)
			if errMem == nil {
				deleteMem := "pcs resource delete health_mem-clone"
				_, deleteMemErr := utils.RunCommand(deleteMem)
				if deleteMemErr == nil {
					deleteResT = append(deleteResT, "mem")
					attrdUpdateGreen("#health-mem")
				} else {
					deleteResF = append(deleteResF, "mem")
				}
			}
		}

		if healthData.DiskFlag == true {
			//创建或更新资源
			crmResourceDisk := "crm_resource  --resource sysinfo-clone --query-xml "
			_, errDisk := utils.RunCommand(crmResourceDisk)
			if errDisk != nil {
				//创建资源
				createDisk := "pcs resource create sysinfo ocf:pacemaker:SysInfo" + " " + "disks=" + healthData.Disks + " " + "min_disk_free=" + healthData.DiskRed + " " + "yellow_disk_free=" + healthData.DiskYellow + " " + "clone" + " " + "&& pcs resource meta sysinfo-clone allow-unhealthy-nodes=true"
				_, createDiskErr := utils.RunCommand(createDisk)
				if createDiskErr == nil {
					createResT = append(createResT, "disk")
				} else {
					createResF = append(createResF, "disk")
				}
			} else {
				//更新资源
				updateDisk := "pcs resource update sysinfo ocf:pacemaker:SysInfo" + " " + "disks=" + healthData.Disks + " " + "min_disk_free=" + healthData.DiskRed + " " + "yellow_disk_free=" + healthData.DiskYellow + " " + "clone" + " " + "&& pcs resource meta sysinfo-clone allow-unhealthy-nodes=true"
				_, updateDiskErr := utils.RunCommand(updateDisk)
				if updateDiskErr == nil {
					updateResT = append(updateResT, "disk")
					attrdUpdateGreen("#health_disk")
				} else {
					updateResF = append(updateResF, "disk")
				}
			}
		} else {
			//删除资源
			crmResourceDisk := "crm_resource  --resource sysinfo-clone --query-xml "
			_, errDisk := utils.RunCommand(crmResourceDisk)
			if errDisk == nil {
				deleteDisk := "pcs resource delete  sysinfo-clone"
				_, deleteDiskErr := utils.RunCommand(deleteDisk)
				if deleteDiskErr == nil {
					deleteResT = append(deleteResT, "disk")
					attrdUpdateGreen("#health_disk")
				} else {
					deleteResF = append(deleteResF, "disk")
				}
			}
		}

		//总结创建、更新、删除资源结果，并返回到前端
		createConT := ""
		createOutT := ""
		if len(createResT) != 0 {
			for _, value := range createResT {
				createConT = createConT + " " + value
			}
			createOutT = gettext.Gettext("Succeed to create health indicator") + createConT + " "
		}

		createConF := ""
		createOutF := ""
		if len(createResF) != 0 {
			for _, value := range createResF {
				createConF = createConF + " " + value
			}
			createOutF = gettext.Gettext("Failed to create health indicator") + createConF + " "
		}

		updateConT := ""
		updateOutT := ""
		if len(updateResT) != 0 {
			for _, value := range updateResT {
				updateConT = updateConT + " " + value
			}
			updateOutT = gettext.Gettext("Succeed to update health indicator") + updateConT + " "
		}

		updateConF := ""
		updateOutF := ""
		if len(updateResF) != 0 {
			for _, value := range updateResF {
				updateConF = updateConF + " " + value
			}
			updateOutF = gettext.Gettext("Failed to update health indicator") + updateConF + " "
		}

		deleteConT := ""
		deleteOutT := ""
		if len(deleteResT) != 0 {
			for _, value := range deleteResT {
				deleteConT = deleteConT + " " + value
			}
			deleteOutT = gettext.Gettext("Succeed to delete health indicator") + deleteConT + " "
		}

		deleteConF := ""
		deleteOutF := ""
		if len(deleteResF) != 0 {
			for _, value := range deleteResF {
				deleteConF = deleteConF + " " + value
			}
			deleteOutF = gettext.Gettext("Failed to delete health indicator") + deleteConF + " "
		}

		if createOutF != "" || updateOutF != "" || deleteOutF != "" {
			result.Action = false
			result.Error = createOutF + updateOutF + deleteOutF + createOutT + updateOutT + deleteOutT
		} else {
			result.Action = true
			result.Info = gettext.Gettext("Smart migration Settings are successful!")
		}
	}
	return result
}

// 更新属性为green
func attrdUpdateGreen(indexName string) {
	res := "crm_node -l | awk -F ' ' '{print $2}'"
	nodes, _ := utils.RunCommand(res)
	nodesName := strings.Split(strings.TrimSuffix(string(nodes), "\n"), "\n")
	for _, nodeName := range nodesName {
		attrd := ATTRD_UPDATER + "\"" + indexName + "\"" + " -U  \"green\" --node " + nodeName
		_, _ = utils.RunCommand(attrd)
	}
}

func GetHealthInfo(nodeName string) string {
	var healthStatus string
	var pronInuse []string

	cmdGetProperty := "cibadmin --query --scope crm_config | grep node-health-strategy |awk -F 'value=\"|\"/>' '{print $2}'"
	pro, _ := utils.RunCommand(cmdGetProperty)
	if string(pro) == "none" {
		healthStatus = "healthy"
	} else {
		pronInuse = getHealthInuseList(nodeName)
		if len(pronInuse) == 0 {
			healthStatus = "healthy"
		} else {
			healthStatus = getHealthInuseNonull(strings.TrimSpace(string(pro)), nodeName, pronInuse)
		}
	}

	return healthStatus
}

func getHealthInuseNonull(pro, nodeName string, proInuse []string) string {
	var healthStatus string
	healthStatus = "healthy"
	for _, value := range proInuse {
		attrdQ := ATTRD_UPDATER + value + " -Q --node " + nodeName + " | awk -F 'value=\"' '{print $2}' | awk -F '\"' '{print $1}'"
		attrdQRes, err := utils.RunCommand(attrdQ)

		if pro == "migrate-on-red" && err == nil && strings.TrimSpace(string(attrdQRes)) == "red" {
			healthStatus = "unhealthy"
			break
		}
		if pro == "only-green" && err == nil && (strings.TrimSpace(string(attrdQRes)) == "yellow" || strings.TrimSpace(string(attrdQRes)) == "red") {
			healthStatus = "unhealthy"
			break
		}
	}

	return healthStatus
}

func getHealthInuseList(nodeName string) []string {
	var pronInuse []string

	healthType := make(map[string]string)
	healthType["SysInfo"] = "\"#health_disk\""
	healthType["HealthMEM"] = "\"#health-mem\""
	healthType["HealthCPU"] = "\"#health-cpu\""

	cmdGetType := "cibadmin -Q -o resources | xmllint --xpath '//primitive/@type' - | awk -F 'type=\"|\"' '{print $2}'"
	getType, _ := utils.RunCommand(cmdGetType)

	for key, value := range healthType {
		if strings.Contains(string(getType), key) {
			pronInuse = append(pronInuse, value)
		} else {
			attrdCg := ATTRD_UPDATER + value + " -U  \"green\" --node " + nodeName
			_, _ = utils.RunCommand(attrdCg)
		}
	}
	return pronInuse
}

type HealthTestData struct {
	Action bool     `json:"action"`
	Health []string `json:"health"`
}

func HealthTest() HealthTestData {
	var healthTestData HealthTestData

	res := []string{}
	listType := []string{"HealthCPU", "HealthMEM", "SysInfo"}
	healthDeleteList, _ := GetResource()
	for _, value := range healthDeleteList {
		if utils.Contains(listType, value) && !utils.Contains(res, value) {
			res = append(res, value)
		}
	}

	healthTestData.Action = true
	healthTestData.Health = res
	return healthTestData
}

func healthDelete(healthList map[string]string) string {
	leap := ""
	if len(healthList) != 0 {
		for key, value := range healthList {
			healthCmd := "pcs resource delete " + key + " --force"
			_, err := utils.RunCommand(healthCmd)
			if err != nil {
				leap = key
			} else {
				attrdUpdateGreen(value)
			}
		}
	}
	return leap
}

func GetResource() (map[string]string, map[string]string) {
	healthDeleteList := map[string]string{}
	healthRevList := map[string]string{}

	cmdGetType := "cibadmin -Q -o resources | xmllint --xpath '//primitive/@type' - | awk -F 'type=\"|\"' '{print $2}'"
	cmdGetResource := "cibadmin -Q -o resources | xmllint --xpath '//primitive/@id' - | awk -F 'type=\"|\"' '{print $2}'"
	getType, _ := utils.RunCommand(cmdGetType)
	getResource, _ := utils.RunCommand(cmdGetResource)
	getTypeRes := utils.SplitLinesScanner(string(getType))
	getResourceRes := utils.SplitLinesScanner(string(getResource))

	resourceDict := map[string]string{}
	for i, k := range getResourceRes {
		if i < len(getTypeRes) {
			resourceDict[k] = getTypeRes[i]
		}
	}

	listType := []string{"HealthCPU", "HealthMEM", "SysInfo"}
	resBak := []string{}
	for key, value := range resourceDict {
		if utils.Contains(listType, value) {
			resBak = append(resBak, key)
		}
	}

	cmdGetSingleName := "crm_resource | grep 'Clone Set:' | awk -F '[' '{print $2}' | awk -F ']:' '{print $1}'"
	getSingleName, _ := utils.RunCommand(cmdGetSingleName)
	singleName := utils.SplitLinesScanner(string(getSingleName))

	listName := []string{"sysinfo", "health_mem", "health_cpu"}
	listTest := []string{}

	for _, single := range singleName {
		if utils.Contains(listName, single) {
			listTest = append(listTest, single)
		}
	}

	for _, i := range listTest {
		for key, value := range resourceDict {
			if i == key && ((key == "sysinfo" && value == "SysInfo") || (key == "health_mem" && value == "HealthMEM") || (key == "health_cpu" && value == "HealthCPU")) {
				resBak = removeAll(resBak, i)
				healthRevList[key] = value
			}
		}
	}

	for key, value := range resourceDict {
		if utils.Contains(resBak, key) {
			healthDeleteList[key] = value
		}
	}
	return healthDeleteList, healthRevList
}

func removeAll(slice []string, value string) []string {
	result := slice[:0] // 复用原切片的底层数组
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
