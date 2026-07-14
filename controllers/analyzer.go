/*
 * Copyright (c) KylinSoft  Co., Ltd. 2027.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: xuxiaojuan <xuxiaojuan@kylinos.cn>
 * Date: Wed July 8 13:56:40 2026 +0800
 */
package controllers

import (
        "gitee.com/openeuler/ha-api/models"
        "github.com/beego/beego/v2/server/web"
)

type SystemInfoController struct {
        web.Controller
}

func (c *SystemInfoController) Get() {
        sysInfo := models.GetSystemInfo()
        c.Data["json"] = sysInfo
        c.ServeJSON()
}

type ServiceStatusController struct {
        web.Controller
}

func (sst *ServiceStatusController) Get() {
        serviceStatus := models.GetServiceStatus()
        sst.Data["json"] = serviceStatus
        sst.ServeJSON()
}
