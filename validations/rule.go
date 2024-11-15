/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 		http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: bizhiyuan
 * Date: 2024-05-16 16:23:42
 * LastEditTime: 2024-05-21 14:06:55
 * Description:input validation of Rule
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
