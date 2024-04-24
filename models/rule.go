/*******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
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
		goto ret
	}
	if len(rules) == 0 {
		rules := doc.FindElements("//rsc_location")
		for _, ruleElem := range rules {
			if rsc := ruleElem.SelectAttrValue("rsc", ""); rsc == rscName {
				rule := ruleElem.SelectElement("rule")
				if rule != nil {
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
func RulesDelete(ruleids map[string][]string) RuleDeleteResponse {
	ruleIdList := ruleids["ids"]
	var res []map[string]interface{}

	for _, id := range ruleIdList {
		delRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, id)
		_, err := utils.RunCommand(delRuleCmd)
		if err != nil {
			res = append(res, map[string]interface{}{
				"id":    id,
				"error": "编号为 " + id + "的规则找不到", // error信息特殊处理？
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

func RuleAdd(data map[string]string) utils.GeneralResponse {
	var cmdAddRule string
	if _, ok := data["ruleid"]; ok {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAddWithId, data["rsc"], data["score"], data["ruleid"])
	} else {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAdd, data["rsc"], data["score"])
	}

	if _, ok := data["value"]; ok {
		cmdAddRule = cmdAddRule + " " + data["attribute"] + " " + data["operation"] + " " + data["value"]
	} else {
		cmdAddRule = cmdAddRule + " " + data["operation"] + " " + data["attribute"]
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

func RuleUpdate(data map[string]string) utils.GeneralResponse {
	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return utils.HandleCmdError(err.Error(), false)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(out); err != nil {
		return utils.HandleXmlError(err.Error(), false)
	}
	r := make(map[string]string)
	rules := doc.FindElements("//rsc_location/rule")
	for _, ruleElem := range rules {
		if ruleId := ruleElem.SelectAttrValue("id", ""); ruleId == data["ruleid"] {
			r["ruleid"] = ruleElem.SelectAttrValue("id", "")
			r["score"] = ruleElem.SelectAttrValue("score", "")
			r["attribute"] = ruleElem.SelectElement("expression").SelectAttrValue("attribute", "")
			r["operation"] = ruleElem.SelectElement("expression").SelectAttrValue("operation", "")
			r["value"] = ruleElem.SelectElement("expression").SelectAttrValue("value", "")
		}
	}
	// delete old rule
	deleteRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, data["ruleid"])
	_, err = utils.RunCommand(deleteRuleCmd)
	if err != nil {
		return utils.HandleCmdError("更新规则失败，原id为"+data["ruleid"]+"规则找不到", false)
	}
	// add new rule
	resp := RuleAdd(data)
	if !resp.Action {
		// recovery the update op
		RuleAdd(r)
		return utils.HandleCmdError(gettext.Gettext("Update rule failed, duplicate constraint already exists"), false)
	}

	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Update rule success"),
	}

}
