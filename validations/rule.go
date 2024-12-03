/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Tue May 21 17:24:27 2024 +0800
 */
package validations

import (
	"gitee.com/openeuler/ha-api/utils"
	"github.com/beego/beego/v2/core/validation"
)

type RuleS struct {
	Attribute string `json:"attribute" valid:"Required"`
	Operation string `json:"operation" valid:"Required"`
	Rsc       string `json:"rsc" valid:"Required"`
	Score     string `json:"score" valid:"Required"`
	Value     string `json:"value,omitempty"`
	RuleID    string `json:"ruleid,omitempty"`
}

type DeleteRuleS struct {
	RuleIDs []string `json:"ids,omitempty"`
}

func (r *RuleS) Valid(v *validation.Validation) {
	operationsRange := []string{"gt", "lt", "gte", "lte", "eq", "ne", "defined", "not_defined"}
	if !utils.Contains(operationsRange, r.Operation) {
		v.SetError("Operation", "operation input should in [ gt lte gte le eq ne defined not_defined ]")
	}
}

func (r *DeleteRuleS) Valid(v *validation.Validation) {
	ruleIDs := r.RuleIDs
	if len(ruleIDs) == 0 {
		v.SetError("RuleIDs", "ids to been deleted is empty")
	}
	for _, ruleId := range ruleIDs {
		if ruleId == "" {
			v.SetError("RuleIDs", "some ruleId in RuleIDs is \"\"")
			break
		}
	}
}
