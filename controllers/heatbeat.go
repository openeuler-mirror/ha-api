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
 * LastEditTime: 2022-04-19 17:37:55
 * Description: 磁盘心跳控制器
 ******************************************************************************/
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type HeartBeatController struct {
	web.Controller
}

func (hbc *HeartBeatController) Get() {
	result := map[string]interface{}{}

	data, err := models.GetHeartBeatConfig()
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = data
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

func (hbc *HeartBeatController) Post() {
	result := map[string]interface{}{}

	data := hbc.Ctx.Input.RequestBody
	if err := models.EditHeartbeatInfo(data); err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = "Change cluster success"
	}

	hbc.Data["json"] = &result
	hbc.ServeJSON()
}

type HeartBeatStatusController struct {
	web.Controller
}

func (hbsc *HeartBeatStatusController) Get() {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = 0

	hbsc.Data["json"] = &result
	hbsc.ServeJSON()
}

type DiskHeartBeatController struct {
	web.Controller
}

func (dhbc *DiskHeartBeatController) Get() {
	result := map[string]interface{}{}

	result["action"] = true
	result["data"] = []string{}

	dhbc.Data["json"] = &result
	dhbc.ServeJSON()
}
