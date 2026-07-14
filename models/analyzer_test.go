/*
 * Copyright (c) KylinSoft  Co., Ltd. 2027.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: xuxiaojuan <xuxiaojuan@kylinos.cn>
 * Date: Wed July 8 13:56:40 2026 +0800
 */

package models

import (
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

var analyzerOrigRunCommand = utils.RunCommand

func analyzerMockRunCommand(mock func(string) ([]byte, error)) {
	utils.RunCommand = mock
}

func analyzerRestoreRunCommand() {
	utils.RunCommand = analyzerOrigRunCommand
}

// ==================== ParseVmstat ====================

func TestParseVmstat_ValidOutput(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  1  1024  50000  12000  30000    0    1    10    20  100  200  5  3 90  2  0`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	assert.Equal(t, 2, stat.Procs.R)
	assert.Equal(t, 1, stat.Procs.B)
	assert.Equal(t, 1024, stat.Memory.Swpd)
	assert.Equal(t, 50000, stat.Memory.Free)
	assert.Equal(t, 12000, stat.Memory.Buff)
	assert.Equal(t, 30000, stat.Memory.Cache)
	assert.Equal(t, 0, stat.Swap.Si)
	assert.Equal(t, 1, stat.Swap.So)
	assert.Equal(t, 10, stat.IO.Bi)
	assert.Equal(t, 20, stat.IO.Bo)
	assert.Equal(t, 100, stat.IO.In)
	assert.Equal(t, 200, stat.IO.Cs)
	assert.InDelta(t, 5.0, stat.CPU.Us, 0.01)
	assert.InDelta(t, 3.0, stat.CPU.Sy, 0.01)
	assert.InDelta(t, 90.0, stat.CPU.Id, 0.01)
	assert.InDelta(t, 2.0, stat.CPU.Wa, 0.01)
}

func TestParseVmstat_EmptyInput(t *testing.T) {
	stat, err := ParseVmstat("")
	assert.NoError(t, err)
	assert.Equal(t, 0, stat.Procs.R)
	assert.Equal(t, 0, stat.Procs.B)
	assert.Equal(t, 0, stat.Memory.Swpd)
}

func TestParseVmstat_OnlyHeaders(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	assert.Equal(t, 0, stat.Procs.R)
}

func TestParseVmstat_InsufficientFields(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  1  1024  50000`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	// 字段不足的行会被跳过，所有值保持为零
	assert.Equal(t, 0, stat.Procs.R)
}

func TestParseVmstat_InvalidNumericFields(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 a  b  cdef  50000  12000  30000    0    1    10    20  100  200  5  3 90  2  0`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	// 解析失败的行会被跳过
	assert.Equal(t, 0, stat.Procs.R)
}

func TestParseVmstat_MultipleLines_LastWins(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0   512  40000  10000  25000    0    0     5    10   50  100  2  1 95  2  0
 3  2  2048  60000  15000  35000    1    2    15    25  150  250  8  5 85  2  0`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	// 最后一行的数据会覆盖前面的
	assert.Equal(t, 3, stat.Procs.R)
	assert.Equal(t, 2, stat.Procs.B)
	assert.Equal(t, 2048, stat.Memory.Swpd)
	assert.Equal(t, 60000, stat.Memory.Free)
}

func TestParseVmstat_BlankLinesAndHeaders(t *testing.T) {
	output := `
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st

 4  0     0  80000   5000  40000    0    0     0     0   80  160  1  1 98  0  0
`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	assert.Equal(t, 4, stat.Procs.R)
	assert.Equal(t, 0, stat.Memory.Swpd)
}

func TestParseVmstat_FloatCPUFields(t *testing.T) {
	output := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 0  0     0  99999      0      0    0    0     0     0    0    0 12.5 6.3 80.1 1.1 0`

	stat, err := ParseVmstat(output)
	assert.NoError(t, err)
	assert.InDelta(t, 12.5, stat.CPU.Us, 0.01)
	assert.InDelta(t, 6.3, stat.CPU.Sy, 0.01)
	assert.InDelta(t, 80.1, stat.CPU.Id, 0.01)
	assert.InDelta(t, 1.1, stat.CPU.Wa, 0.01)
}

// ==================== parseStatusOutput ====================

func TestParseStatusOutput_FullStatus(t *testing.T) {
	output := `● corosync.service - Corosync Cluster Engine
   Loaded: loaded (/usr/lib/systemd/system/corosync.service; enabled; vendor preset: disabled)
   Active: active (running) since Mon 2026-07-07 10:30:00 CST; 2 days ago
 Main PID: 12345 (corosync)
   Memory: 25.6M
   CGroup: /system.slice/corosync.service

Jul 07 10:30:00 node1 corosync[12345]: Starting corosync
Jul 07 10:30:01 node1 corosync[12345]: error: failed to connect to cluster
Jul 08 09:15:00 node1 corosync[12345]: Cluster is operational`

	status, err := parseStatusOutput("corosync", output)
	assert.NoError(t, err)
	assert.Equal(t, "corosync", status.Name)
	assert.Equal(t, "active (running)", status.Active)
	assert.Equal(t, "Mon 2026-07-07 10:30:00 CST", status.RunningSince)
	assert.Equal(t, "12345", status.ProcessID)
	assert.Equal(t, "25.6M", status.MemoryUsage)
	assert.Len(t, status.Logs, 1)
	assert.Contains(t, status.Logs[0], "error: failed to connect to cluster")
}

func TestParseStatusOutput_ActiveFallback(t *testing.T) {
	output := `● pacemaker.service - Pacemaker High Availability Cluster Manager
   Loaded: loaded (/usr/lib/systemd/system/pacemaker.service; enabled)
   Active: inactive (dead)`

	status, err := parseStatusOutput("pacemaker", output)
	assert.NoError(t, err)
	assert.Equal(t, "pacemaker", status.Name)
	assert.Equal(t, "inactive (dead)", status.Active)
	assert.Empty(t, status.RunningSince)
}

func TestParseStatusOutput_NoActiveField(t *testing.T) {
	output := `● corosync.service - Corosync Cluster Engine
   Loaded: loaded (/usr/lib/systemd/system/corosync.service; enabled)
 Main PID: 9999 (corosync)
   Memory: 10.0M`

	status, err := parseStatusOutput("corosync", output)
	assert.Nil(t, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Active field found")
}

func TestParseStatusOutput_NoMainPID(t *testing.T) {
	output := `● pcsd.service - PCS GUI and REST interface
   Loaded: loaded (/usr/lib/systemd/system/pcsd.service; disabled)
   Active: active (running) since Tue 2026-07-08 08:00:00 CST; 1 day ago`

	status, err := parseStatusOutput("pcsd", output)
	assert.NoError(t, err)
	assert.Equal(t, "active (running)", status.Active)
	assert.Empty(t, status.ProcessID)
}

func TestParseStatusOutput_NoMemory(t *testing.T) {
	output := `● ha-api.service - HA REST API
   Loaded: loaded (/usr/lib/systemd/system/ha-api.service; enabled)
   Active: active (running) since Wed 2026-07-09 12:00:00 CST; 5h ago
 Main PID: 54321 (ha-api)`

	status, err := parseStatusOutput("ha-api", output)
	assert.NoError(t, err)
	assert.Equal(t, "active (running)", status.Active)
	assert.Equal(t, "54321", status.ProcessID)
	assert.Empty(t, status.MemoryUsage)
}

func TestParseStatusOutput_LogsWithFail(t *testing.T) {
	output := `● corosync.service - Corosync Cluster Engine
   Active: active (running) since Mon 2026-07-07 10:30:00 CST; 2 days ago

Jul 07 10:30:00 node1 corosync[100]: INFO: cluster started
Jul 07 10:31:00 node1 corosync[100]: Failed to join the cluster
Jul 07 10:32:00 node1 corosync[100]: WARNING: timeout occurred
Jul 07 10:33:00 node1 corosync[100]: ERROR: quorum lost`

	status, err := parseStatusOutput("corosync", output)
	assert.NoError(t, err)
	// "Failed" 和 "ERROR" 会被捕获，"WARNING" 不会
	assert.Len(t, status.Logs, 2)
	assert.Contains(t, status.Logs[0], "Failed to join the cluster")
	assert.Contains(t, status.Logs[1], "ERROR: quorum lost")
}

func TestParseStatusOutput_EmptyLogs(t *testing.T) {
	output := `● pacemaker.service - Pacemaker
   Active: active (running) since Mon 2026-07-07 10:30:00 CST; 2 days ago

Jul 07 10:30:00 node1 pacemaker[200]: INFO: all good
Jul 07 10:31:00 node1 pacemaker[200]: INFO: resources started`

	status, err := parseStatusOutput("pacemaker", output)
	assert.NoError(t, err)
	assert.Empty(t, status.Logs)
}

func TestParseStatusOutput_EmptyOutput(t *testing.T) {
	status, err := parseStatusOutput("corosync", "")
	assert.Nil(t, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Active field found")
}

// ==================== ParseSystemctlStatus ====================

func TestParseSystemctlStatus_DisallowedServiceName(t *testing.T) {
	status, err := ParseSystemctlStatus("sshd")
	assert.Nil(t, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service name not allowed")
}

func TestAllowedServiceNames_AllPresent(t *testing.T) {
	for _, name := range []string{"corosync", "pacemaker", "pcsd", "ha-api"} {
		assert.True(t, allowedServiceNames[name], "%s should be allowed", name)
	}
}

func TestAllowedServiceNames_Rejected(t *testing.T) {
	disallowed := []string{"nginx", "mysql", "sshd", "docker", ""}
	for _, name := range disallowed {
		assert.False(t, allowedServiceNames[name], "%q should not be allowed", name)
	}
}

// ==================== allowedServiceNames ====================

func TestAllowedServiceNames_ExactMatch(t *testing.T) {
	expected := map[string]bool{
		"corosync":  true,
		"pacemaker": true,
		"pcsd":      true,
		"ha-api":    true,
	}
	assert.Equal(t, expected, allowedServiceNames)
}

// ==================== getRunResult ====================

func TestGetRunResult_ValidCommand(t *testing.T) {
	defer analyzerRestoreRunCommand()
	analyzerMockRunCommand(func(cmd string) ([]byte, error) {
		return []byte("  Linux 5.10.0-60 \n"), nil
	})

	result, err := getRunResult("uname -r")
	assert.NoError(t, err)
	assert.Equal(t, "Linux 5.10.0-60", result)
}

func TestGetRunResult_CommandError(t *testing.T) {
	defer analyzerRestoreRunCommand()
	analyzerMockRunCommand(func(cmd string) ([]byte, error) {
		return nil, assert.AnError
	})

	result, err := getRunResult("nonexistent_command")
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestGetRunResult_SingleLineOutput(t *testing.T) {
	defer analyzerRestoreRunCommand()
	analyzerMockRunCommand(func(cmd string) ([]byte, error) {
		return []byte("single line output"), nil
	})

	result, err := getRunResult("some_command")
	assert.NoError(t, err)
	assert.Equal(t, "single line output", result)
}

func TestGetRunResult_MultiLineOutput_TakesFirstLine(t *testing.T) {
	defer analyzerRestoreRunCommand()
	analyzerMockRunCommand(func(cmd string) ([]byte, error) {
		return []byte("first line\nsecond line\nthird line"), nil
	})

	result, err := getRunResult("some_command")
	assert.NoError(t, err)
	assert.Equal(t, "first line", result)
}
