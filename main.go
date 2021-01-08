package main

import (
	_ "openkylin.com/ha-api/routers"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func main() {

	logs.SetLogger("console")
	logs.SetLevel(logs.LevelDebug)

	web.Run()
}
