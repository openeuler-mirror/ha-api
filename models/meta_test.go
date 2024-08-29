/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-08-20 16:43:13
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-08-20 17:05:11
 * @FilePath: /ha-api/models/meta_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package models

import "testing"

func TestGetAllResourceMetas(t *testing.T) {
	result := GetAllResourceMetas()
	if result["action"] != true {
		t.Fatal("Get all resource metas failed")
	}
}

func TestGetResourceMetas(t *testing.T) {
	result := GetResourceMetas("heartbeat", "Dummy", "ocf")
	if result["action"] != true {
		t.Fatal("Get ocf heartbeat resource metas failed")
	}

	result2 := GetResourceMetas("stonith", "fence_sbd", "")
	if result2["action"] != true {
		t.Fatal("Get stonith resource metas failed")
	}
}
