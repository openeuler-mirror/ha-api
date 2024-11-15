/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-19 17:37:51
 * Description: 日志下载控制器
 ******************************************************************************/
package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"gitee.com/openeuler/ha-api/models"
)

type LogController struct {
	web.Controller
}

func (lc *LogController) Get() {
	logs.Debug("handle resource GET request")
	result := models.GenerateLog()
	lc.Data["json"] = &result
	lc.ServeJSON()
}

func (lc *LogController) Put() {
	logs.Debug("handle resource PUT request")

	lc.ServeJSON()
}

func (lc *LogController) Post() {
	logs.Debug("handle resource POST request")

	lc.ServeJSON()
}

type LogDownloadController struct {
	web.Controller
}

func (ldc *LogDownloadController) Get() {
	logs.Debug("handle log download GET request")
	fileTail := ldc.Ctx.Input.Param(":filetail")
	// result := models.GenerateLog()
	// lc.Data["json"] = &result
	// lc.ServeJSON()

	const filePath = "/usr/share/heartbeat-gui/ha-api/static/"
	filePrefix := "kylinha-log-"
	http.ServeFile(ldc.Ctx.ResponseWriter, ldc.Ctx.Request, filePath+filePrefix+fileTail)
}
