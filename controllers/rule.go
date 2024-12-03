/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Tue Mar 12 10:13:58 2024 +0800
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
