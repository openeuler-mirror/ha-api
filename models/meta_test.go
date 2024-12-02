/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Aug 29 15:31:29 2024 +0800
 */
package models

import "testing"

func TestGetAllResourceMetas(t *testing.T) {
	result := GetAllResourceMetas()
	if result["action"] != true {
		t.Fatal("Get all resource metas failed")
	}
}

func TestGetResourceMetas(t *testing.T) {
	result := GetResourceMetas("heartbeat", "Dummy", "ocf")
	if result["action"] != true {
		t.Fatal("Get ocf heartbeat resource metas failed")
	}

	result2 := GetResourceMetas("stonith", "fence_sbd", "")
	if result2["action"] != true {
		t.Fatal("Get stonith resource metas failed")
	}
}
