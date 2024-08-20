/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2024-08-20 10:30:17
 * Description: 集群测试用例
 ******************************************************************************/
package models

import "testing"

func TestGetClusterPropertiesInfo(t *testing.T) {
	result := GetClusterPropertiesInfo()
	if result["action"] != true {
		t.Fatal("Get cluster properties failed")
	}
}

func TestUpdateClusterProperties(t *testing.T) {
	clusterPropJson := map[string]interface{}{}
	clusterPropJson["no-quorum-policy"] = "ignore"
	print(clusterPropJson)
	res := UpdateClusterProperties(clusterPropJson)
	if res["action"] != true {
		t.Fatal("Update cluster properties failed")
	}
}
