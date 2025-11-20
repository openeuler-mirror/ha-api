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
	disksListDataTemp, _ := utils.RunCommand(cmdGetDisksList)
	if len(disksListDataTemp) != 0 {
		disksListData = strings.Split(strings.TrimSuffix(string(disksListDataTemp), "\n"), "\n")
	}
	var diskname string
	var diskcap int
	diskslist := make([]map[string]interface{}, 0)
	for _, value := range disksListData {
		disklist := make(map[string]interface{})
		diskname = strings.Split(value, " ")[0]
		diskcap, _ = strconv.Atoi(strings.Split(value, " ")[1])
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