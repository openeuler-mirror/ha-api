/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: liqiuyu <liqiuyu@kylinos.cn>
 * Date: Tue Jan 12 09:51:22 2021 +0800
 */
package controllers

import (
	"net/http"
	"path/filepath"
        "strings"

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
        // Prevent path traversal: extract only the base filename
        fileTail = filepath.Base(fileTail)
        // Validate file extension
        if !strings.HasSuffix(fileTail, ".tar") {
                ldc.Ctx.Output.SetStatus(http.StatusBadRequest)
                ldc.Data["json"] = map[string]interface{}{"action": false, "error": "invalid filename"}
                ldc.ServeJSON()
                return
        }

        const filePath = "/usr/share/heartbeat-gui/ha-api/static/"
        filePrefix := "kylinha-log-"
        fullPath := filepath.Join(filePath, filePrefix+fileTail)
        // Verify the resolved path is still within the expected directory
        if !strings.HasPrefix(fullPath, filePath) {
                ldc.Ctx.Output.SetStatus(http.StatusBadRequest)
                ldc.Data["json"] = map[string]interface{}{"action": false, "error": "invalid filename"}
                ldc.ServeJSON()
                return
        }
        http.ServeFile(ldc.Ctx.ResponseWriter, ldc.Ctx.Request, fullPath)
}
