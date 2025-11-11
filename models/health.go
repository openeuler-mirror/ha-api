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