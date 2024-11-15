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
 * LastEditTime: 2024-04-22 09:23:45
 * Description: ha-api的入口
 ******************************************************************************/
package main

import (
	"net/http"

	_ "gitee.com/openeuler/ha-api/routers"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

	"gitee.com/openeuler/ha-api/utils"
)

func pageNotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("page not found"))
}

func detectUserLanguage() string {
	// Todo: check user language preferences, check environment variables, user profiles, or HTTP request headers
	return "zh_CN"
}

func init() {
	gettext.BindLocale(gettext.New("ha-api", "locale"))
	switch language := detectUserLanguage(); language {
	case "zh_CN":
		gettext.SetLanguage("zh_CN")
	default:
		gettext.SetLanguage("zh_CN") // default chinese
	}
}

func main() {

	logs.SetLogger("console")
	logs.SetLevel(logs.LevelNotice)

	web.BConfig.CopyRequestBody = true

	web.SetStaticPath("/static", "views/static")

	// web.SetStaticPath("/4.12.13", "views/static/4.12.13")
	// web.SetStaticPath("/static", "views/static/static")

	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 300

	web.ErrorHandler("404", pageNotFoundHandler)

	port, _ := utils.ReadPortFromConfig()
	web.Run(":" + port)
}
