/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

// func TestGetNodesInfo(t *testing.T) {
// 	_, err := GetNodesInfo()
// 	if err != nil {
// 		t.Fatal("Get Node Info failed")
// 	}
// }

// func TestGetNodeIDInfo(t *testing.T) {
// 	_, err := GetNodeIDInfo("host1")
// 	if err != nil {
// 		t.Fatal("Get Node ID Info failed")
// 	}
// }

// func TestHandleNodeAction(t *testing.T) {
// 	cmd := handleNodeAction("start", "primitive", "test", "")
// 	if cmd != "pcs cluster start test &sleep 5" {
// 		t.Fatal("Handle Node Action test1 failed")
// 	}

// 	cmd = handleNodeAction("stop", "primitive", "test", "")
// 	if cmd != "pcs cluster stop test &sleep 5" {
// 		t.Fatal("Handle Node Action test2 failed")
// 	}

// 	cmd = handleNodeAction("start", "remote", "test", "")
// 	if cmd != "pcs resource enable test &sleep 5" {
// 		t.Fatal("Handle Node Action test3 failed")
// 	}

// 	cmd = handleNodeAction("stop", "remote", "test", "")
// 	if cmd != "pcs resource disable test" {
// 		t.Fatal("Handle Node Action test4 failed")
// 	}
// 	cmd = handleNodeAction("start", "guest", "test", "res")
// 	if cmd != "pcs resource enable res &sleep 5" {
// 		t.Fatal("Handle Node Action test3 failed")
// 	}
// 	cmd = handleNodeAction("stop", "guest", "test", "res")
// 	if cmd != "pcs resource disable res" {
// 		t.Fatal("Handle Node Action test5 failed")
// 	}
// }

import (
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock 工具类
type MockCmdRunner struct{ mock.Mock }

func (m *MockCmdRunner) Run(cmd string) ([]byte, error) {
	args := m.Called(cmd)
	return args.Get(0).([]byte), args.Error(1)
}

func TestGetNodeIDInfo_Success(t *testing.T) {
	// 备份并替换 RunCommand
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 模拟成功的命令执行
	utils.RunCommand = func(cmd string) ([]byte, error) {
		xml := `
	192.0.2.1`
		return []byte(xml), nil
	}
	result, err := GetNodeIDInfo("node1")

	// 验证
	assert.NoError(t, err)
	// assert.Equal(t, 2, len(result["ips"]))
	assert.Contains(t, result["ips"], "192.0.2.1")
	// assert.Contains(t, result["ips"], "192.168.1.2")
}

func TestGetNodesInfo_Success(t *testing.T) {
	// 备份并替换 RunCommand
	originalRunCommand := utils.RunCommand
	defer func() { utils.RunCommand = originalRunCommand }()

	// 模拟成功的命令执行
	utils.RunCommand = func(cmd string) ([]byte, error) {
		xml := `
	<crm_mon version="2.1.9-1.p04.ky11">
  <summary>
    <stack type="corosync" pacemakerd-state="running"/>
    <current_dc present="true" version="2.1.9-1.p04.ky11-49aab9983" name="ha1" id="1" with_quorum="true" mixed_version="false"/>
    <last_update time="Thu May 29 11:05:37 2025" origin="ha1"/>
    <last_change time="Thu May 29 09:45:30 2025" user="root" client="root" origin="ha1"/>
    <nodes_configured number="2"/>
    <resources_configured number="1" disabled="0" blocked="0"/>
    <cluster_options stonith-enabled="true" symmetric-cluster="true" no-quorum-policy="stop" maintenance-mode="false" stop-all-resources="false" stonith-timeout-ms="60000" priority-fencing-delay-ms="0"/>
  </summary>
  <nodes>
    <node name="node1" id="1" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" health="green" feature_set="3.19.6" shutdown="false" expected_up="true" is_dc="true" resources_running="0" type="member"/>
    <node name="node2" id="2" online="true" standby="true" standby_onfail="false" maintenance="false" pending="false" unclean="false" health="green" shutdown="false" expected_up="false" is_dc="false" resources_running="0" type="member"/>
  </nodes>
  <resources>
    <resource id="d" resource_agent="ocf::heartbeat:Dummy" role="Stopped" active="false" orphaned="false" blocked="false" maintenance="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="0"/>
  </resources>
  <node_history/>
</crm_mon>
`
		return []byte(xml), nil
	}

	// 执行测试
	result, err := GetNodesInfo()

	// 验证
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	// 验证第一个节点
	assert.Equal(t, "node1", result[0]["id"])
	assert.Equal(t, "Master", result[0]["status"])
	assert.Equal(t, "true", result[0]["is_dc"])

	// 验证第二个节点
	assert.Equal(t, "node2", result[1]["id"])
	assert.Equal(t, "Standby", result[1]["status"])
	assert.Equal(t, "false", result[1]["is_dc"])
}

func TestDoNodeAction(t *testing.T) {
	mockRunner := new(MockCmdRunner)
	utils.RunCommand = mockRunner.Run // 注入Mock

	testCases := []struct {
		name       string
		nodeID     string
		action     string
		data       map[string]string
		mockCmd    string // 预期执行的命令
		mockErr    error  // 模拟命令错误
		wantAction bool   // 期望action结果
		wantErrMsg string // 期望的错误信息
	}{
		// 成功场景
		{
			name:       "standby_success",
			nodeID:     "node1",
			action:     "standby",
			mockCmd:    utils.CmdNodeStandby + "node1",
			wantAction: true,
		},
		{
			name:       "unstandy_success",
			nodeID:     "node1",
			action:     "unstandby",
			mockCmd:    utils.CmdNodeUnStandby + "node1",
			wantAction: true,
		},
		{
			name:       "start_success",
			nodeID:     "node1",
			action:     "start",
			data:       map[string]string{"type": "primitive", "res": ""},
			mockCmd:    handleNodeAction("start", "node1", "primitive", ""),
			wantAction: true,
		},
		{
			name:       "stop_success",
			nodeID:     "node1",
			action:     "stop",
			data:       map[string]string{"type": "primitive", "res": ""},
			mockCmd:    handleNodeAction("stop", "node1", "primitive", ""),
			wantAction: true,
		},
		// {
		// 	name:       "start_primitive",
		// 	action:     "start",
		// 	nodeID:     "node2",
		// 	data:       map[string]string{"type": "primitive", "res": ""},
		// 	mockCmd:    "pcs cluster start node2" + utils.DefaultSleep,
		// 	wantAction: true,
		// },
		// 失败场景
		// {
		// 	name:       "command_failure",
		// 	action:     "unstandby",
		// 	nodeID:     "node3",
		// 	mockCmd:    utils.CmdNodeUnStandby + "node3",
		// 	mockErr:    errors.New("connection timeout"),
		// 	wantAction: false,
		// 	wantErrMsg: "Change node status Failed",
		// },
		// 特殊逻辑
		// {
		// 	name:       "stop_guest_no_sleep",
		// 	action:     "stop",
		// 	nodeID:     "node4",
		// 	data:       map[string]string{"type": "guest", "res": "res1"},
		// 	mockCmd:    "pcs resource disable res1",
		// 	wantAction: true,
		// },
		{
			name:       "restart_chain_commands",
			action:     "restart",
			nodeID:     "node5",
			mockCmd:    utils.CmdStopCluster + "node5 && " + utils.CmdStopCluster + "node5",
			wantAction: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// 设置Mock预期
			if tt.mockCmd != "" {
				mockRunner.On("Run", tt.mockCmd).Return([]byte{}, tt.mockErr)
			}

			// 执行测试
			result := DoNodeAction(tt.nodeID, tt.action, tt.data)

			// 验证结果
			assert.Equal(t, tt.wantAction, result["action"])
			if tt.wantErrMsg != "" {
				assert.Contains(t, result["error"], tt.wantErrMsg)
			} else {
				assert.Nil(t, result["error"])
			}

			// 确保所有Mock期望被调用
			mockRunner.AssertExpectations(t)
		})
	}
}

func TestHandleNodeAction(t *testing.T) {
	testCases := []struct {
		name     string
		action   string
		nodeType string
		nodeID   string
		nodeRes  string
		wantCmd  string
	}{
		// 正常分支
		{
			name:     "start_remote",
			action:   "start",
			nodeType: "remote",
			nodeID:   "node1",
			wantCmd:  "pcs resource enable node1" + utils.DefaultSleep,
		},
		{
			name:     "stop_primitive",
			action:   "stop",
			nodeType: "primitive",
			nodeID:   "node2",
			wantCmd:  "pcs cluster stop node2 --force &sleep 5",
		},
		// 异常分支
		{
			name:     "invalid_action",
			action:   "invalid",
			nodeType: "primitive",
			wantCmd:  "invalid action: invalid",
		},
		{
			name:     "invalid_type",
			action:   "start",
			nodeType: "unknown",
			wantCmd:  "invalid node type: unknown for action: start",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			cmd := handleNodeAction(tt.action, tt.nodeID, tt.nodeType, tt.nodeRes)
			assert.Equal(t, tt.wantCmd, cmd)
		})
	}
}
