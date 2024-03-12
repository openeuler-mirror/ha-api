/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-06 16:23:42
 * LastEditTime: 2024-03-06 16:23:55
 * Description:规则Controler层
 */

package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"gitee.com/openeuler/ha-api/models"
)

type RuleController struct {
	web.Controller
}

func (rc *RuleController) Get() {
	logs.Debug("handle rule GET request")
	rscName := rc.GetString("rscname")
	rc.Data["json"] = models.RulesGet(rscName)
	rc.ServeJSON()
}

func (rc *RuleController) Post() {
	logs.Debug("handle rule POST request")
	data := make(map[string]string)
	if err := json.Unmarshal(rc.Ctx.Input.RequestBody, &data); err != nil {
		rc.handleJsonError(err.Error(), false)
	}
	res := models.RuleAdd(data)
	rc.Data["json"] = &res
	rc.ServeJSON()
}
