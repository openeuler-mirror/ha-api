/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Fri Nov 15 16:03:14 2024 +0800
 */

package models

import "testing"

func TestGetNodeList(t *testing.T) {
	result := getNodeList()
	if len(result) == 0 {
		t.Fatal("Get node list failed")
	}
}
