/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Tue May 21 16:25:25 2024 +0800
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
