/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/gettext-go"

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
		result["info"] = gettext.Gettext("Action on resource success")
	}

	rac.Data["json"] = &result
	rac.ServeJSON()
}

type ResourceMetaAttributesController struct {
	web.Controller
}

func (rc *ResourceMetaAttributesController) Get() {
	category := rc.Ctx.Input.Param(":category")
	result := models.GetResourceMetaAttributes(category)
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
		result["error"] = gettext.Gettext("invalid input data")
	} else {
		err = models.UpdateResourceAttributes(rscID, reqData)
		if err != nil {
			result["action"] = false
			result["error"] = err.Error()
		} else {
			result["action"] = true
			result["info"] = gettext.Gettext("Update resource attributes success")
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
