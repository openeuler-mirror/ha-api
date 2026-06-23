/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package filters

import (
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
)

const (
	apiPrefix     = "/api"
	versionPrefix = "/v1"
)

var LoginFilter = func(ctx *context.Context) {
	// implement login check
	username := ctx.Input.Session("username")
	if username == nil {
		if !strings.Contains(ctx.Request.RequestURI, "/login") && !strings.Contains(ctx.Request.RequestURI, "/remote") && !strings.Contains(ctx.Request.RequestURI, "/public-key") {
			if strings.HasPrefix(ctx.Request.RequestURI, apiPrefix) {
				ctx.Redirect(403, "session timeout")
			}
		}
	}
}

var RemoteClusterRedirectFilter = func(ctx *context.Context) {
	// 1. 检查路径是否包含 :cluster_name（仅处理需要远程跳转的路由）
	clusterName := ctx.Input.Param(":cluster_name")
	if clusterName == "" {
		return
	}
	nodeName := ctx.Input.Param(":node_name")
	if nodeName != "" {
		return
	}

	if utils.IsLocalCluster(clusterName) {
		return
	}

	currentPath := ctx.Request.URL.String() // 修改携带请求参数
	slog.Debug("Redirect URL", "url", currentPath, "clusterName", clusterName)

	if strings.Contains(currentPath, "logs") {
		return
	}

	var result map[string]interface{}
	var err error
	result, err = models.UrlRedirect(
		clusterName,
		currentPath,
		ctx.Request.Method,
		ctx.Input.RequestBody,
		nil,
	)
	if err != nil {
		ctx.Output.JSON(result, false, false)
		return
	}
	ctx.Output.JSON(result, false, false)
}

// 全局过滤器函数
func LogRequestControllerFilter(ctx *context.Context) {
	// 获取控制器名和方法名（如：RuleController.Get）
	startTime := time.Now()
	ctx.Input.SetData("startTime", startTime)
	url := ctx.Input.URL()
	method := ctx.Input.Method()
	host := ctx.Input.Host()
	// ip := ctx.Input.IP()
	if !strings.Contains(ctx.Request.RequestURI, "/remote") {
		slog.Info("handle request",
			"path", url,
			"method", method,
			"host", host,
		)
	} else {
		slog.Debug("handle request",
			"path", url,
			"method", method,
			"host", host,
		)

	}
}

func LogResponseControllerFilter(ctx *context.Context) {
	// 获取响应状态码
	statusCode := ctx.Output.Status
	duration := time.Since(ctx.Input.GetData("startTime").(time.Time))

	// 基础日志字段
	logFields := []interface{}{
		"status", statusCode,
		"path", ctx.Request.URL.Path,
		"method", ctx.Request.Method,
		"duration", duration.String(),
	}

	// 500错误特殊处理
	if statusCode >= 500 {
		// 获取堆栈跟踪
		stack := debug.Stack()

		// 添加额外字段
		logFields = append(logFields,
			"stack", string(stack),
		)
		slog.Error("server error", logFields...)
	} else {
		if !strings.Contains(ctx.Request.RequestURI, "/remote") {
			slog.Info("request completed", logFields...)
		} else {
			slog.Debug("request completed", logFields...)
		}

	}

}
