/*
* @Author: bixiaoyan bixiaoyan@kylinos.cn
* @Date: 2024-08-20 13:52:22
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-09-02 15:26:58
* @FilePath: /ha-api/models/alarms_test.go
* @Description: 高可用集群报警配置模块单元测试
 */
package models

import "testing"

func TestAlarmsGet(t *testing.T) {
	result := AlarmsGet()
	if result["action"] != true {
		t.Fatal("Get alarm failed")
	}
}
