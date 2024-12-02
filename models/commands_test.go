/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Tue Aug 20 16:27:26 2024 +0800
 */
package models

import "testing"

func TestGetCommandsList(t *testing.T) {
	result := GetCommandsList()
	if result["action"] != true {
		t.Fatal("Get commands list failed")
	}
}

func TestRunBuiltinCommand(t *testing.T) {
	_, res := RunBuiltinCommand(1)
	if res != nil {
		t.Fatal("Get command result failed")
	}
}
