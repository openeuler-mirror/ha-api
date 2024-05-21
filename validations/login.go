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
 * Date: 2024-05-16 16:23:42
 * LastEditTime: 2024-05-21 14:06:55
 * Description:input validation of Login
 */
package validations

import (
	"gitee.com/openeuler/ha-api/settings"
	"github.com/beego/beego/v2/core/validation"
	"github.com/chai2010/gettext-go"
)

type PasswordS struct {
	Password string `json:"password" valid:"Required;"`
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
