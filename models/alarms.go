/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 10:05:59
 * Description: 告警功能
 ******************************************************************************/
package models

import (
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

// 报警信息展示
/*
{"action": true, "data": {"sender": "hatest@cs2c.com.cn", "smtp": "mailgw.cs2c.com.cn", "flag": "on", "receiver": ["hatest@cs2c.com.cn", ".com"], "password": "hhaest", "port": "25"}}
*/
func AlarmsGet() map[string]interface{} {
	var value map[string]interface{}
	var doc *etree.Document
	var alertsJson []*etree.Element
	sender := ""
	smtp := ""
	password := ""
	port := ""
	switCh := ""

	out, err := utils.RunCommand(utils.CmdCibQueryConfig)

	if err != nil {
		logs.Error("get alert message failed", err)
		goto ret
	}
	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		logs.Error("parse xml comfig error", err)
		goto ret
	}

	if len(alertsJson) == 0 {
		alertsJson := doc.FindElements("/configuration/crm_config/cluster_property_set/nvpair")
		for _, v := range alertsJson {
			if value[v.SelectAttr("name").Value] == "emailSender" {
				sender = v.SelectAttr("value").Value
			}
			if value[v.SelectAttr("name").Value] == "emailServer" {
				smtp = v.SelectAttr("value").Value
			}
			if value[v.SelectAttr("name").Value] == "password" {
				password = v.SelectAttr("value").Value
			}
			if value[v.SelectAttr("name").Value] == "port" {
				port = v.SelectAttr("value").Value
			}
			if value[v.SelectAttr("name").Value] == "switCh" {
				switCh = v.SelectAttr("value").Value
			}
		}
		var recipients []string
		var vaLue string
		receivers := doc.SelectElements("recipient")
		for _, v := range receivers {
			vaLue = v.SelectAttr("value").Value
			recipients = append(recipients, vaLue)
		}

		cmdStr := "/usr/bin/pwd_decode" + string(password)
		out, _ := utils.RunCommand(cmdStr)
		mailPassword := string(out)

		data := map[string]interface{}{}
		data["sender"] = sender
		data["smtp"] = smtp
		data["lag"] = switCh
		data["receiver"] = recipients
		data["password"] = mailPassword
		data["port"] = port
	}

ret:
	data := map[string]interface{}{}
	result := make(map[string]interface{})
	if len(data) == 0 {
		data["flag"] = false
		data["smtp"] = ""
		data["port"] = ""
		data["sender"] = ""
		data["password"] = ""
		data["receiver"] = []string{}
	}

	result["action"] = true
	result["data"] = data
	return result
}

func AlarmsSet(data map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	var receiver string

	utils.RunCommand(utils.CmdDeleteAlert)
	sender := ""
	smtp := ""
	password := ""
	port := ""
	switCh := ""

	if len(data) != 0 {
		if _, ok := data["flag"]; ok {
			if data["flag"] != "" {
				switCh = "on"
			} else {
				switCh = "off"
			}
		}
		if _, ok := data["smtp"]; ok {
			smtp = string(data["smtp"])
		}
		if _, ok := data["sender"]; ok {
			sender = string(data["sender"])
		}
		if _, ok := data["password"]; ok {
			password = string(data["password"])
		}
		if _, ok := data["port"]; ok {
			port = string(data["port"])
		}
		if _, ok := data["receiver"]; ok {
			receiver = data["receiver"]
		}
	}

	opsStr := " options email_sender=" + sender + " email_server=" + smtp + " password=" + password + " port=" + port + " switCh=" + switCh
	cmdStr := utils.CmdCreateAlert + opsStr
	_, err := utils.RunCommand(cmdStr)
	if err != nil {
		result["action"] = false
		result["error"] = gettext.Gettext("Set alarm failed")
		return result
	}

	for _, recipient := range receiver {
		reveiverStr := utils.CmdAddAlert + " value=" + string(recipient) + " --force"
		_, err := utils.RunCommand(reveiverStr)
		if err != nil {
			result["action"] = false
			result["error"] = gettext.Gettext("Set alarm failed")
			return result
		}
	}

	result["action"] = true
	result["info"] = gettext.Gettext("Set alarm success")
	return result
}
