/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"fmt"
	"log/slog"

	"gitee.com/openeuler/ha-api/settings"
	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
	"github.com/chai2010/gettext-go"
)

// 登录接口
func Login(data validations.UserS, authChecker utils.AuthChecker) map[string]interface{} {
	resp := map[string]interface{}{}
	var crypto utils.CryptoProvider
	var originPasswd string
	// RSA算法解密
	crypto, err := utils.NewCryptoProvider(utils.AlgorithmRSA, map[string]string{
		"publicKeyPath":  settings.RSA_PUBLIC_KEY,
		"privateKeyPath": settings.RSA_PRIVATE_KEY,
	})
	if err != nil {
		slog.Error(err.Error())
		resp["action"] = false
		resp["error"] = gettext.Gettext("Username or password error")
		return resp
	}

	originPasswd, err = crypto.Decrypt(data.Password)
	if err != nil {
		slog.Error(err.Error())
		resp["action"] = false
		resp["error"] = gettext.Gettext("Username or password error")
		return resp
	}

	// check password
	if authChecker.CheckAuth(data.UserName, originPasswd) {
		clusterInfo := CheckIsClusterExist()
		if action, ok := clusterInfo["action"].(bool); ok && action {
			resp["cluster_exist"] = true
			localClusterName, _ := clusterInfo["cluster_name"].(string)
			resp["local_cluster_name"] = localClusterName
		}
		resp["action"] = true
		resp["info"] = gettext.Gettext("Login success")
	} else {
		resp["action"] = false
		resp["error"] = gettext.Gettext("Username or password error")
		slog.Error("Login failed: username or password error")
	}
	return resp
}

func PasswordChange(data validations.PasswordS, authChecker utils.AuthChecker) utils.GeneralResponse {
	// authChecker := new(utils.PamAuthChecker)
	result := new(utils.GeneralResponse)
	var crypto utils.CryptoProvider
	var oldDecryptedPasswd string
	var newDecryptedPasswd string
	var cmd string = ""
	// RSA算法解密
	crypto, err := utils.NewCryptoProvider(utils.AlgorithmRSA, map[string]string{
		"publicKeyPath":  settings.RSA_PUBLIC_KEY,
		"privateKeyPath": settings.RSA_PRIVATE_KEY,
	})
	if err != nil {
		slog.Error(err.Error())
		result.Action = false
		result.Error = gettext.Gettext("Change password failed")
		return *result
	}

	oldDecryptedPasswd, err = crypto.Decrypt(data.OldPwd)

	if err != nil {
		slog.Error(err.Error())
		result.Action = false
		result.Error = gettext.Gettext("Change password failed")
		return *result
	}

	// 校验旧的passwd
	if !authChecker.CheckAuth(settings.PacemakerUname, oldDecryptedPasswd) {
		slog.Error("old password incorrect")
		result.Action = false
		result.Error = gettext.Gettext("Old Password incorrect")
		return *result
	}

	// 新的passwd
	newDecryptedPasswd, err = crypto.Decrypt(data.NewPwd)
	slog.Debug("Passwd Change", "oldPasswd", oldDecryptedPasswd)
	slog.Debug("Passwd Change", "newPasswd", newDecryptedPasswd)
	if err != nil {
		utils.Error(err, err.Error())
		result.Action = false
		result.Error = gettext.Gettext("New Password incorrect")
	}

	cmd = fmt.Sprintf(utils.CmdChangePwd, newDecryptedPasswd)
	if _, err := utils.RunCommand(cmd); err != nil {
		slog.Error(fmt.Sprintf("change password failed: %s", err.Error()))
		result.Action = false
		result.Error = gettext.Gettext("Change password failed")
		return *result
	} else {
		result.Action = true
		result.Info = gettext.Gettext("Change password success")
	}

	return *result
}
