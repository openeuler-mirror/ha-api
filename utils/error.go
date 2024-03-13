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
 * Date: 2024-03-12 15:54:56
 * LastEditTime: 2024-03-12 15:55:19
 * Description: utils 层进行错误处理响应
 */
package utils

type GeneralResponse struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
	Info   string `json:"info,omitempty"`
}

func HandleCmdError(errMsg string, action bool) GeneralResponse {
	return GeneralResponse{
		Action: action,
		Error:  errMsg,
	}
}

func HandleXmlError(errMsg string, action bool) GeneralResponse {
	return GeneralResponse{
		Action: action,
		Error:  errMsg,
	}
}
