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
 * LastEditTime: 2022-04-19 17:37:53
 * Description: 网页控制器
 ******************************************************************************/
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
		// default find templete from views folder
		rc.TplName = "static/index.html"
		return
	}
	rc.Redirect("/login", http.StatusOK)
}

type IndexController struct {
	web.Controller
}

func (ic *IndexController) Get() {
	// default find templete from views folder
	ic.TplName = "static/index.html"
}
