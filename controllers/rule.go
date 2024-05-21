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
 * LastEditTime: 2024-05-21 16:23:55
 * Description:规则Controler层
 */

package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"gitee.com/openeuler/ha-api/models"
	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
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

	result := utils.GeneralResponse{}
	requestInput := new(validations.RuleS)

	err := validations.UnmarshalAndValidation(rc.Ctx.Input.RequestBody, requestInput)
	if err != nil {
		rc.handleJsonError(err.Error(), false)
		return
	}

	result = models.RuleAdd(requestInput)
	rc.Data["json"] = &result
	rc.ServeJSON()
}

func (rc *RuleController) Delete() {
	logs.Debug("handle rule DELETE request")
	requestInput := new(validations.DeleteRuleS)
	if err := validations.UnmarshalAndValidation(rc.Ctx.Input.RequestBody, requestInput); err != nil {
		rc.handleJsonError(err.Error(), false)
		return
	}

	res := models.RulesDelete(requestInput)
	rc.Data["json"] = &res
	rc.ServeJSON()
}

func (rc *RuleController) Put() {
	logs.Debug("handle rule PUT request")
	requestInput := new(validations.RuleS)
	if err := validations.UnmarshalAndValidation(rc.Ctx.Input.RequestBody, requestInput); err != nil {
		rc.handleJsonError(err.Error(), false)
		return
	}

	res := models.RuleUpdate(requestInput)
	rc.Data["json"] = &res
	rc.ServeJSON()
}
