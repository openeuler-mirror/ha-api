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
 * Date: 2024-03-07 14:49:51
 * LastEditTime: 2024-03-07 17:37:48
 * Description: 集群利用率
 ******************************************************************************/
package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
)

type UtilizationController struct {
	web.Controller
}

func (ul *UtilizationController) Get() {

	ul.Data["json"] = models.GetUtilization()
	ul.ServeJSON() 
}

func (ul *UtilizationController) Post() {
	logs.Debug("handle utilization POST request")

	jsonStr := ul.Ctx.Input.RequestBody
	ul.Data["json"] = models.SetUtilization(jsonStr)
	ul.ServeJSON()
}

func (ul *UtilizationController) Delete() {
	logs.Debug("handle utilization DELETE request")

	jsonStr := ul.Ctx.Input.RequestBody
	ul.Data["json"] = models.DelUtilization(jsonStr)
	ul.ServeJSON()
}


 