/*******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-03-06 10:42:50
 * LastEditTime: 2024-03-06 11:23:09
 * Description:规则
 ******************************************************************************/
package models

import (
	"fmt"

	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/etree"
	"github.com/chai2010/gettext-go"
)

type Rule struct {
	Rsc       string `json:"rsc"`
	RuleId    string `json:"ruleid"`
	Score     string `json:"score"`
	Attribute string `json:"attribute"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type RuleGetResponse struct {
	utils.GeneralResponse
	Data []Rule `json:"data"`
}

func RulesGet(rscName string) RuleGetResponse {

	var doc *etree.Document
	var rules []*etree.Element

	rulelist := []Rule{}

	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		goto ret
	}

	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		logs.Error("RulesGet: Read xml  : %v", err)
		goto ret
	}

	rules = doc.FindElements("//rsc_location")
	for _, ruleElem := range rules {
		if rsc := ruleElem.SelectAttrValue("rsc", ""); rsc == rscName {
			rule := ruleElem.SelectElement("rule")
			if rule != nil {
				// 降级返回
				r := Rule{
					Rsc:       rsc,
					RuleId:    rule.SelectAttrValue("id", ""),
					Score:     rule.SelectAttrValue("score", ""),
					Attribute: rule.SelectElement("expression").SelectAttrValue("attribute", ""),
					Operation: rule.SelectElement("expression").SelectAttrValue("operation", ""),
					Value:     rule.SelectElement("expression").SelectAttrValue("value", ""),
				}
				rulelist = append(rulelist, r)
			}

		}
	}

	return RuleGetResponse{
		Data: rulelist,
		GeneralResponse: utils.GeneralResponse{
			Action: true,
		},
	}

ret:
	return RuleGetResponse{
		GeneralResponse: utils.GeneralResponse{
			Action: false,
			Error:  gettext.Gettext("Get rule failed"),
		},
	}
}

type RuleDeleteResponse struct {
	Action bool                     `json:"action"`
	Info   string                   `json:"info,omitempty"`
	Error  []map[string]interface{} `json:"error,omitempty"`
}

// Todo:return list ?
func RulesDelete(ruleids *validations.DeleteRuleS) RuleDeleteResponse {
	ruleIdList := ruleids.RuleIDs
	var res []map[string]interface{}

	for _, id := range ruleIdList {
		delRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, id)
		_, err := utils.RunCommand(delRuleCmd)
		if err != nil {
			res = append(res, map[string]interface{}{
				"id":    id,
				"error": fmt.Sprintf(gettext.Gettext("Rule %s could not be found"), id),
			})
		}
	}
	if len(res) != 0 {
		// Todo： some rules was deleted failed, recovery？
		return RuleDeleteResponse{
			Action: false,
			Error:  res,
		}
	}
	return RuleDeleteResponse{
		Action: true,
		Info:   gettext.Gettext("Delete rule success"),
	}
}

func RuleAdd(data *validations.RuleS) utils.GeneralResponse {
	var cmdAddRule string
	if data.RuleID != "" {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAddWithId, data.Rsc, data.Score, data.RuleID)
	} else {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAdd, data.Rsc, data.Score)
	}

	if data.Value != "" {
		cmdAddRule = cmdAddRule + " " + data.Attribute + " " + data.Operation + " " + data.Value
	} else {
		cmdAddRule = cmdAddRule + " " + data.Operation + " " + data.Attribute
	}
	_, err := utils.RunCommand(cmdAddRule)
	if err != nil {
		return utils.HandleCmdError(gettext.Gettext("Add rule failed, duplicate constraint already exists"), false)
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Add rule success"),
	}
}

func RuleUpdate(data *validations.RuleS) utils.GeneralResponse {
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return utils.HandleCmdError(err.Error(), false)
	}

	oldRule := new(validations.RuleS)
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(out); err != nil {
		return utils.HandleXmlError(err.Error(), false)
	}
	rules := doc.FindElements("//rsc_location/rule")
	for _, ruleElem := range rules {
		ruleId := ruleElem.SelectAttrValue("id", "")
		fmt.Println("ruleId: ", ruleId)
		if ruleId == data.RuleID {
			oldRule.RuleID = ruleElem.SelectAttrValue("id", "")
			oldRule.Score = ruleElem.SelectAttrValue("score", "")
			oldRule.Attribute = ruleElem.SelectElement("expression").SelectAttrValue("attribute", "")
			oldRule.Operation = ruleElem.SelectElement("expression").SelectAttrValue("operation", "")
			oldRule.Value = ruleElem.SelectElement("expression").SelectAttrValue("value", "")
		}
	}
	// delete old rule
	deleteRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, data.RuleID)
	fmt.Println(deleteRuleCmd)
	_, err = utils.RunCommand(deleteRuleCmd)
	if err != nil {
		return utils.HandleCmdError(
			fmt.Sprintf(gettext.Gettext("Update rule failed, rule %s not found"), data.RuleID), false)
	}

	// add new rule
	resp := RuleAdd(data)
	if !resp.Action {
		// recovery the update op
		RuleAdd(oldRule)
		return utils.HandleCmdError(gettext.Gettext("Update rule failed, duplicate constraint already exists"), false)
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Update rule success"),
	}
}
