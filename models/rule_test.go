/*
 * Copyright (c) KylinSoft  Co., Ltd. 2027.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: xuxiaojuan <xuxiaojuan@kylinos.cn>
 * Date: Wed July 8 13:56:40 2026 +0800
 */

package models

import (
	"errors"
	"strings"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"gitee.com/openeuler/ha-api/validations"
	"github.com/stretchr/testify/assert"
)

func ruleMockCmd(t *testing.T, fn func(string) ([]byte, error)) {
	t.Helper()
	orig := utils.RunCommand
	utils.RunCommand = fn
	t.Cleanup(func() { utils.RunCommand = orig })
}

// ==================== RulesGet ====================

func TestRulesGet_Success(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		xml := `<constraints>
  <rsc_location id="location-dummy-node1" rsc="dummy" node="node1" score="INFINITY">
    <rule id="location-dummy-node1-rule" score="INFINITY">
      <expression id="expr-1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
  <rsc_location id="location-vip-node2" rsc="vip" node="node2" score="100">
    <rule id="location-vip-node2-rule" score="100">
      <expression id="expr-2" attribute="#uname" operation="eq" value="node2"/>
    </rule>
  </rsc_location>
</constraints>`
		return []byte(xml), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "dummy", result.Data[0].Rsc)
	assert.Equal(t, "location-dummy-node1-rule", result.Data[0].RuleId)
	assert.Equal(t, "INFINITY", result.Data[0].Score)
	assert.Equal(t, "#uname", result.Data[0].Attribute)
	assert.Equal(t, "eq", result.Data[0].Operation)
	assert.Equal(t, "node1", result.Data[0].Value)
}

func TestRulesGet_MultipleRulesForSameResource(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		xml := `<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="rule-1" score="100">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
  <rsc_location id="loc-2" rsc="dummy">
    <rule id="rule-2" score="-INFINITY">
      <expression id="e2" attribute="#uname" operation="eq" value="node2"/>
    </rule>
  </rsc_location>
</constraints>`
		return []byte(xml), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "rule-1", result.Data[0].RuleId)
	assert.Equal(t, "100", result.Data[0].Score)
	assert.Equal(t, "rule-2", result.Data[1].RuleId)
	assert.Equal(t, "-INFINITY", result.Data[1].Score)
}

func TestRulesGet_NoMatchingResource(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		xml := `<constraints>
  <rsc_location id="loc-1" rsc="vip">
    <rule id="rule-vip" score="INFINITY">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
</constraints>`
		return []byte(xml), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Empty(t, result.Data)
}

func TestRulesGet_EmptyConstraints(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte(`<constraints></constraints>`), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Empty(t, result.Data)
}

func TestRulesGet_CommandError(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte("Error: cibadmin failed"), errors.New("command failed")
	})

	result := RulesGet("dummy")

	assert.False(t, result.Action)
	assert.NotEmpty(t, result.Error)
}

func TestRulesGet_XMLParseError(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte(`<invalid xml><<<`), nil
	})

	result := RulesGet("dummy")

	assert.False(t, result.Action)
	assert.NotEmpty(t, result.Error)
}

func TestRulesGet_RuleWithoutRuleElement(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		xml := `<constraints>
  <rsc_location id="loc-1" rsc="dummy" node="node1" score="INFINITY"/>
</constraints>`
		return []byte(xml), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Empty(t, result.Data)
}

func TestRulesGet_RuleWithNilExpression(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		xml := `<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="rule-no-expr" score="INFINITY"/>
  </rsc_location>
</constraints>`
		return []byte(xml), nil
	})

	result := RulesGet("dummy")

	assert.True(t, result.Action)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "rule-no-expr", result.Data[0].RuleId)
	assert.Equal(t, "INFINITY", result.Data[0].Score)
	assert.Empty(t, result.Data[0].Attribute)
	assert.Empty(t, result.Data[0].Operation)
	assert.Empty(t, result.Data[0].Value)
}

// ==================== RulesDelete ====================

func TestRulesDelete_Success(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte(""), nil
	})

	result := RulesDelete(&validations.DeleteRuleS{RuleIDs: []string{"rule-1", "rule-2"}})

	assert.True(t, result.Action)
	assert.Equal(t, "Delete rule success", result.Info)
	assert.Empty(t, result.Error)
}

func TestRulesDelete_SingleRule(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		assert.Contains(t, cmd, "pcs constraint rule delete")
		return []byte(""), nil
	})

	result := RulesDelete(&validations.DeleteRuleS{RuleIDs: []string{"rule-1"}})

	assert.True(t, result.Action)
}

func TestRulesDelete_PartialFailure(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "rule-bad") {
			return []byte("Error: rule not found"), errors.New("not found")
		}
		return []byte(""), nil
	})

	result := RulesDelete(&validations.DeleteRuleS{RuleIDs: []string{"rule-good", "rule-bad"}})

	assert.False(t, result.Action)
	assert.Len(t, result.Error, 1)
	assert.Equal(t, "rule-bad", result.Error[0]["id"])
}

func TestRulesDelete_AllFail(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte("Error: rule not found"), errors.New("not found")
	})

	result := RulesDelete(&validations.DeleteRuleS{RuleIDs: []string{"rule-1", "rule-2"}})

	assert.False(t, result.Action)
	assert.Len(t, result.Error, 2)
}

func TestRulesDelete_EmptyList(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		t.Fatal("should not call RunCommand for empty list")
		return nil, nil
	})

	result := RulesDelete(&validations.DeleteRuleS{RuleIDs: []string{}})

	assert.True(t, result.Action)
	assert.Equal(t, "Delete rule success", result.Info)
}

// ==================== RuleAdd ====================

func TestRuleAdd_WithRuleID(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		assert.Contains(t, cmd, "pcs constraint location")
		assert.Contains(t, cmd, "score=")
		assert.Contains(t, cmd, "id=")
		return []byte(""), nil
	})

	result := RuleAdd(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.True(t, result.Action)
	assert.Equal(t, "Add rule success", result.Info)
}

func TestRuleAdd_WithoutRuleID(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		assert.Contains(t, cmd, "pcs constraint location")
		assert.Contains(t, cmd, "score=")
		assert.NotContains(t, cmd, "id=")
		return []byte(""), nil
	})

	result := RuleAdd(&validations.RuleS{
		Rsc: "dummy", Score: "100",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.True(t, result.Action)
	assert.Equal(t, "Add rule success", result.Info)
}

func TestRuleAdd_EmptyValue(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		assert.Contains(t, cmd, "'defined' '#uname'")
		return []byte(""), nil
	})

	result := RuleAdd(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY",
		Attribute: "#uname", Operation: "defined", Value: "",
	})

	assert.True(t, result.Action)
	assert.Equal(t, "Add rule success", result.Info)
}

func TestRuleAdd_CommandFailure(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte("Error: duplicate constraint"), errors.New("duplicate")
	})

	result := RuleAdd(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.False(t, result.Action)
	assert.NotEmpty(t, result.Error)
}

// ==================== RuleUpdate ====================

func TestRuleUpdate_Success(t *testing.T) {
	callCount := 0
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		callCount++
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="my-rule" score="100">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
</constraints>`), nil
		}
		return []byte(""), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "200", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node2",
	})

	assert.True(t, result.Action)
	assert.Equal(t, "Update rule success", result.Info)
	assert.Equal(t, 3, callCount) // cibadmin + delete + add
}

func TestRuleUpdate_QueryCommandFails(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte("Error: cibadmin failed"), errors.New("cibadmin error")
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.False(t, result.Action)
	assert.NotEmpty(t, result.Error)
}

func TestRuleUpdate_XMLParseError(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		return []byte(`<invalid xml><<<`), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.False(t, result.Action)
	assert.NotEmpty(t, result.Error)
}

func TestRuleUpdate_DeleteFails(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="my-rule" score="100">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
</constraints>`), nil
		}
		if strings.Contains(cmd, "delete") {
			return []byte("Error: rule not found"), errors.New("not found")
		}
		return []byte(""), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "200", RuleID: "nonexistent-rule",
		Attribute: "#uname", Operation: "eq", Value: "node2",
	})

	assert.False(t, result.Action)
	assert.Contains(t, result.Error, "not found")
}

func TestRuleUpdate_AddNewRuleFails_TriggersRecovery(t *testing.T) {
	var executedCmds []string
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="my-rule" score="100">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
</constraints>`), nil
		}
		executedCmds = append(executedCmds, cmd)
		if strings.Contains(cmd, "delete") {
			return []byte(""), nil
		}
		if strings.Contains(cmd, "node2") {
			return []byte("Error: duplicate"), errors.New("duplicate")
		}
		return []byte(""), nil

	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "200", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node2",
	})

	assert.False(t, result.Action)
	assert.Contains(t, result.Error, "duplicate constraint")
	assert.Len(t, executedCmds, 3) // delete + add(new) + add(recovery)
	assert.Contains(t, executedCmds[0], "rule delete")
	assert.Contains(t, executedCmds[1], "score='200'")
	assert.Contains(t, executedCmds[1], "'node2'")
	assert.Contains(t, executedCmds[2], "'dummy'")
	assert.Contains(t, executedCmds[2], "score='100'")
	assert.Contains(t, executedCmds[2], "'node1'")
}

func TestRuleUpdate_RuleIDNotFoundInXML(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="other-rule" score="100">
      <expression id="e1" attribute="#uname" operation="eq" value="node1"/>
    </rule>
  </rsc_location>
</constraints>`), nil
		}
		return []byte(""), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "200", RuleID: "nonexistent-rule",
		Attribute: "#uname", Operation: "eq", Value: "node2",
	})

	assert.True(t, result.Action)
}

func TestRuleUpdate_RuleWithNilExpression(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints>
  <rsc_location id="loc-1" rsc="dummy">
    <rule id="my-rule" score="100"/>
  </rsc_location>
</constraints>`), nil
		}
		return []byte(""), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "200", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node2",
	})

	assert.True(t, result.Action)
	assert.Equal(t, "Update rule success", result.Info)
}

func TestRuleUpdate_EmptyConstraints(t *testing.T) {
	ruleMockCmd(t, func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			return []byte(`<constraints></constraints>`), nil
		}
		if strings.Contains(cmd, "delete") {
			return []byte("Error: not found"), errors.New("not found")
		}
		return []byte(""), nil
	})

	result := RuleUpdate(&validations.RuleS{
		Rsc: "dummy", Score: "INFINITY", RuleID: "my-rule",
		Attribute: "#uname", Operation: "eq", Value: "node1",
	})

	assert.False(t, result.Action)
}
