/*
 * @Author: bixiaoyan bixiaoyan@kylinos.cn
 * @Date: 2024-03-13 15:04:41
 * @LastEditors: bixiaoyan bixiaoyan@kylinos.cn
 * @LastEditTime: 2024-06-17 09:22:06
 * @FilePath: /ha-api/controllers/tag.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
 * Date: 2024-03-13 15:05:51
 * LastEditTime: 2024-03-13 18:37:48
 * Description: tag
 ******************************************************************************/
package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"gitee.com/openeuler/ha-api/models"
	"github.com/beego/beego/v2/core/logs"
)

type TagController struct {
	web.Controller
}

func (tc *TagController) Get() {
	tc.Data["json"] = models.GetTag()
	tc.ServeJSON() 
}

func (tc *TagController) Post() {
	logs.Debug("handle Tag POST request")
	jsonStr := tc.Ctx.Input.RequestBody
	tc.Data["json"] = models.SetTag(jsonStr)
	tc.ServeJSON()
}

type TagUpdateController struct {
	web.Controller
}

func (tuc *TagUpdateController) Put() {
	logs.Debug("handle TagUpdate POST request")
	jsonStr := tuc.Ctx.Input.RequestBody
	tagName := tuc.Ctx.Input.Param(":tag_name")
	tuc.Data["json"] = models.UpdateTag(tagName, jsonStr)
	tuc.ServeJSON()
}

type TagActionController struct {
	web.Controller
}

func (tuc *TagActionController) Put() {
	tagName := tuc.Ctx.Input.Param(":tag_name")
	atcion := tuc.Ctx.Input.Param(":action") 
	tuc.Data["json"] = models.TagAction(tagName,atcion)
	tuc.ServeJSON()
}



 
  