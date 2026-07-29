/*
 * Copyright (c) KylinSoft  Co., Ltd. 2026.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Fri Jul 04 15:54:28 2026 +0800
 */

 package controllers

 import (
	 "gitee.com/openeuler/ha-api/models"
	 "github.com/beego/beego/v2/server/web"
 )
 
 type ResourceModelConfig struct {
	 web.Controller
 }
 
 func (rmc *ResourceModelConfig) Get() {
	 rmc.Data["json"] = models.ResourceModelDeploy()
	 rmc.ServeJSON()
 }
 