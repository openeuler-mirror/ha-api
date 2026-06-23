/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package validations

import (
	"gitee.com/openeuler/ha-api/settings"
	"github.com/beego/beego/v2/core/validation"
	"github.com/chai2010/gettext-go"
)

type PasswordS struct {
	OldPwd string `json:"old_password" valid:"Required;"`
	NewPwd string `json:"new_password" valid:"Required;"`
}
type UserS struct {
	UserName string `json:"username" valid:"Required;"`
	Password string `json:"password" valid:"Required;"`
}

func (u *UserS) Valid(v *validation.Validation) {
	if u.UserName != settings.PacemakerUname {
		v.SetError("UserName", gettext.Gettext("username is not allowed to login"))
	}
}
