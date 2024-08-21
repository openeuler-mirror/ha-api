/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-08-20 14:54:10
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-08-20 14:57:23
 * @FilePath: /ha-api/models/log_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package models

import "testing"

func TestGenerateLog(t *testing.T) {
	result := GenerateLog()
	if result["action"] != true {
		t.Fatal("Generate log failed")
	}
}
