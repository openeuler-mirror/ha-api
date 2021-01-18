package routers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/controllers"
)

func init() {
	web.Router("/test", &controllers.TestController{})

	web.Router("/", &controllers.RootController{})
	web.Router("/resource", &controllers.IndexController{})
	// the same logic as /resource
	web.Router("/login", &controllers.IndexController{})

	// web.Router("/api/haclusters/1/nodes", &controllers.HAClustersController{})
	// web.Router("/api/haclusters/1/nodes/:nodeid", &controllers.HAClustersController{})

	ns := web.NewNamespace("/api/v1",
		web.NSRouter("/haclusters/1", &controllers.HAClustersController{}),
		web.NSRouter("/login", &controllers.LoginController{}),
		web.NSRouter("/logout", &controllers.LogoutController{}),

		web.NSRouter("/haclusters/1/resources", &controllers.ResourceController{}),
		web.NSRouter("/haclusters/1/resources/:rscID/:action", &controllers.ResourceController{}),

		web.NSRouter("/haclusters/1/configs", &controllers.HeartBeatController{}),

		web.NSRouter("/haclusters/1/nodes", &controllers.NodesController{}),
		web.NSRouter("/haclusters/1/nodes/:nodeID", &controllers.NodeController{}),
		web.NSRouter("/haclusters/1/nodes/:nodeID/:action", &controllers.NodeActionController{}),
		web.NSRouter("/haclusters/1/localnodes/:action", &controllers.LocalHaOperation{}),
		web.NSRouter("/haclusters/1/logs", &controllers.LogController{}),
		web.NSRouter("/haclusters/1/alarms", &controllers.AlarmConfig{}),
	)
	web.AddNamespace(ns)
}
