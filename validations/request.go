/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-05-16 16:23:42
 * LastEditTime: 2024-05-21 14:06:55
 * Description:
 */
package validations

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/chai2010/gettext-go"
)

// unmarshal the input object into struct and validate input
func UnmarshalAndValidation(input []byte, data interface{}) error {
	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	valid := validation.Validation{}
	validCheck, _ := valid.Valid(data)
	if !validCheck {
		data := make(map[string]string)
		for _, err := range valid.Errors {
			data[err.Key] = err.Message
			logs.Warn(err.Message)
		}

		return fmt.Errorf(gettext.Gettext("input validation failed"))
	}
	return nil
}
