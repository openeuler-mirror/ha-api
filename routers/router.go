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
 * Date: 2024-03-19 11:19:33
 * LastEditTime: 2024-03-19 15:00:25
 * Description: 网页路由
 ******************************************************************************/
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
		web.NSRouter("/haclusters/1", &controllers.ClustersController{}),
		web.NSRouter("/haclusters/1/cluster_status", &controllers.ClustersStatusController{}),
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
		web.NSRouter("/managec/is_cluster_exist", &controllers.IsClusterExist{}),

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

		web.NSRouter("/haclusters/1/utilization", &controllers.UtilizationController{}),
		web.NSRouter("/haclusters/1/tag", &controllers.TagController{}),
		web.NSRouter("/haclusters/1/rules", &controllers.RuleController{}),
		web.NSRouter("/haclusters/1/scripts", &controllers.ScriptsController{}),
		web.NSRouter("/remotescripts", &controllers.ScriptsRemoteController{}),
	)
	nr := web.NewNamespace("/remote/api/v1",
		web.NSRouter("/sync_config", &controllers.Sync_configController{}),
		web.NSRouter("/nodes/add_nodes", &controllers.LocalAddNodesController{}),
		web.NSRouter("/managec/local_cluster_info", &controllers.LocalClusterInfoController{}),
		web.NSRouter("/destroy_cluster", &controllers.LocalClusterDestroyController{}),
	)
	web.AddNamespace(ns)
	web.AddNamespace(nr)
}
