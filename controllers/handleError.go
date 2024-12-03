/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Tue Mar 12 16:14:04 2024 +0800
 */
package controllers

type ErrorInfo struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
}

func (rc *RuleController) handleJsonError(errMsg string, action bool) {
	result := ErrorInfo{
		Action: action,
		Error:  errMsg,
	}
	rc.Data["json"] = &result
	rc.ServeJSON()
}
