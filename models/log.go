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
 * LastEditTime: 2022-04-20 10:35:08
 * Description: 日志功能
 ******************************************************************************/
package models

import (
	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

func GenerateLog() map[string]interface{} {
	result := map[string]interface{}{}
	file := map[string]string{}

	var out []byte
	var err error
	if out, err = utils.RunCommand(utils.CmdGenLog); err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("Get neokylinha log failed")
		return result
	}

	filePath := string(out)
	file["filepath"] = filePath
	result["action"] = true
	result["data"] = file
	return result
}
