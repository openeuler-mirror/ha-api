/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Wed Aug 21 09:42:19 2024 +0800
 */
package models

import "testing"

func TestGenerateLog(t *testing.T) {
	result := GenerateLog()
	if result["action"] != true {
		t.Fatal("Generate log failed")
	}
}
