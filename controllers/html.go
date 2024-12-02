/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type RootController struct {
	web.Controller
}

func (rc *RootController) Get() {
	logs.Debug("handle index request")

	// TODO: check session
	if true {
		// default find template from views folder
		rc.TplName = "static/index.html"
		return
	}
	rc.Redirect("/login", http.StatusOK)
}

type IndexController struct {
	web.Controller
}

func (ic *IndexController) Get() {
	// default find template from views folder
	ic.TplName = "static/index.html"
}
