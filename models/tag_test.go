/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Dec 05 15:26:24 2024 +0800
 */

package models

import (
	"testing"
)

func TestGetTag(t *testing.T) {
	result := GetTag()
	if result.Action != true {
		t.Fatal("Get tag failed")
	}
}
