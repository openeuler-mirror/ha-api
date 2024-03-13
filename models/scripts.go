/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-13 14:16:22
 * LastEditTime: 2024-03-13 14:57:14
 * Description:脚本生成器模块相关接口
 */
package models

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
)

const (
	pacemakerAgentsCmd = "pcs resource agents ocf:pacemaker"
)

func IsScriptExist(scriptName string) utils.GeneralResponse {
	out, err := utils.RunCommand(pacemakerAgentsCmd)
	if err != nil {
		return utils.HandleCmdError("查询脚本命令执行失败", false)
	}
	scripts := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, script := range scripts {
		if script == scriptName {
			logs.Warn(fmt.Sprintf("脚本 %s 已存在于pacemaker目录下", scriptName))
			return utils.GeneralResponse{
				Action: false,
				Error:  "脚本已经存在于pacemaker目录下",
			}
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   "脚本不存在",
	}
}
