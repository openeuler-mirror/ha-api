/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 10:30:21
 * Description: 磁盘心跳测试用例
 ******************************************************************************/
package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetHbByHost(t *testing.T) {
	// 	[root@ha1 ~]# cat /var/lib/pcsd/known-hosts
	jsonData := `{
		"format_version": 1,
		"data_version": 14,
		"known_hosts": {
		  "ha1": {
			"dest_list": [
			  {
				"addr": "ha1",
				"port": 2224
			  }
			],
			"token": "0d406755-42ff-42e9-a44d-52f90243e3e9"
		  },
		  "ha2": {
			"dest_list": [
			  {
				"addr": "ha2",
				"port": 2224
			  }
			],
			"token": "8ea37a23-cebb-43e6-97ae-9ea0ccb64da1"
		  },
		  "ha3": {
			"dest_list": [
			  {
				"addr": "ha3",
				"port": 2224
			  }
			],
			"token": "f6d2247a-ac96-4516-afb5-bdae6a076fb2"
		  }
		}
	  }`

	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestHeartBeatInfo(t *testing.T) {
	jsonData := `{
		"hbaddrs1": [
			{
				"ip": "192.168.100.187",
				"nodeid": "ns187"
			},
			{
				"ip": "192.168.100.188",
				"nodeid": "ns188"
			}
		],
		"hbaddrs2": [
			{
				"ip": "192.168.100.187",
				"nodeid": "ns187"
			},
			{
				"ip": "192.168.100.188",
				"nodeid": "ns188"
			}
		],
		"hbaddrs2_enabled": 1
	}`

	// data := map[string]([]HostInfo){}
	data := struct {
		Hbaddrs1        []HostInfo `json:"hbaddrs1"`
		Hbaddrs2        []HostInfo `json:"hbaddrs2"`
		Hbaddrs2Enabled int        `json:"hbaddrs2_enabled"`
	}{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
