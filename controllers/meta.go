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
 * LastEditTime: 2022-04-19 17:37:48
 * Description: 元控制器
 ******************************************************************************/
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"openkylin.com/ha-api/models"
)

type MetaController struct {
	web.Controller
}

func (mc *MetaController) Get() {
	rscClass := mc.Ctx.Input.Param(":rsc_class")
	rscType := mc.Ctx.Input.Param(":rsc_type")
	rscProvider := mc.Ctx.Input.Param(":rsc_provider")
	mc.Data["json"] = models.GetResourceMetas(rscClass, rscType, rscProvider)
	mc.ServeJSON()
}

type MetasController struct {
	web.Controller
}

func (mc *MetasController) Get() {
	mc.Data["json"] = models.GetAllResourceMetas()
	mc.ServeJSON()
}
