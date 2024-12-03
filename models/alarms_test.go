/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Mon Sep 2 15:58:49 2024 +0800
 */
package models

import "testing"

func TestAlarmsGet(t *testing.T) {
	result := AlarmsGet()
	if result["action"] != true {
		t.Fatal("Get alarm failed")
	}
}
