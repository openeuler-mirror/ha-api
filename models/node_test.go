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

func TestHandleNodeAction(t *testing.T) {
	cmd := handleNodeAction("start", "primitive", "test", "")
	if cmd != "pcs cluster start test &sleep 5" {
		t.Fatal("Handle Node Action test1 failed")
	}

	cmd = handleNodeAction("stop", "primitive", "test", "")
	if cmd != "pcs cluster stop test &sleep 5" {
		t.Fatal("Handle Node Action test2 failed")
	}

	cmd = handleNodeAction("start", "remote", "test", "")
	if cmd != "pcs resource enable test &sleep 5" {
		t.Fatal("Handle Node Action test3 failed")
	}

	cmd = handleNodeAction("stop", "remote", "test", "")
	if cmd != "pcs resource disable test" {
		t.Fatal("Handle Node Action test4 failed")
	}
	cmd = handleNodeAction("start", "guest", "test", "res")
	if cmd != "pcs resource enable res &sleep 5" {
		t.Fatal("Handle Node Action test3 failed")
	}
	cmd = handleNodeAction("stop", "guest", "test", "res")
	if cmd != "pcs resource disable res" {
		t.Fatal("Handle Node Action test5 failed")
	}
}
