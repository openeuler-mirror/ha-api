package main

import (
	_ "openkylin.com/ha-api/routers"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func main() {

	logs.SetLogger("console")
	logs.SetLevel(logs.LevelDebug)

	web.BConfig.CopyRequestBody = true
	web.BConfig.Listen.HTTPAddr = "172.30.30.94"
	web.SetStaticPath("/static", "views/static")

	// web.SetStaticPath("/4.12.13", "views/static/4.12.13")
	// web.SetStaticPath("/static", "views/static/static")

	web.Run()
}
