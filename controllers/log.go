/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package controllers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
)

type LogController struct {
	web.Controller
}

func (lc *LogController) Get() {
	slog.Debug("handle resource GET request")

	clusterName := lc.Ctx.Input.Param(":cluster_name")

	if !utils.IsLocalCluster(clusterName) {
		currentPath := lc.Ctx.Input.URL()
		localConf := models.GetLocalConf()
		remoteNodes := localConf.GetNodes(clusterName)
		if len(remoteNodes) == 0 {
			res := map[string]interface{}{}
			res["action"] = false
			res["error"] = "Can not get remote nodes"
			lc.Data["json"] = &res
			lc.ServeJSON()
			return
		}
		for _, node := range remoteNodes {
			url := utils.GenerateRemoteRequestURL(node, currentPath)
			res, err := utils.SendRequest(url, lc.Ctx.Request.Method, lc.Ctx.Input.RequestBody)
			if err != nil {
				slog.Warn("request to node failed (retry next node)", "url", url, "node", node, "error", err)
				continue
			}
			defer res.Body.Close()
			respData, err := io.ReadAll(res.Body)
			if err != nil {
				slog.Error("failed to read response body", "url", currentPath, "node", node, "error", err)
				continue
			}
			contentDispostion := res.Header.Get("content-disposition")
			parts := strings.Split(contentDispostion, "filename=")
			if len(parts) < 2 {
				slog.Error("content-disposition header missing filename", "header", contentDispostion)
				continue
			}
			fileName := strings.TrimPrefix(parts[1], "\"")
			fileName = strings.TrimSuffix(fileName, "\"")

			lc.Ctx.Output.Header("content-type", "application/octet-stream")
			lc.Ctx.Output.Header("content-transfer-encoding", "binary")
			lc.Ctx.Output.Header("content-disposition", "attachment;filename="+fileName)

			lc.Ctx.Output.Body(respData)
			return

		}
		slog.Error("UrlRedirect failed (all nodes )", "cluster", clusterName, "url", currentPath)
		res := map[string]interface{}{}
		res["action"] = false
		res["error"] = "no nodes succeeded"
		lc.Data["json"] = &res
		lc.ServeJSON()
		return
	}

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
