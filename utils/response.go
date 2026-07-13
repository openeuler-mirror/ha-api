/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

type GeneralResponse struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
	Info   string `json:"info,omitempty"`
}

type RecoveryResponse struct {
	Action bool     `json:"action"`
	Error  string   `json:"error,omitempty"`
	Stack  []string `json:"stack,omitempty"`
}
