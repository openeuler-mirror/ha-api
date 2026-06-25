/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Tue Jun 23 15:54:28 2026 +0800
 */

package models

// 资源模板包含的资源种类、每个种类包含的资源，可修改
var category_default = []string{"DISK", "FileSystem", "Database", "Middleware", "VIP"}
var category_DISK_default = []string{"LVM"}
var category_FileSystem_default = []string{"Filesystem"}
var category_Database_default = []string{"DMDB8"}
var category_Middleware_default = []string{"TongWeb8", "BES9.5"}
var category_VIP_default = []string{"IPaddr2", "IPaddr"}

// 每个资源包含的属性及默认值，可修改
var res_LVM_default = map[string]interface{}{
	"volgrpname": "",
	"tag":        "",
}

var res_Filesystem_default = map[string]interface{}{
	"device":    "",
	"directory": "",
	"fstype":    "",
	"options":   "",
}

var res_DMDB8_default = map[string]interface{}{
	"datadir":     "/opt/dmdbms/",
	"instancedir": "/opt/dmdbms/data/DAMENG",
}

var res_TongWeb8_default = map[string]interface{}{
	"tongweb_path": "",
}

var res_BES9_default = map[string]interface{}{
	"BES_HOME":        "/opt/BES",
	"iastoolUser":     "",
	"iastoolPassword": "",
	"JAVA_HOME":       "",
}

var res_IPaddr_default = map[string]interface{}{
	"ip":           "",
	"nic":          "eth0",
	"cidr_netmask": "",
	"broadcast":    "",
}

var res_IPaddr2_default = map[string]interface{}{
	"ip":           "",
	"nic":          "",
	"cidr_netmask": "",
	"broadcast":    "",
}

type ResourceModelGetResponse struct {
	Action bool              `json:"action"`
	Data   ResourceModelData `json:"data"`
}

type ResourceModelData struct {
	Category_DISK       []string `json:"category_disk"`
	Category_FileSystem []string `json:"category_filesystem"`
	Category_Database   []string `json:"category_database"`
	Category_Middleware []string `json:"category_middleware"`
	Category_VIP        []string `json:"category_vip"`
}

func ResourceModelGet() ResourceModelGetResponse {
	var resourceModelGetResonpe ResourceModelGetResponse
	var resourceModelData ResourceModelData

	resourceModelData.Category_DISK = category_DISK_default
	resourceModelData.Category_FileSystem = category_FileSystem_default
	resourceModelData.Category_Database = category_Database_default
	resourceModelData.Category_Middleware = category_Middleware_default
	resourceModelData.Category_VIP = category_VIP_default

	resourceModelGetResonpe.Action = true
	resourceModelGetResonpe.Data = resourceModelData
	return resourceModelGetResonpe
}
