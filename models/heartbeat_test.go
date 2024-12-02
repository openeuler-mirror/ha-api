/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Thu Jan 14 13:33:38 2021 +0800
 */
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
