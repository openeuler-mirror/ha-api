/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Mon Mar 25 17:23:29 2024 +0800
 */
package utils

type GeneralResponse struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
	Info   string `json:"info,omitempty"`
}

type Response struct {
	Action bool        `json:"action"`
	Error  interface{} `json:"error,omitempty"`
	Info   string      `json:"info,omitempty"`
}
