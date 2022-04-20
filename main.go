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
 * LastEditTime: 2022-04-20 14:21:08
 * Description: ha-api的入口
 ******************************************************************************/
package main

import (
	"net/http"

	_ "openkylin.com/ha-api/routers"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func pageNotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("page not found"))
}

func main() {

	logs.SetLogger("console")
	logs.SetLevel(logs.LevelNotice)

	web.BConfig.CopyRequestBody = true
	//	web.BConfig.Listen.HTTPAddr = "172.30.30.94"
	web.SetStaticPath("/static", "views/static")

	// web.SetStaticPath("/4.12.13", "views/static/4.12.13")
	// web.SetStaticPath("/static", "views/static/static")

	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 300

	web.ErrorHandler("404", pageNotFoundHandler)

	web.Run()
}
