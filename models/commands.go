/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 10:30:19
 * Description: 集群命令
 ******************************************************************************/
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

	return "", errors.New("Invalid command index")
}
