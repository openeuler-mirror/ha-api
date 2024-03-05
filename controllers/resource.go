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
 * LastEditTime: 2022-04-19 17:37:45
 * Description: 资源控制器
 ******************************************************************************/
package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"gitee.com/openeuler/ha-api/models"
)

type ResourceController struct {
	web.Controller
}

func (rc *ResourceController) Get() {
	logs.Debug("handle resource GET request")

	rc.Data["json"] = models.GetResourceInfo()
	rc.ServeJSON()
}

func (rc *ResourceController) Post() {
	logs.Debug("handle resource POST request")

	jsonStr := rc.Ctx.Input.RequestBody
	rc.Data["json"] = models.CreateResource(jsonStr)
	rc.ServeJSON()
}

func (rc *ResourceController) Put() {
	logs.Debug("handle resource Put request")
	logs.Debug("do nothing")

	rc.ServeJSON()
}

type ResourceActionController struct {
	web.Controller
}

func (rac *ResourceActionController) Put() {
	logs.Debug("handle resource action Put request")

	rscID := rac.Ctx.Input.Param(":rscID")
	action := rac.Ctx.Input.Param(":action")
	data := rac.Ctx.Input.RequestBody

	result := map[string]interface{}{}
	if err := models.ResourceAction(rscID, action, data); err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["info"] = "Action on resource success"
	}

	rac.Data["json"] = &result
	rac.ServeJSON()
}

type ResourceMetaAttributesController struct {
	web.Controller
}

func (rc *ResourceMetaAttributesController) Get() {
	catagory := rc.Ctx.Input.Param(":catagory")
	result := models.GetResourceMetaAttributes(catagory)
	rc.Data["json"] = &result
	rc.ServeJSON()
}

type ResourceOpsById struct {
	web.Controller
}

func (robi *ResourceOpsById) Get() {
	result := map[string]interface{}{}
	rscID := robi.Ctx.Input.Param(":rscID")
	rst, err := models.GetResourceInfoByrscID(rscID)
	if err != nil {
		result["action"] = false
		result["err"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = rst
	}
	robi.Data["json"] = &result
	robi.ServeJSON()
}

func (robi *ResourceOpsById) Put() {
	rscID := robi.Ctx.Input.Param(":rscID")
	result := map[string]interface{}{}
	reqData := make(map[string]interface{})
	if err := json.Unmarshal(robi.Ctx.Input.RequestBody, &reqData); err != nil {
		result = make(map[string]interface{})
		result["action"] = false
		result["error"] = "invalid input data"
	} else {
		err = models.UpdateResourceAttributes(rscID, reqData)
		if err != nil {
			result["action"] = false
			result["error"] = err.Error()
		} else {
			result["action"] = true
			result["info"] = "Update resource attributes Success"
		}
	}

	robi.Data["json"] = &result
	robi.ServeJSON()
}

type ResourceRelationsController struct {
	web.Controller
}

func (rrc *ResourceRelationsController) Get() {
	logs.Debug("handle resource relation GET request")
	rscID := rrc.Ctx.Input.Param(":rscID")
	relation := rrc.Ctx.Input.Param(":relation")

	result := map[string]interface{}{}

	retData, err := models.GetResourceConstraints(rscID, relation)
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
	} else {
		result["action"] = true
		result["data"] = retData
	}

	rrc.Data["json"] = &result
	rrc.ServeJSON()
}
