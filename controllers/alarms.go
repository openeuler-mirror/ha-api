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
 * LastEditTime: 2022-04-19 17:37:59
 * Description: 告警控制器
 ******************************************************************************/
package controllers

import (
	"github.com/beego/beego/v2/server/web"

	"encoding/json"

	"gitee.com/openeuler/ha-api/models"
)

type AlarmConfig struct {
	web.Controller
}

func (ac *AlarmConfig) Get() {
	ac.Data["json"] = models.AlarmsGet()
	ac.ServeJSON()
}

func (ac *AlarmConfig) Post() {
	var result map[string]interface{}

	reqData := make(map[string]string)
	if err := json.Unmarshal(ac.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		result = models.AlarmsSet(reqData)
	}
	ac.Data["json"] = &result
	ac.ServeJSON()
}
