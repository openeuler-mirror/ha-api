/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-08-20 14:00:58
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-08-20 14:49:01
 * @FilePath: /ha-api/models/commands_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
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
