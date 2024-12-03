/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: 江新宇 <jiangxinyu@kylinos.cn>
 * Date: Tue Jan 19 22:19:26 2021 +0800
 */
package models

import (
	"errors"

	"gitee.com/openeuler/ha-api/utils"
)

// DO NOT EDIT ANYWHERE!!!
var commandTable = map[int]string{
	1: "crm_mon -1 -o",
	2: "crm_simulate -Ls",
	3: "pcs config show",
	4: "corosync-cfgtool -s",
}

func GetCommandsList() map[string]interface{} {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = commandTable
	return result
}

func RunBuiltinCommand(cmdID int) (string, error) {
	if _, ok := commandTable[cmdID]; ok {
		cmd := commandTable[cmdID]
		output, err := utils.RunCommand(cmd)
		if err != nil {
			return "", errors.New("Run command failed: " + err.Error())
		}
		return string(output), nil
	}

	return "", errors.New("invalid command index")
}
