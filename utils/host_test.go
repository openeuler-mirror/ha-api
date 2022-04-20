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
 * LastEditTime: 2022-04-20 11:18:32
 * Description: corosync配置获取的测试用例
 ******************************************************************************/
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
