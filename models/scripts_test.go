/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Mon Dec 02 17:44:18 2024 +0800
 */
 package models

 import (
	 "testing"
 )

 func TestIsScriptExist(t *testing.T) {
	scriptName := "/usr/lib/ocf/resource.d/pacemaker/Dummy"
	result := IsScriptExist();
	if result["action"] != true {
		t.Fatal("Get alarm failed")
	}
}