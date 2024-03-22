/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bixiaoyan
 * Date: 2024-03-13 15:05:51
 * LastEditTime: 2024-03-22 17:24:41
 * Description: tag
 ******************************************************************************/
package models

import (
	"encoding/xml"
    "strings"
	"gitee.com/openeuler/ha-api/utils"
    "github.com/beego/beego/v2/core/logs"
    
)

type TagInfo struct {
    ID       string      `json:"id"`
    TagsRsc  interface{} `json:"tagsrsc"`
}

type Tag struct {
    ID       string    `xml:"id,attr"`
    ObjRefs  []ObjRef  `xml:"obj_ref"`
}
type ObjRef struct {
    ID      string   `xml:"id,attr"`
}
type Tags struct {
    Tags    []Tag    `xml:"tag"`
}

type TagGetResult struct {
    Action  bool        `json:"action"`
    Data    []TagInfo   `json:"data"`
    ResList []string    `json:"res_list"`
    Error   string      `json:"error,omitempty"`
}



func GetTag() TagGetResult {
    resList := []string{}
    cmd := "cibadmin --query --scope tags"
    rsc_info := GetResourceInfo()
    
    out, err := utils.RunCommand(cmd)
    for _, res := range rsc_info["data"].([]map[string]interface{}) {        
        resList = append(resList, res["id"].(string))
    }

    if err != nil{
        if strings.Contains(string(out), "No such device or address") {
            return TagGetResult{Action: true, Data: []TagInfo{}, ResList: resList}
        }else {
            return TagGetResult{Action: false, Error: "tag信息获取失败"}
        }
        
    }
    
    // 解析tag数据
    // confJSON := map[string]interface{}{}
    var tags Tags
    if err := xml.Unmarshal([]byte(out), &tags); err != nil {
        return TagGetResult{Action: false, Error: "JSON解析失败"}
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
