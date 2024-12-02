/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Tue Jan 5 09:33:00 2021 +0800
 */
package utils

import (
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
	"github.com/pkg/errors"
)

const (
	CmdNodeStatus                   = "crm_node -l"
	CmdClusterStatus                = "crm_mon -1"
	CmdClusterStatusAsXML           = CmdClusterStatus + " --as-xml"
	CmdHostName                     = "hostname"
	CmdCountClustersConfigsBackuped = "ls /usr/share/heartbeat-gui/ha-api/ClustersInfo.conf.* | wc -l"
	CmdCibQueryConfig               = "cibadmin --query --scope configuration"
	CmdDeleteAlert                  = "pcs alert delete alert_Kylin"
	CmdCreateAlert                  = "pcs alert create id=alert_log path=/usr/share/pacemaker/alerts/alert_log.sh"
	CmdAddAlert                     = "pcs alert recipient add alert_Kylin"
	CmdUpdateResourceStickness      = "crm_attribute  -t rsc_defaults -n resource-stickiness -v "
	CmdUpdateCrmConfig              = "crm_attribute -t crm_config -n "
	CmdQueryCIB                     = "cibadmin -Q"
	CmdQueryCrmConfig               = "cibadmin --query --scope crm_config"
	CmdQueryResources               = "cibadmin --query --scope resources"
	CmdQueryConstraints             = "cibadmin --query --scope constraints"
	CmdQueryResourcesById           = "cibadmin --query --scope resources|grep 'id=\"%s\"'"

	CmdStartCluster             = "pcs cluster start "
	CmdStopClusterLocal         = "pcs cluster stop "
	CmdStopCluster              = "pcs cluster stop --all"
	CmdSetupCluster             = "pcs cluster setup hacluster"
	CmdSetupClusterStandard     = "pcs cluster setup %s %s totem token=8000 --start"
	CmdDestroyCluster           = "pcs cluster destroy --all"
	CmdNodeAdd                  = "pcs cluster node add %s "
	CmdNodeAddStart             = CmdNodeAdd + "--start"
	CmdNodeStandby              = "pcs node standby "
	CmdNodeUnStandby            = "pcs node unstandby "
	CmdHostAuthNode             = "pcs host auth %s -u hacluster -p '%s'"
	CmdHostAuthNodeWithAddr     = "pcs host auth %s addr=%s -u hacluster -p '%s'"
	CmdDefaultResourceStickness = "pcs resource defaults | grep resource-stickiness --color=never"
	CmdGetPcsdAuthFile          = "cat /var/lib/pcsd/known-hosts"
	CmdGetSbdStatus             = "pcs stonith sbd status --full"

	CmdSaveCIB                 = "pcs cluster cib ra-cfg"
	CmdPushFileToCIB           = "pcs cluster cib-push ra-cfg"
	CmdResourceStop            = "pcs resource disable "
	CmdResourceStart           = "pcs resource enable "
	CmdResourceCleanup         = "pcs resource cleanup"
	CmdResourceUnclone         = "pcs resource unclone "
	CmdResourceUngroup         = "pcs resource ungroup "
	CmdResourceDelete          = "pcs resource delete %s"
	CmdResourceDeleteForce     = CmdResourceDelete + " --force"
	CmdResourceUpdateMeta      = "pcs resource update %s meta %s=%v"
	CmdResourceUpdateMetaForce = CmdResourceUpdateMeta + "--force"
	CmdResourceUpdate          = "pcs resource update %s %s"
	CmdResourceUpdateForce     = CmdResourceUpdate + " --force"

	CmdResourceOpDelete    = "pcs resource op delete %s"
	CmdResourceClone       = "pcs resource clone %s"
	CmdResourceGroupAdd    = "pcs resource group add %s"
	CmdResourceGroupRemove = "pcs resource group remove %s"
	CmdResourceMetaAdd     = "pcs resource meta %s %s=%s"

	CmdLocationDelete           = "pcs constraint location delete %s"
	CmdLocationAdd              = "pcs constraint location %s prefers %s=%s"
	CmdOrderDelete              = "pcs constraint order delete %s"
	CmdColationDelete           = "pcs constraint colocation delete %s %s"
	CmdColationAdd              = "pcs constraint colocation add %s with %s INFINITY"
	CmdRuleDelete               = "pcs constraint rule delete %s"
	CmdRuleAdd                  = "pcs constraint location %s rule score=%s"
	CmdRuleAddWithId            = "pcs constraint location %s rule score=%s id=%s"
	CmdOrderAdd                 = "pcs constraint order start %s then %s"
	CmdListResourceStandards    = "crm_resource --list-standards"
	CmdListOcfProviders         = "crm_resource --list-ocf-providers"
	CmdListOcfResourceAgent     = "crm_resource --list-agents ocf:%s"
	CmdListResourceAgent        = "crm_resource --list-agents %s"
	CmdShowMetaData             = "crm_resource  --show-metadata %s:%s"
	CmdShowMetaDataWithProvider = "crm_resource  --show-metadata %s:%s:%s"
	CmdCrmResource              = "crm_resource --resource %s"
	CmdQueryResourceAsXml       = "crm_resource --resource %s --query-xml"

	CmdGenLog          = "/usr/share/heartbeat-gui/ha-api/loggen.sh"
	CmdPacemakerAgents = "pcs resource agents ocf:pacemaker"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	logs.Debug("Running command: %s", c)
	command := exec.Command("bash", "-c", c)
	out, err := command.CombinedOutput()
	if err != nil {
		// logs.Error("Run command failed!, command: " + c + " out: " + string(out) + " err: " + err.Error())
		return out, errors.Wrapf(err, "Run command failed!, command: "+c+" out: "+string(out)+" err: "+err.Error())
	}
	return out, nil
}
