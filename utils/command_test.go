/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

import (
	"fmt"
	"testing"
)

func TestRunCommand(t *testing.T) {
	cmd := "echo \"hello world\""
	out, err := RunCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(out))

	cmd2 := "echo \"hello world\" | grep -o world"
	out, err = RunCommand(cmd2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(out))
}
