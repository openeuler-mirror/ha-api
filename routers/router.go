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

	// http://172.30.30.94:8080/kylinha-log-ha2-20210130175130.tar
	web.Router("/kylinha-log-:filetail(.*\\.tar$)", &controllers.LogDownloadController{})

	ns := web.NewNamespace("/api/v1",
		web.NSRouter("/haclusters/1", &controllers.HAClustersController{}),
		web.NSRouter("/login", &controllers.LoginController{}),
		web.NSRouter("/logout", &controllers.LogoutController{}),

		web.NSRouter("/haclusters/1/resources", &controllers.ResourceController{}),
		web.NSRouter("/haclusters/1/resources/:rscID/:action", &controllers.ResourceActionController{}),
		web.NSRouter("/haclusters/1/resources/meta_attributes/:catagory", &controllers.ResourceMetaAttributesController{}),
		web.NSRouter("/haclusters/1/resources/:rscID/relations/:relation", &controllers.ResourceRelationsController{}),
		web.NSRouter("/haclusters/1/resources/:rscID", &controllers.ResourceOpsById{}),

		web.NSRouter("/haclusters/1/metas", &controllers.MetasController{}),
		web.NSRouter("/haclusters/1/metas/:rsc_class/:rsc_type", &controllers.MetaController{}),
		web.NSRouter("/haclusters/1/metas/:rsc_class/:rsc_type/:rsc_provider", &controllers.MetaController{}),

		web.NSRouter("/haclusters/1/hbstatus", &controllers.HeartBeatStatusController{}),
		web.NSRouter("/haclusters/1/get_diskhb_device", &controllers.DiskHeartBeatController{}),

		web.NSRouter("/haclusters/1/configs", &controllers.HeartBeatController{}),

		web.NSRouter("/haclusters/1/nodes", &controllers.NodesController{}),
		web.NSRouter("/haclusters/1/nodes/:nodeID", &controllers.NodeController{}),
		web.NSRouter("/haclusters/1/nodes/:nodeID/:action", &controllers.NodeActionController{}),
		web.NSRouter("/haclusters/1/localnodes/:action", &controllers.LocalHaOperation{}),

		web.NSRouter("/haclusters/1/logs", &controllers.LogController{}),

		web.NSRouter("/haclusters/1/alarms", &controllers.AlarmConfig{}),

		web.NSRouter("/haclusters/1/commands", &controllers.CommandsController{}),
		web.NSRouter("/haclusters/1/commands/:cmd_type", &controllers.CommandsRunnerController{}),
	)
	web.AddNamespace(ns)
}
