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

func TestAlarmsSet(t *testing.T) {
	alarmData := map[string]interface{}{
		"flag":     true,
		"smtp":     "smtp.163.com",
		"port":     25,
		"sender":   "testUser@163.com",
		"password": "testPasswd",
	}
	alarmData["receiver"] = []string{"test@163.com", "test2@163.com"}
	result := AlarmsSet(alarmData)
	if result["action"] != true {
		t.Fatal("Set alarm failed")
	}
}
