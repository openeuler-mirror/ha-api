/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liqiuyu <liqiuyu@kylinos.cn>
 * Date: Tue Jan 12 09:51:22 2021 +0800
 */
package models

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/chai2010/gettext-go"
)

func GenerateLog() map[string]interface{} {
	result := map[string]interface{}{}
	file := map[string]string{}

	var out []byte
	var err error
	if out, err = utils.RunCommand(utils.CmdGenLog); err != nil {
		fmt.Println(err)
		result["action"] = false
		result["error"] = gettext.Gettext("Get neokylinha log failed")
		return result
	}

	filePath := strings.Split(string(out), "\n")[0]
	file["filepath"] = filePath
	result["action"] = true
	result["data"] = file
	return result
}
