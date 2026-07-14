/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Tue Jan 5 09:33:00 2021 +0800
 */
package utils

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"github.com/pkg/errors"
)

// ShellEscape safely escapes a string for use as a single shell argument.
// It uses single-quote wrapping, which is the most robust method for POSIX shells.
// Single quotes within the value are escaped by ending the quote, adding an escaped
// single quote, and reopening the quote: 'it'\''s' => it's
func ShellEscape(s string) string {
       return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

const (
	CmdNodeStatus                   = "crm_node -l"
	CmdCrmMon                       = "crm_mon "
	CmdCrmMonAsXML                  = CmdCrmMon + " --as-xml"
	CmdClusterStatus                = "crm_mon -1"
	CmdClusterStatusAsXML           = CmdClusterStatus + " --as-xml"
	CmdHostName                     = "hostname"
	CmdGetCurrentNodeName           = "cat /etc/hostname"
	CmdHostExists                   = "cat /etc/hosts | grep -w %s"
	CmdHostSearch                   = "cat /etc/hosts | grep -i -w %s"
	CmdCountClustersConfigsBackuped = "ls /usr/share/heartbeat-gui/ha-api/ClustersInfo.conf.* | wc -l"
	CmdCibQueryConfig               = "cibadmin --query --scope configuration"
	CmdDeleteAlert                  = "pcs alert delete alert_Kylin"
	CmdDeleteMailAlert              = "pcs alert delete alert_Kylin_mail"
	CmdCreateAlert                  = "pcs alert create id=alert_Kylin path=/usr/share/pacemaker/alerts/alert_log.sh"
	CmdCreateMailAlert              = "pcs alert create id=alert_Kylin_mail path=/usr/share/pacemaker/alerts/kylin_alert.sh"
	CmdAddAlert                     = "pcs alert recipient add alert_Kylin_mail"
	CmdSendEmail                    = "/usr/share/pacemaker/alerts/python_email.py"
	CmdUpdateResourceStickness      = "crm_attribute  -t rsc_defaults -n resource-stickiness -v "
	CmdUpdateCrmConfig              = "crm_attribute -t crm_config -n "
	CmdQueryCIB                     = "cibadmin -Q"
	CmdQueryCrmConfig               = "cibadmin --query --scope crm_config"
	CmdQueryResources               = "cibadmin --query --scope resources"
	CmdQueryConstraints             = "cibadmin --query --scope constraints"
	CmdQueryResourcesById           = "cibadmin --query --scope resources|grep 'id=\"%s\"'"

	CmdStartCluster             = "pcs cluster start --all "
	CmdStartClusterNode         = "pcs cluster start %s"
	CmdStartClusterRemoteNode   = "pcs resource enable %s"
	CmdStopClusterLocal         = "pcs cluster stop "
	CmdStopClusterNode          = "pcs cluster stop %s"
	CmdStopClusterNodeForce     = CmdStopClusterNode + " --force"
	CmdStopClusterRemoteNode    = "pcs resource disable %s"
	CmdStopCluster              = "pcs cluster stop --all"
	CmdSetupCluster             = "pcs cluster setup hacluster"
	CmdSetupClusterStandard     = "pcs cluster setup %s %s totem token=8000 --start"
	CmdDestroyCluster           = "pcs cluster destroy --all"
	CmdDestroyClusterForce      = "pcs cluster destroy --all --force"
	CmdNodeAdd                  = "pcs cluster node add %s "
	CmdNodeAddStart             = CmdNodeAdd + "--start"
	CmdNodeStandby              = "pcs node standby "
	CmdNodeUnStandby            = "pcs node unstandby "
	CmdHostAuthNode             = "pcs host auth %s -u hacluster -p '%s'"
	CmdHostAuthNodeWithAddr     = "pcs host auth %s addr=%s -u hacluster -p '%s'"
	CmdHostAuthLocal            = "pcs client local-auth -u hacluster -p '%s'"
	CmdDefaultResourceStickness = "pcs resource defaults | grep resource-stickiness --color=never"
	CmdGetPcsdAuthFile          = "cat /var/lib/pcsd/known-hosts"
	CmdGetSbdStatus             = "pcs stonith sbd status --full"

	CmdSaveCIB                  = "pcs cluster cib ra-cfg"
	CmdPushFileToCIB            = "pcs cluster cib-push ra-cfg"
	CmdResourceStop             = "pcs resource disable "
	CmdStonithResourceCreate    = "pcs stonith create "
	CmdResourceCreate           = "pcs resource create "
	CmdResourceStart            = "pcs resource enable "
	CmdResourceCleanup          = "pcs resource cleanup"
	CmdResourceUnclone          = "pcs resource unclone "
	CmdResourceUngroup          = "pcs resource ungroup "
	CmdResourceDelete           = "pcs resource delete %s"
	CmdResourceDeleteForce      = CmdResourceDelete + " --force"
	CmdResourceDeleteGuestForce = "pcs cluster node delete-guest %s && pcs resource delete %s --force" // 删除guest资源
	CmdResourceUpdateMeta       = "pcs resource update %s meta %s=%v"
	// CmdResourceUpdateMeta2      = "pcs resource update %s meta %s"
	CmdResourceUpdateMetaForce = CmdResourceUpdateMeta + " --force"
	CmdResourceUpdate          = "pcs resource update %s %s"
	CmdResourceUpdateForce     = CmdResourceUpdate + " --force"
	CmdResourceUpdateOp        = "pcs resource update %s op %s"

	CmdResourceOpDelete    = "pcs resource op delete %s"
	CmdResourceClone       = "pcs resource clone %s"
	CmdResourceGroupAdd    = "pcs resource group add %s"
	CmdResourceGroupRemove = "pcs resource group remove %s"
	CmdResourceMetaAdd     = "pcs resource meta %s %s=%s"

	CmdLocationDelete = "pcs constraint location delete %s"
	CmdMetaDelete     = "crm_resource -r %s -m --delete-parameter %s" // 删除普通资源元属性
	CmdInstanceDelete = "crm_resource -r %s --delete-parameter %s"    // 删除普通资源实例属性
	CmdLocationAdd    = "pcs constraint location %s prefers %s=%s"
	CmdOrderDelete    = "pcs constraint order delete %s"
	CmdColationDelete = "pcs constraint colocation delete %s %s"

	CmdColationAddIN            = "pcs constraint colocation add %s with %s INFINITY"
	CmdColationAddNEIN          = "pcs constraint colocation add %s with %s -INFINITY"
	CmdColationAdd              = "pcs constraint colocation add %s with %s %s"
	CmdColationAddWithRole      = "pcs constraint colocation add %s with %s %s with-rsc-role=%s"
	CmdRuleDelete               = "pcs constraint rule delete %s"
	CmdRuleAdd                  = "pcs constraint location %s rule score=%s"
	CmdRuleAddWithId            = "pcs constraint location %s rule score=%s id=%s"
	CmdOrderAdd                 = "pcs constraint order %s %s then %s %s"
	CmdListResourceStandards    = "crm_resource --list-standards"
	CmdListOcfProviders         = "crm_resource --list-ocf-providers"
	CmdListOcfResourceAgent     = "crm_resource --list-agents ocf:%s"
	CmdListResourceAgent        = "crm_resource --list-agents %s"
	CmdShowMetaData             = "crm_resource  --show-metadata %s:%s"
	CmdShowMetaDataWithProvider = "crm_resource  --show-metadata %s:%s:%s"
	CmdCrmResource              = "crm_resource --resource %s"
	CmdQueryResourceAsXml       = "crm_resource --resource %s --query-xml"

	CmdGenLog              = "/usr/share/heartbeat-gui/ha-api/loggen.sh"
	CmdPacemakerAgents     = "pcs resource agents ocf:pacemaker"
	DefaultSleep           = " &sleep 5"
	CmdDeleteLinks         = "pcs cluster link delete %s"
	CmdAddLink             = "pcs cluster link add %s"
	CmdAddLinkForce        = "pcs cluster link add %s --force"
	CmdAddLinksWithLinkNum = "pcs cluster link add %s  options linknumber=%s --force"
	CmdHbStatus            = "corosync-cfgtool -n"
	CmdHbStatusS           = "corosync-cfgtool -s"
	CmdUpdateLink          = "pcs cluster link update %s %s"
	CmdUpdateLinkForce     = CmdUpdateLink + " --force"
	CmdSyncCorosyncConf    = "pcs cluster sync"
	CmdPcsdStatus          = "pcs status pcsd"
	CmdResourceStatus      = "pcs resource status"
	CmdCheckResourceExist  = "crm_resource --locate --resource %s"

	CmdChangePwd = "echo 'hacluster:%s'|chpasswd >/dev/null 2>&1"
)

// RunChangePwd safely changes the hacluster password without shell injection risk.
// The password is passed through a pipe, never embedded in shell syntax.
var RunChangePwd = func(password string) ([]byte, error) {
	slog.Debug("Running change password command")

	chpasswd := exec.Command("chpasswd")
	chpasswd.Env = append(chpasswd.Environ(), "LANG=C")

	stdin, err := chpasswd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stdin pipe")
	}

	if err := chpasswd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start chpasswd")
	}

	go func() {
		fmt.Fprintf(stdin, "hacluster:%s\n", password)
		stdin.Close()
	}()

	out, err := chpasswd.CombinedOutput()
	if err != nil {
		slog.Error("Change password failed", "out", string(out))
		return out, errors.Wrapf(err, "change password failed, out: %s", string(out))
	}
	return out, nil
}

// RunCommand runs the command and get the result
var RunCommand = func(c string) ([]byte, error) {
	slog.Debug("Running command", "cmd", c)

	command := exec.Command("bash", "-c", c)
	command.Env = append(command.Environ(), "LANG=C")
	out, err := command.CombinedOutput()
	if err != nil {
		slog.Error("Run command failed!", "cmd", c, "out", string(out))
		return out, errors.Wrapf(err, "Run command failed!, command: "+c+" out: "+string(out))
	}
	return out, nil
}

// RunCommandWithArgs executes a command without shell interpretation.
// The binary and each argument are passed directly to exec.Command,
// completely eliminating shell injection risk.
// Use this instead of RunCommand when no shell features (pipes, redirects) are needed.
func RunCommandWithArgs(binary string, args ...string) ([]byte, error) {
        slog.Debug("Running command with args", "binary", binary, "args", args)

        command := exec.Command(binary, args...)
        command.Env = append(command.Environ(), "LANG=C")
        out, err := command.CombinedOutput()
        if err != nil {
                slog.Error("Run command failed!", "binary", binary, "args", args, "out", string(out))
                return out, errors.Wrapf(err, "Run command failed!, binary: %s args: %v out: %s", binary, args, string(out))
        }
        return out, nil
}

