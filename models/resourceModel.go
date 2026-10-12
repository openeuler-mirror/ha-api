/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Tue Jun 23 15:54:28 2026 +0800
 */

package models

import (
	"fmt"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

// ==================== 配置定义（唯一数据源） ====================

// CategoryConfig 存储所有分类及其子项配置
var CategoryConfig = map[string]Category{
	"DISK": {
		SubItems:   []string{"LVM"},
		DefaultLen: 1,
		ResourceDefaults: map[string]map[string]interface{}{
			"LVM": {
				"id":         "",
				"volgrpname": "",
				"tag":        "",
			},
		},
	},
	"FileSystem": {
		SubItems:   []string{"Filesystem"},
		DefaultLen: 1,
		ResourceDefaults: map[string]map[string]interface{}{
			"Filesystem": {
				"id":        "",
				"device":    "",
				"directory": "",
				"fstype":    "",
				"options":   "",
			},
		},
	},
	"Database": {
		SubItems:   []string{"DMDB8"},
		DefaultLen: 1,
		ResourceDefaults: map[string]map[string]interface{}{
			"DMDB8": {
				"id":          "",
				"datadir":     "/opt/dmdbms/",
				"instancedir": "/opt/dmdbms/data/DAMENG",
			},
		},
	},
	"Middleware": {
		SubItems:   []string{"TongWeb8", "BES9.5"},
		DefaultLen: 1,
		ResourceDefaults: map[string]map[string]interface{}{
			"TongWeb8": {
				"id":           "",
				"tongweb_path": "",
			},
			"BES9.5": {
				"id":              "",
				"BES_HOME":        "/opt/BES",
				"iastoolUser":     "",
				"iastoolPassword": "",
				"JAVA_HOME":       "",
			},
		},
	},
	"VIP": {
		SubItems:   []string{"IPaddr2", "IPaddr"},
		DefaultLen: 1,
		ResourceDefaults: map[string]map[string]interface{}{
			"IPaddr2": {
				"id":           "",
				"ip":           "",
				"nic":          "",
				"cidr_netmask": "",
				"broadcast":    "",
			},
			"IPaddr": {
				"id":           "",
				"ip":           "",
				"nic":          "eth0",
				"cidr_netmask": "",
				"broadcast":    "",
			},
		},
	},
}

// Category 表示一个资源分类
type Category struct {
	SubItems         []string                          // 该分类下的资源类型列表
	DefaultLen       int                               // 默认部署数量
	ResourceDefaults map[string]map[string]interface{} // 各资源类型的属性默认值
}

type ResourceModelDeployResponse struct {
	Content   string              // 首页展示内容
	GroupName string              // 资源组名
	Category  map[string]Category // 所有资源信息
}

func ResourceModelDeploy() ResourceModelDeployResponse {
	return ResourceModelDeployResponse{
		Content:   "this is resource model",
		GroupName: "",
		Category:  CategoryConfig,
	}
}

type ResourceModelRequest struct {
	GroupName string              `json:"group_name"`
	Data      []ResourceModelData `json:"data"`
	IsTrue    bool                `json:"isTrue"`
}

type ResourceModelPostResponse struct {
	Action bool        `json:"action,omitempty"`
	Error  interface{} `json:"error,omitempty"`
	Info   interface{} `json:"info,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type ResourceModelData struct {
	CategoryKey  string                 `json:"categoryKey"`
	CategoryName string                 `json:"categoryName"`
	InstanceId   string                 `json:"instanceId"`
	ResourceType string                 `json:"resourceType"`
	ResourceAttr map[string]interface{} `json:"resourceAttr"`
}

// 创建资源命令拼接
func buildCommands(req ResourceModelRequest) []string {
	cmds := make([]string, 0, len(req.Data)+1)
	resourcesName := make([]string, 0, len(req.Data))
	for _, singleRes := range req.Data {
		resourcesName = append(resourcesName, singleRes.InstanceId)
		cmd := utils.CmdResourceCreate + singleRes.InstanceId + " " + singleRes.ResourceType
		// 遍历 ResourceAttr，拼接所有属性和属性值
		for key, value := range singleRes.ResourceAttr {
			cmd += " " + key + "=" + fmt.Sprintf("%v", value)
		}
		cmds = append(cmds, cmd)
	}
	cmd := "pcs resource group add " + req.GroupName
	for _, resourceName := range resourcesName {
		cmd += resourceName + " "
	}
	cmds = append(cmds, cmd)
	return cmds
}

func ResourceCreate(req ResourceModelRequest) ResourceModelPostResponse {
	var res ResourceModelPostResponse
	cmds := buildCommands(req)

	if !req.IsTrue {
		res.Data = cmds
		return res
	}

	CmdTestCIB := "pcs cluster cib test.xml"
	output, err := utils.RunCommand(CmdTestCIB)
	if err != nil {
		res.Action = false
		res.Error = gettext.Gettext("pcs cluster cib command failed: " + string(output))
		return res
	}

	results := make([]string, 0, len(cmds)+1)
	for _, cmd := range cmds {
		out, err := utils.RunCommand(cmd)
		if err != nil {
			results = append(results, string(out))
		}
	}
	if len(results) == 0 {
		res.Action = true
		res.Info = gettext.Gettext("resources create success")
		return res
	}
	res.Action = false
	res.Error = results
	return res
}
