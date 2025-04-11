/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liqiuyu <liqiuyu@kylinos.cn>
 * Date: Mon Jan 18 11:44:18 2021 +0800
 */
package models

import (
	"fmt"
	"reflect"
	"strconv"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

// 报警信息展示
/*
{"action": true, "data": {"sender": "hatest@cs2c.com.cn", "smtp": "mailgw.cs2c.com.cn", "flag": "on", "receiver": ["hatest@cs2c.com.cn", ".com"], "password": "hhaest", "port": "25"}}
*/
type AlarmResponse struct {
	Action bool      `json:"action"`
	Data   AlarmData `json:"data"`
}

type AlarmData struct {
	Flag      bool     `json:"flag"`
	Password  string   `json:"password"`
	Port      float64  `json:"port"`
	Receiver  []string `json:"receiver"` // 允许 null 值
	Sender    string   `json:"sender"`
	Smtp      string   `json:"smtp"`
	Recipient []string `json:"recipient"` // 允许 null 值
}

func AlarmsGet() AlarmResponse {
	var result AlarmResponse
	var doc *etree.Document
	var alertsJson []*etree.Element
	sender := ""
	smtp := ""
	password := ""
	switCh := ""

	out, err := utils.RunCommand(utils.CmdCibQueryConfig)
	if err != nil {
		logs.Error("get alert message failed", err)
		goto ret
	}
	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		logs.Error("parse xml config error", err)
		goto ret
	}

	if len(alertsJson) == 0 {
		alertsJson = doc.FindElements("/configuration/alerts/alert/instance_attributes/nvpair")

		for _, v := range alertsJson {
			fmt.Println(v)
			if v.SelectAttr("name").Value == "email_sender" {
				sender = v.SelectAttr("value").Value
				fmt.Println(sender)
			}
			if v.SelectAttr("name").Value == "email_server" {
				smtp = v.SelectAttr("value").Value
			}
			if v.SelectAttr("name").Value == "password" {
				password = v.SelectAttr("value").Value
			}
			if v.SelectAttr("name").Value == "port" {
				portV := v.SelectAttr("value").Value
				// if portFloat, ok := portV.(float64); ok {

				port, _ := strconv.ParseFloat(portV, 64)
				result.Data.Port = port
				// }
			}
			if v.SelectAttr("name").Value == "switCh" {
				switCh = v.SelectAttr("value").Value
			}
		}
		if switCh == "on" {
			result.Data.Flag = true
		} else {
			result.Data.Flag = false
		}
		var recipients []string
		var vaLue string
		receivers := doc.FindElements("/configuration/alerts/alert/recipient")
		fmt.Println("get receivers: ", receivers)
		for _, v := range receivers {
			fmt.Println(v)
			vaLue = v.SelectAttr("value").Value
			recipients = append(recipients, vaLue)
		}

		cmdStr := "/usr/bin/pwd_decode " + string(password)
		out, _ := utils.RunCommand(cmdStr)
		fmt.Println(password)
		mailPassword := ""
		if string(out) != "the parameter is less\n" {
			mailPassword = string(out)
		}
		result.Data.Sender = sender
		result.Data.Smtp = smtp
		result.Data.Receiver = recipients
		result.Data.Password = mailPassword
	}

ret:
	result.Action = true
	return result
}

func isDataEmpty(data AlarmData) bool {
	zeroValue := AlarmData{}
	return reflect.DeepEqual(data, zeroValue)
}
func AlarmsSet(data AlarmData) map[string]interface{} {
	result := make(map[string]interface{})

	utils.RunCommand(utils.CmdDeleteAlert)
	switCh := ""

	if data.Flag != true {
		switCh = "on"
	} else {
		switCh = "off"
	}

	fmt.Println("get receiver: ", data.Receiver)
	port := strconv.Itoa(int(data.Port))
	opsStr := " options email_sender=" + data.Sender + " email_server=" + data.Smtp + " password=" + data.Password + " port=" + port + " switCh=" + switCh
	cmdStr := utils.CmdCreateAlert + opsStr
	fmt.Println(cmdStr)
	_, err := utils.RunCommand(cmdStr)
	if err != nil {
		result["action"] = false
		fmt.Println("err1")
		fmt.Println(err)
		result["error"] = gettext.Gettext("Set alarm failed")
		return result
	}

	for _, recipient := range data.Receiver {
		fmt.Println("set recipient")
		reveiverStr := utils.CmdAddAlert + " value=" + string(recipient) + " --force"
		fmt.Println(reveiverStr)
		_, err := utils.RunCommand(reveiverStr)
		if err != nil {
			result["action"] = false
			fmt.Println("err2")
			fmt.Println(err)
			result["error"] = gettext.Gettext("Set alarm failed")
			return result
		}
	}

	result["action"] = true
	result["info"] = gettext.Gettext("Set alarm success")
	return result
}

func AlarmsTest() utils.GeneralResponse {
	var result utils.GeneralResponse
	alarmConfig := AlarmsGet().Data
	for _, recipient := range alarmConfig.Receiver {
		port := strconv.Itoa(int(alarmConfig.Port))
		reveiverStr := utils.CmdSendEmail + alarmConfig.Smtp + "' '" + alarmConfig.Sender + "' '" + alarmConfig.Password + "' '" + recipient + "' '此邮件为测试邮件' " + port
		fmt.Println(reveiverStr)
		out, err := utils.RunCommand(reveiverStr)
		if err != nil {
			fmt.Println(out)
			fmt.Println(err)
			testcmd := "echo send mail to " + string(recipient) + " failed:" + string(out) + " >>/var/log/mailtest.log"
			out1, err1 := utils.RunCommand(testcmd)
			if err1 != nil {
				fmt.Println(out1)
			}
			result.Action = false
			result.Error = gettext.Gettext("Send alarm test failed")
			return result

		} else {
			testcmd := "echo send mail to " + string(recipient) + " success >>/var/log/mailtest.log"
			utils.RunCommand(testcmd)
		}
	}
	result.Action = true
	result.Info = gettext.Gettext("Send alarm test success")
	return result

}
