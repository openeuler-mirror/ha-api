/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-08-20 14:25:08
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-11-15 14:32:44
 * @FilePath: /ha-api/models/config_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package models

import "testing"

func TestGetNodeList(t *testing.T) {
	result := getNodeList()
	if len(result) == 0 {
		t.Fatal("Get node list failed")
	}
}
