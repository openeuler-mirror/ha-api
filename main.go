package main

import (
	_ "openkylin.com/ha-api/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.Run()
}
