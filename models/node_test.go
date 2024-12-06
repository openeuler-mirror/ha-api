/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"testing"
)

func TestGetNodesInfo(t *testing.T) {
	_, err := GetNodesInfo()
	if err != nil {
		t.Fatal("Get Node Info failed")
	}
}

func TestGetNodeIDInfo(t *testing.T) {
	_, err := GetNodeIDInfo("host1")
	if err != nil {
		t.Fatal("Get Node ID Info failed")
	}
}
