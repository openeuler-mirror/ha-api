/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Tue Jan 5 09:33:00 2021 +0800
 */
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
