/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-13 14:16:22
 * LastEditTime: 2024-03-13 14:57:14
 * Description:脚本生成器模块相关接口
 */
package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
)

const (
	pacemakerAgentsCmd = "pcs resource agents ocf:pacemaker"
)

type ScriptResponse struct {
	utils.GeneralResponse
	Data map[string]string `json:"data,omitempty"`
}

func IsScriptExist(scriptName string) utils.GeneralResponse {
	out, err := utils.RunCommand(pacemakerAgentsCmd)
	if err != nil {
		return utils.HandleCmdError("查询脚本命令执行失败", false)
	}
	scripts := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, script := range scripts {
		if script == scriptName {
			logs.Warn(fmt.Sprintf("脚本 %s 已存在于pacemaker目录下", scriptName))
			return utils.GeneralResponse{
				Action: false,
				Error:  "脚本已经存在于pacemaker目录下",
			}
		}
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   "脚本不存在",
	}
}

func GenerateLocalScript(data map[string]string) utils.GeneralResponse {
	name := data["name"]
	startCommand := data["start"]
	stopCommand := data["stop"]
	monitorCommand := data["monitor"]
	filePath := filepath.Join("/usr/lib/ocf/resource.d/pacemaker/", name)
	scriptTemplete :=
		`#!/bin/sh
#
# Author:	Kylin
# License:      GNU General Public License (GPL) 
#
#	usage: $0 {start|stop|status|monitor|validate-all|meta-data}
#  OCF parameters:
#
##########################################################################
# Initialization:
#OCF_ROOT="/usr/lib/ocf"
: ${OCF_FUNCTIONS_DIR=${OCF_ROOT}/lib/heartbeat}
. ${OCF_FUNCTIONS_DIR}/ocf-shellfuncs
#########################################################################

name="` + name + `"
start_path="` + startCommand + `"
stop_path="` + stopCommand + `"
meta_data() {
	cat <<END
<?xml version="1.0"?>
<!DOCTYPE resource-agent SYSTEM "ra-api-1.dtd">
<resource-agent name="$name">
<version>1.0</version>
<longdesc lang="en">
This script manages $name
</longdesc>
<shortdesc lang="en">OCF Resource Agent compliant $name.</shortdesc>
<parameters>
</parameters>
<actions>
<action name="start"    timeout="50" />
<action name="stop"     timeout="60" />
<action name="monitor"  depth="10"  timeout="15s" interval="30s" start-delay="1s" />
<action name="validate-all"  timeout="30s" />
<action name="meta-data"  timeout="5s" />
</actions>
</resource-agent>
END
}

anything_status() {` + "\n\tret0=`ps -ef |grep -v grep |grep " + monitorCommand + "|wc -l" + "`" + `
	if [ 1 -a $ret0 -gt 0 ]; then
		return $OCF_SUCCESS;
	elif [ 1 -a $ret0 -eq 0 ]; then
		return $OCF_NOT_RUNNING;
	else
		return $OCF_ERR_GENERIC;
	fi
}

anything_start() {
	anything_status
	ret=$?
	if [ $ret -eq 0 ]; then
		ocf_log info "$name agent already running"
		return $OCF_SUCCESS;
	else
		$start_path >/dev/null 2>&1
		anything_monitor
		ret1=$?
		if [ $ret1 -eq 0 ]; then
			ocf_log info "$name agent start ok"
			return $OCF_SUCCESS;
		else
			ocf_log info "$name agent start failed"
			return $OCF_ERR_GENERIC;
			fi
		fi
}

anything_stop()
{
	anything_status
	ret=$?
	if [ $ret -ne 7 ]; then
		# ocf_log info $ret
		$stop_path >/dev/null 2>&1
		sleep 3 
		anything_monitor
		ret=$?
		if [ $ret -eq 7 ]; then
		ocf_log info "$name agent stop ok"
		return $OCF_SUCCESS;
		else
		ocf_log info "$name agent stop failed"
		return $OCF_ERR_GENERIC;
		fi
	else
		ocf_log info "$name agent not running"
		return $OCF_SUCCESS;
	fi
	
}
	
anything_monitor()
{
	anything_status
	ret=$?
	#ocf_log info $ret
	if [ $ret -eq 0 ]; then
	ocf_log info " $name is running!"
	return $OCF_SUCCESS
	elif [ $ret -eq 7 ]; then
	ocf_log info "$name is not running" 
	return $OCF_NOT_RUNNING
	else
	ocf_log info "$name is failed"
	return $OCF_ERR_GENERIC
	fi
}
	
anything_validate()
{
	return $OCF_SUCCESS
}
	
case $1 in
	start)	anything_start
		;;
	
	stop)	anything_stop
		;;
		
	status)	anything_status
		;;
		
	monitor)	anything_monitor
		;;
		
	validate-all)	anything_validate
		exit  $OCF_SUCCESS
		;;
			
	meta-data)	meta_data
		exit  $OCF_SUCCESS
		;;
	*)
		ocf_log err "$0 was called with unsupported arguments: $*"
		exit $OCF_ERR_UNIMPLEMENTED
		;;
esac
exit $?`

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return utils.GeneralResponse{
			Action: false, Error: err.Error()}
	}
	defer file.Close() // 确保文件在函数结束时关闭

	if _, err := file.WriteString(scriptTemplete); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}

	if err := os.Chmod(filePath, 0755); err != nil {
		return utils.GeneralResponse{Action: false, Error: err.Error()}
	}
	// 返回成功信息
	return utils.GeneralResponse{Action: true, Info: "Generate script success!"}
}

func GenerateScript(data map[string]string) ScriptResponse {
	// 获取所有节点
	out, err := utils.RunCommand(CmdNodeStatus)
	if err != nil {
		return ScriptResponse{
			GeneralResponse: utils.GeneralResponse{
				Action: true,
				Error:  err.Error(),
			},
		}
	}
	// 获取当前节点名称
	nodes := strings.Split(strings.TrimSpace(string(out)), "\n")
	out, err = utils.RunCommand(CmdHostName)
	if err != nil {
		return ScriptResponse{
			GeneralResponse: utils.GeneralResponse{
				Action: true,
				Error:  err.Error(),
			},
		}
	}

	result := make(map[string]string)
	localHostName := strings.TrimSpace(string(out))
	for _, node := range nodes {
		nodeName := strings.Split(node, " ")[1]
		if nodeName == localHostName {
			// 生成当前节点的脚本
			res := GenerateLocalScript(data)
			if !res.Action {
				result[nodeName] = "Generate script failed"
				logs.Error(fmt.Sprintf("Generate local script failed, err msg: %s", res.Error))
				continue
			} else {
				result[nodeName] = "Generate script success"
			}
		} else {
			// 生成集群其他节点的脚本
			remoteUrl := fmt.Sprintf("http://%s:%d/api/v1/remotescripts", nodeName, settings.HAAPI_DEFAULT_PORT)
			data["nodecall"] = localHostName
			body, err := json.Marshal(data)
			if err != nil {
				return ScriptResponse{
					GeneralResponse: utils.GeneralResponse{
						Action: true,
						Error:  err.Error(),
					},
				}
			}
			headers := map[string]string{
				"Content-Type": "application/json",
			}

			resp, err := http.Post(remoteUrl, headers["Content-Type"], bytes.NewBuffer(body))
			if err != nil {
				result[nodeName] = "Generate script failed"
				logs.Error(err)
				continue
			}
			defer resp.Body.Close() // 确保请求在函数结束时关闭

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				result[nodeName] = "Generate script failed"
				continue
			}

			var retMap map[string]interface{}
			if err = json.Unmarshal(body, &retMap); err != nil {
				result[nodeName] = "Generate script failed!"
				continue
			}
			retMap["action"] = retMap["action"].(bool)
			if retMap["action"].(bool) {
				result[nodeName] = "Generate script success!"
			} else {
				result[nodeName] = "Generate script failed"
			}
		}
	}
	return ScriptResponse{
		Data: result,
		GeneralResponse: utils.GeneralResponse{
			Action: true,
		},
	}

}
