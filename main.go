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
	logs.SetLevel(logs.LevelDebug)

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
