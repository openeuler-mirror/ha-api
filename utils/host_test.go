/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 22 17:27:48 2021 +0800
 */
package utils

import (
	"fmt"
	"testing"
)

func TestGetCorosyncConfig(t *testing.T) {
	result, err := getCorosyncConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
