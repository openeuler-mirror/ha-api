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
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
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
	GeneralResponse
	Data []Rule `json:"data"`
}

type GeneralResponse struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
	Info   string `json:"info,omitempty"`
}

const cstQueryCmd string = "cibadmin --query --scope constraints"

func RulesGet(rscName string) RuleGetResponse {

	var doc *etree.Document
	var rules []*etree.Element

	rulelist := []Rule{}

	out, err := utils.RunCommand(cstQueryCmd)
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
		GeneralResponse: GeneralResponse{
			Action: true,
		},
	}

ret:
	return RuleGetResponse{
		GeneralResponse: GeneralResponse{
			Action: true,
			Error:  err.Error(),
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
		delRuleCmd := "pcs constraint rule delete " + id
		_, err := utils.RunCommand(delRuleCmd)
		if err != nil {
			res = append(res, map[string]interface{}{
				"id":    id,
				"error": err.Error(),
			})
		}
	}
	if len(res) != 0 {
		return RuleDeleteResponse{
			Action: true,
			Error:  res,
		}
	}
	return RuleDeleteResponse{
		Action: true,
		Info:   "Delete rule success!",
	}
}

func RuleAdd(data map[string]string) GeneralResponse {
	var cmdAddRule string
	if _, ok := data["ruleid"]; ok {
		cmdAddRule = "pcs constraint location " + data["rsc"] + " rule score=" + data["score"] + " id=" + data["ruleid"]
	} else {
		cmdAddRule = "pcs constraint location " + data["rsc"] + " rule score=" + data["score"]
	}

	if _, ok := data["value"]; ok {
		cmdAddRule = cmdAddRule + " " + data["attribute"] + " " + data["operation"] + " " + data["value"]
	} else {
		cmdAddRule = cmdAddRule + " " + data["operation"] + " " + data["attribute"]
	}
	_, err := utils.RunCommand(cmdAddRule)
	if err != nil {
		return GeneralResponse{
			Action: false,
			Error:  err.Error(),
		}
	}
	return GeneralResponse{
		Action: true,
		Info:   "Add rule success!",
	}

}
