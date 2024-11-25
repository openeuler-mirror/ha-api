/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bixiaoyan
 * Date: 2024-11-22 14:25:51
 * LastEditTime: 2024-11-22 14:37:48
 * Description: 集群利用率功能单元测试
 ******************************************************************************/

package models

import "testing"

func TestGetUtilization(t *testing.T) {
	result := GetUtilization()
	if result.Action != true {
		t.Fatal("Get utilization failed")
	}
}
