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

	"github.com/stretchr/testify/assert"
)

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
