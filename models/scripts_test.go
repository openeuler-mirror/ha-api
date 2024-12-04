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

func TestIsScriptExistSuccess(t *testing.T) {
	scriptName := "Dummy"
	result := IsScriptExist(scriptName)
	if result.Action == true {
		t.Fatal("IsScriptExistSuccess test failed")
	}
}

func TestIsScriptExistFailed(t *testing.T) {
	scriptName := "Dummy-not-exist"
	result := IsScriptExist(scriptName)
	if result.Action == false {
		t.Fatal("IsScriptExistFailed test failed")
	}
}
func TestGenerateLocalScriptSuccess(t *testing.T) {
	scriptData := map[string]string{
		"name":    "testScript",
		"start":   "systemctl start testScript",
		"stop":    "systemctl stop testScript",
		"monitor": "systemctl status testScript",
	}
	err := GenerateLocalScript(scriptData)
	if err != nil {
		t.Fatal("GenerateLocalScript test failed")
	}
}
