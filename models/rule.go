/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Tue Mar 12 10:13:58 2024 +0800
 */
package models

import (
	"fmt"
	"log/slog"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
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
		slog.Error("RulesGet: Read xml failed", "error", err)
		goto ret
	}

	rules = doc.FindElements("/constraints/rsc_location")
	for _, ruleElem := range rules {
		if rsc := ruleElem.SelectAttrValue("rsc", ""); rsc == rscName {
			rule := ruleElem.SelectElement("rule")
			if rule != nil {
				expr := rule.SelectElement("expression")
				r := Rule{
					Rsc:    rsc,
					RuleId: rule.SelectAttrValue("id", ""),
					Score:  rule.SelectAttrValue("score", ""),
				}
				if expr != nil {
					r.Attribute = expr.SelectAttrValue("attribute", "")
					r.Operation = expr.SelectAttrValue("operation", "")
					r.Value = expr.SelectAttrValue("value", "")
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
	slog.Debug("RulesDelete", "ids", ruleIdList)
	var res []map[string]interface{}

	for _, id := range ruleIdList {
		delRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, utils.ShellEscape(id))
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
	slog.Debug("RuleAdd",
		"id", data.RuleID,
		"value", data.Value,
	)
	var cmdAddRule string
	if data.RuleID != "" {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAddWithId, utils.ShellEscape(data.Rsc), utils.ShellEscape(data.Score), utils.ShellEscape(data.RuleID))
	} else {
		cmdAddRule = fmt.Sprintf(utils.CmdRuleAdd, utils.ShellEscape(data.Rsc), utils.ShellEscape(data.Score))
	}

	if data.Value != "" {
		cmdAddRule = cmdAddRule + " " + utils.ShellEscape(data.Attribute) + " " + utils.ShellEscape(data.Operation) + " " + utils.ShellEscape(data.Value)
	} else {
		cmdAddRule = cmdAddRule + " " + utils.ShellEscape(data.Operation) + " " + utils.ShellEscape(data.Attribute)
	}
	_, err := utils.RunCommand(cmdAddRule)
	if err != nil {
		return utils.HandleCmdError(gettext.Gettext("Add rule failed, duplicate constraint already exists"), false)
	}
	slog.Info("Add rule success", "id", data.RuleID)
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
		if ruleId == data.RuleID {
			if parent := ruleElem.Parent(); parent != nil {
				oldRule.Rsc = parent.SelectAttrValue("rsc", "")
			}
			oldRule.RuleID = ruleElem.SelectAttrValue("id", "")
			oldRule.Score = ruleElem.SelectAttrValue("score", "")
			if expr := ruleElem.SelectElement("expression"); expr != nil {
				oldRule.Attribute = expr.SelectAttrValue("attribute", "")
				oldRule.Operation = expr.SelectAttrValue("operation", "")
				oldRule.Value = expr.SelectAttrValue("value", "")
			}
		}
	}
	// delete old rule
	deleteRuleCmd := fmt.Sprintf(utils.CmdRuleDelete, utils.ShellEscape(data.RuleID))
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
	slog.Info("Update rule success", "id", data.RuleID)
	return utils.GeneralResponse{
		Action: true,
		Info:   gettext.Gettext("Update rule success"),
	}
}

// StaleConstraint 残留约束信息，由 pcs resource move/ban 自动生成的 location 约束
type StaleConstraint struct {
	ConstraintID string `json:"constraint_id"`
	Resource     string `json:"resource"`
	Node         string `json:"node"`
	Score        string `json:"score"`
	Role         string `json:"role,omitempty"`
	Type         string `json:"type"`
}

type StaleConstraintsResponse struct {
	utils.GeneralResponse
	Data []StaleConstraint `json:"data"`
}

// StaleConstraintsGet 检测由 pcs resource move/ban 残留的 location 约束。
// move/ban 自动生成的约束 ID 以 "cli-" 为前缀（如 cli-ban-vip-on-node1），
// 用户手动创建的约束 ID 为自定义命名，不受影响。
func StaleConstraintsGet() StaleConstraintsResponse {
	staleList := []StaleConstraint{}

	out, err := utils.RunCommand(utils.CmdQueryConstraints)
	if err != nil {
		return StaleConstraintsResponse{
			GeneralResponse: utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Get stale constraints failed"),
			},
		}
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(out); err != nil {
		slog.Error("StaleConstraintsGet: Read xml failed", "error", err)
		return StaleConstraintsResponse{
			GeneralResponse: utils.GeneralResponse{
				Action: false,
				Error:  gettext.Gettext("Get stale constraints failed"),
			},
		}
	}

	for _, loc := range doc.FindElements("/constraints/rsc_location") {
		id := loc.SelectAttrValue("id", "")
		if !strings.HasPrefix(id, "cli-") {
			continue
		}
		// 按前缀区分约束来源：ban 残留 / move 残留
		cType := "move"
		switch {
		case strings.HasPrefix(id, "cli-ban-"):
			cType = "ban"
		case strings.HasPrefix(id, "cli-standby-"), strings.HasPrefix(id, "cli-prefer-"):
			cType = "move"
		}
		staleList = append(staleList, StaleConstraint{
			ConstraintID: id,
			Resource:     loc.SelectAttrValue("rsc", ""),
			Node:         loc.SelectAttrValue("node", ""),
			Score:        loc.SelectAttrValue("score", ""),
			Role:         loc.SelectAttrValue("role", ""),
			Type:         cType,
		})
	}

	return StaleConstraintsResponse{
		GeneralResponse: utils.GeneralResponse{
			Action: true,
		},
		Data: staleList,
	}
}
