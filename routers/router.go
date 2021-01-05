package routers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/controllers"
)

func init() {
	web.Router("/api/v1/haclusters/1", &controllers.HAClustersController{})
}
