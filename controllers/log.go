/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liqiuyu <liqiuyu@kylinos.cn>
 * Date: Tue Jan 12 09:51:22 2021 +0800
 */
package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/settings"
)

type LogController struct {
	web.Controller
}

func (lc *LogController) Get() {
	slog.Debug("handle resource GET request")
	// todo: multi clusters
	result, geterr := models.GenerateLog()

	if geterr != nil {
		res := map[string]interface{}{}
		res["action"] = false
		res["error"] = geterr.Error()
		lc.Data["json"] = &res
		lc.ServeJSON()
	} else {

		FileInfo, err := os.Stat(result)
		if err != nil {
			res := map[string]interface{}{}
			res["action"] = false
			res["error"] = gettext.Gettext(fmt.Sprintf("Can not stat file %s", result))
			slog.Error("Can not stat file", "file", result, "error", err)
			lc.Data["json"] = &res
			lc.ServeJSON()
			return
		}
		slog.Info("Generated log file", "file", result)
		defer os.Remove(result)
		lc.Ctx.Output.Download(result, FileInfo.Name())
	}

}

func (lc *LogController) Put() {
	slog.Debug("handle resource PUT request")

	lc.ServeJSON()
}

func (lc *LogController) Post() {
	slog.Debug("handle resource POST request")

	lc.ServeJSON()
}

type LogDownloadController struct {
	web.Controller
}

func (ldc *LogDownloadController) Get() {
	slog.Debug("handle log download GET request")
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

	filePrefix := "kylinha-log-"
	fullPath := filepath.Join(settings.StaticDir, filePrefix+fileTail)
	// Verify the resolved path is still within the expected directory
	if !strings.HasPrefix(fullPath, settings.StaticDir) {
		ldc.Ctx.Output.SetStatus(http.StatusBadRequest)
		ldc.Data["json"] = map[string]interface{}{"action": false, "error": "invalid filename"}
		ldc.ServeJSON()
		return
	}
	http.ServeFile(ldc.Ctx.ResponseWriter, ldc.Ctx.Request, fullPath)
}
