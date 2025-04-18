/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Fri Mar 22 17:26:24 2024 +0800
 */
package models

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

type TagInfo struct {
	ID      string      `json:"id"`
	TagsRsc interface{} `json:"tagsrsc"`
}

type Tag struct {
	ID      string   `xml:"id,attr"`
	ObjRefs []ObjRef `xml:"obj_ref"`
}
type ObjRef struct {
	ID string `xml:"id,attr"`
}
type Tags struct {
	Tags []Tag `xml:"tag"`
}

type TagPostData struct {
	ID           string   `json:"id"`
	Tag_resource []string `json:"tag_resource"`
}

type TagGetResult struct {
	Action  bool      `json:"action"`
	Data    []TagInfo `json:"data"`
	ResList []string  `json:"res_list"`
	Error   string    `json:"error,omitempty"`
}

func GetTag() TagGetResult {
	resList := []string{}
	cmd := "cibadmin --query --scope tags"
	rsc_info := GetResourceInfo()

	out, err := utils.RunCommand(cmd)
	for _, res := range rsc_info["data"].([]map[string]interface{}) {
		resList = append(resList, res["id"].(string))
	}

	if err != nil {
		if strings.Contains(string(out), "No such device or address") {
			return TagGetResult{Action: true, Data: []TagInfo{}, ResList: resList}
		} else {
			return TagGetResult{Action: false, Error: gettext.Gettext("Get tag failed")}
		}

	}

	// 解析tag数据
	var tags Tags
	if err := xml.Unmarshal([]byte(out), &tags); err != nil {
		return TagGetResult{Action: false, Error: gettext.Gettext("Parsing JSON failed")}
	}

	tagInfos := make([]TagInfo, len(tags.Tags))
	for i, tag := range tags.Tags {
		tagInfos[i] = TagInfo{
			ID:      tag.ID,
			TagsRsc: tag.ObjRefs,
		}
	}

	return TagGetResult{Action: true, Data: tagInfos, ResList: resList}
}

func SetTag(data []byte) utils.GeneralResponse {

	var result utils.GeneralResponse
	// json数据解析
	if data == nil || len(data) == 0 {
		result.Action = false
		result.Error = gettext.Gettext("No input data")
		return result
	}
	var TagData TagPostData
	err := json.Unmarshal(data, &TagData)
	if err != nil {
		result.Action = false
		result.Error = gettext.Gettext("Cannot convert data to json map")
		return result
	}

	cmd := "pcs tag create " + TagData.ID
	for _, res := range TagData.Tag_resource {
		cmd = cmd + " " + string(res)
	}
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = gettext.Gettext("Add tag success")
	} else {
		result.Action = false
		result.Error = string(out)
	}
	return result
}

func UpdateTag(tagName string, data []byte) utils.GeneralResponse {
	var result utils.GeneralResponse
	// json数据解析
	if data == nil || len(data) == 0 {
		result.Action = false
		result.Error = gettext.Gettext("No input data")
		return result
	}
	var TagData TagPostData
	err := json.Unmarshal(data, &TagData)
	if err != nil {
		result.Action = false
		result.Error = gettext.Gettext("Cannot convert data to json map")
		return result
	}
	cmd_create := "pcs tag create " + string(tagName)
	cmd_delete := "pcs tag delete " + string(tagName)
	for _, res := range TagData.Tag_resource {
		cmd_create = cmd_create + " " + string(res)
	}
	out1, err1 := utils.RunCommand(cmd_delete)
	if err1 == nil {
		out2, err2 := utils.RunCommand(cmd_create)
		if err2 == nil {
			result.Action = true
			result.Info = gettext.Gettext("Update tag success")
		} else {
			result.Action = false
			result.Error = string(out2)
		}

	} else {
		result.Action = false
		result.Error = string(out1)
	}
	return result
}

func TagAction(tagName string, action string) utils.GeneralResponse {
	var result utils.GeneralResponse
	cmd := "pcs resource "
	if action == "delete" {
		return DeleteTag(tagName)
	} else if action == "start" {
		cmd = cmd + "enable " + tagName 
	} else if action == "stop" {
		cmd = cmd + "disable " + tagName
	}
	print(cmd)
	out, err := utils.RunCommand(cmd)
	if err == nil {
			result.Action = true
			result.Info = gettext.Gettext("Tag action success")
	} else {
		result.Action = false
		result.Error = string(out)
	}
	return result
}

func DeleteTag(tagName string) utils.GeneralResponse {
	var result utils.GeneralResponse
	cmd := "pcs tag delete " + tagName 
	out, err := utils.RunCommand(cmd)
	if err == nil {
		result.Action = true
		result.Info = gettext.Gettext("Delete tag success")
	} else {
		result.Action = false
		result.Error = string(out)
	}
	return result
}




