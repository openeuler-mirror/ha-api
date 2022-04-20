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
 * LastEditTime: 2022-04-20 11:18:25
 * Description: 运行日志命令获取结果
 ******************************************************************************/
package utils

import (
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	logs.Debug("Running command: %s", c)
	command := exec.Command("bash", "-c", c)
	out, err := command.CombinedOutput()
	if err != nil {
		logs.Error("Run command failed!, command: " + c + " out: " + string(out) + " err: " + err.Error())
	}
	return out, err
}
