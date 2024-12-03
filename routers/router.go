/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Tue Jan 5 09:33:00 2021 +0800
 */
package routers

import (
	"gitee.com/openeuler/ha-api/controllers"
	"github.com/beego/beego/v2/server/web"
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
		web.NSRouter("/:cluster_name/1", &controllers.ClustersController{}),
		web.NSRouter("/:cluster_name/1/cluster_status", &controllers.ClustersStatusController{}),
		web.NSRouter("/login", &controllers.LoginController{}),
		web.NSRouter("/logout", &controllers.LogoutController{}),
		web.NSRouter("/user", &controllers.PasswordChangeController{}),
		web.NSRouter("/managec/cluster_add", &controllers.MultipleClustersController{}),
		//web.NSRouter("/managec/sync_config", &controllers.Sync_configController{}),
		web.NSRouter("/managec/cluster_setup", &controllers.ClusterSetupController{}),
		web.NSRouter("/managec/cluster_destroy", &controllers.ClusterDestroyController{}),
		web.NSRouter("/managec/cluster_remove", &controllers.ClusterRemoveController{}),
		web.NSRouter("/managec/add_nodes", &controllers.AddNodesController{}),
		web.NSRouter("/managec/local_cluster_info", &controllers.LocalClusterInfoController{}),
		web.NSRouter("/managec/is_cluster_exist", &controllers.IsClusterExistController{}),

		web.NSRouter("/:cluster_name/1/resources", &controllers.ResourceController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID/:action", &controllers.ResourceActionController{}),
		web.NSRouter("/:cluster_name/1/resources/meta_attributes/:category", &controllers.ResourceMetaAttributesController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID/relations/:relation", &controllers.ResourceRelationsController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID", &controllers.ResourceOpsById{}),

		web.NSRouter("/:cluster_name/1/metas", &controllers.MetasController{}),
		web.NSRouter("/:cluster_name/1/metas/:rsc_class/:rsc_type", &controllers.MetaController{}),
		web.NSRouter("/:cluster_name/1/metas/:rsc_class/:rsc_type/:rsc_provider", &controllers.MetaController{}),

		web.NSRouter("/:cluster_name/1/hbstatus", &controllers.HeartBeatStatusController{}),
		web.NSRouter("/:cluster_name/1/get_diskhb_device", &controllers.DiskHeartBeatController{}),

		web.NSRouter("/:cluster_name/1/configs", &controllers.HeartBeatController{}),

		web.NSRouter("/:cluster_name/1/nodes", &controllers.NodesController{}),
		web.NSRouter("/:cluster_name/1/nodes/:nodeID", &controllers.NodeController{}),
		web.NSRouter("/:cluster_name/1/nodes/:nodeID/:action", &controllers.NodeActionController{}),
		web.NSRouter("/:cluster_name/1/localnodes/:action", &controllers.LocalHaOperation{}),

		web.NSRouter("/:cluster_name/1/logs", &controllers.LogController{}),

		web.NSRouter("/:cluster_name/1/alarms", &controllers.AlarmConfig{}),

		web.NSRouter("/:cluster_name/1/commands", &controllers.CommandsController{}),
		web.NSRouter("/:cluster_name/1/commands/:cmd_type", &controllers.CommandsRunnerController{}),

		web.NSRouter("/:cluster_name/1/utilization", &controllers.UtilizationController{}),
		web.NSRouter("/:cluster_name/1/tag", &controllers.TagController{}),
		web.NSRouter("/:cluster_name/1/tag/:tag_name", &controllers.TagUpdateController{}),
		web.NSRouter("/:cluster_name/1/tag/:tag_name/:action", &controllers.TagActionController{}),
		web.NSRouter("/:cluster_name/1/rules", &controllers.RuleController{}),
		web.NSRouter("/:cluster_name/1/scripts", &controllers.ScriptsController{}),
		web.NSRouter("/remotescripts", &controllers.ScriptsRemoteController{}),
	)

	nr := web.NewNamespace("/remote/api/v1",
		web.NSRouter("/sync_config", &controllers.Sync_configController{}),
		web.NSRouter("/nodes/add_nodes", &controllers.LocalAddNodesController{}),
		web.NSRouter("/managec/local_cluster_info", &controllers.LocalClusterInfoController{}),
		web.NSRouter("/destroy_cluster", &controllers.LocalClusterDestroyController{}),

		web.NSRouter("/:cluster_name/1", &controllers.ClustersController{}),
		web.NSRouter("/:cluster_name/1/cluster_status", &controllers.ClustersStatusController{}),

		web.NSRouter("/:cluster_name/1/resources", &controllers.ResourceController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID/:action", &controllers.ResourceActionController{}),
		web.NSRouter("/:cluster_name/1/resources/meta_attributes/:category", &controllers.ResourceMetaAttributesController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID/relations/:relation", &controllers.ResourceRelationsController{}),
		web.NSRouter("/:cluster_name/1/resources/:rscID", &controllers.ResourceOpsById{}),

		web.NSRouter("/:cluster_name/1/metas", &controllers.MetasController{}),
		web.NSRouter("/:cluster_name/1/metas/:rsc_class/:rsc_type", &controllers.MetaController{}),
		web.NSRouter("/:cluster_name/1/metas/:rsc_class/:rsc_type/:rsc_provider", &controllers.MetaController{}),

		// desperated
		// web.NSRouter("/:cluster_name/1/hbstatus", &controllers.HeartBeatStatusController{}),
		// web.NSRouter("/:cluster_name/1/get_diskhb_device", &controllers.DiskHeartBeatController{}),

		web.NSRouter("/:cluster_name/1/configs", &controllers.HeartBeatController{}),

		web.NSRouter("/:cluster_name/1/nodes", &controllers.NodesController{}),
		web.NSRouter("/:cluster_name/1/nodes/:nodeID", &controllers.NodeController{}),
		web.NSRouter("/:cluster_name/1/nodes/:nodeID/:action", &controllers.NodeActionController{}),
		web.NSRouter("/:cluster_name/1/localnodes/:action", &controllers.LocalHaOperation{}),

		web.NSRouter("/:cluster_name/1/logs", &controllers.LogController{}),

		web.NSRouter("/:cluster_name/1/alarms", &controllers.AlarmConfig{}),

		web.NSRouter("/:cluster_name/1/commands", &controllers.CommandsController{}),
		web.NSRouter("/:cluster_name/1/commands/:cmd_type", &controllers.CommandsRunnerController{}),

		web.NSRouter("/:cluster_name/1/utilization", &controllers.UtilizationController{}),
		web.NSRouter("/:cluster_name/1/tag", &controllers.TagController{}),
		web.NSRouter("/:cluster_name/1/tag/:tag_name", &controllers.TagUpdateController{}),
		web.NSRouter("/:cluster_name/1/tag/:tag_name/:action", &controllers.TagActionController{}),
		web.NSRouter("/:cluster_name/1/rules", &controllers.RuleController{}),
		web.NSRouter("/:cluster_name/1/scripts", &controllers.ScriptsController{}),
	)
	web.AddNamespace(ns)
	web.AddNamespace(nr)
}
