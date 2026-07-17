/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Mon Sep 2 15:58:49 2024 +0800
 */

package models

import (
	"errors"
	"strings"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

var savedRunCommand = utils.RunCommand

func mockRunCmd(mock func(string) ([]byte, error)) {
	utils.RunCommand = mock
}

func restoreRunCmd() {
	utils.RunCommand = savedRunCommand
}

// ==================== AlarmsGet ====================

func TestAlarmsGet_Success(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="test@example.com"/>
        <nvpair name="email_server" value="smtp.example.com"/>
        <nvpair name="password" value="enc_pwd"/>
        <nvpair name="port" value="587"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
      <recipient value="user1@example.com"/>
      <recipient value="user2@example.com"/>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		// pwd_decode
		return []byte("decoded_password"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.True(t, result.Data.Flag)
	assert.Equal(t, "test@example.com", result.Data.Sender)
	assert.Equal(t, "smtp.example.com", result.Data.Smtp)
	assert.Equal(t, 587.0, result.Data.Port)
	assert.Equal(t, "decoded_password", result.Data.Password)
	assert.Equal(t, []string{"user1@example.com", "user2@example.com"}, result.Data.Receiver)
}

func TestAlarmsGet_CommandError(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		return nil, errors.New("command failed")
	})

	result := AlarmsGet()

	// 即使命令失败，Action 仍被设为 true（goto ret 后赋值）
	assert.True(t, result.Action)
	assert.Empty(t, result.Data.Sender)
	assert.Empty(t, result.Data.Smtp)
}

func TestAlarmsGet_XMLParseError(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		return []byte(`<invalid xml><<<`), nil
	})

	result := AlarmsGet()

	// XML 解析失败走 goto ret，Action 仍为 true
	assert.True(t, result.Action)
	assert.Empty(t, result.Data.Sender)
}

func TestAlarmsGet_SwitchOff(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="a@b.com"/>
        <nvpair name="email_server" value="smtp.b.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="off"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		return []byte("plain_pwd"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.False(t, result.Data.Flag)
}

func TestAlarmsGet_NoRecipients(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="a@b.com"/>
        <nvpair name="email_server" value="smtp.b.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		return []byte("decoded"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.True(t, result.Data.Flag)
	assert.Empty(t, result.Data.Receiver)
}

func TestAlarmsGet_PasswordDecryptReturnsLess(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="password" value="enc_pwd"/>
        <nvpair name="switCh" value="off"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		// pwd_decode 返回 "the parameter is less" 表示解密失败
		return []byte("the parameter is less\n"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.Equal(t, "", result.Data.Password)
}

func TestAlarmsGet_PasswordDecryptFails(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="password" value="enc_pwd"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		// pwd_decode 命令执行失败，返回错误文本和 error
		return []byte("error: decode failed"), errors.New("decode failed")
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	// 失败时密码应为空，而非错误文本
	assert.Equal(t, "", result.Data.Password)
}

func TestAlarmsGet_NoAlertsSection(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration></configuration>`
			return []byte(xml), nil
		}
		return []byte("decoded"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.Empty(t, result.Data.Sender)
	assert.Empty(t, result.Data.Smtp)
	assert.Equal(t, 0.0, result.Data.Port)
	assert.False(t, result.Data.Flag)
}

func TestAlarmsGet_PartialFields(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="only@sender.com"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		return []byte("decoded"), nil
	})

	result := AlarmsGet()

	assert.True(t, result.Action)
	assert.Equal(t, "only@sender.com", result.Data.Sender)
	assert.Empty(t, result.Data.Smtp)
	assert.Equal(t, 0.0, result.Data.Port)
	assert.False(t, result.Data.Flag)
}

// ==================== isDataEmpty ====================

func TestIsDataEmpty_ZeroValue(t *testing.T) {
	assert.True(t, isDataEmpty(AlarmData{}))
}

func TestIsDataEmpty_WithSender(t *testing.T) {
	assert.False(t, isDataEmpty(AlarmData{Sender: "test@example.com"}))
}

func TestIsDataEmpty_WithPortOnly(t *testing.T) {
	assert.False(t, isDataEmpty(AlarmData{Port: 25}))
}

func TestIsDataEmpty_WithReceiverOnly(t *testing.T) {
	assert.False(t, isDataEmpty(AlarmData{Receiver: []string{"a@b.com"}}))
}

func TestIsDataEmpty_WithFlagTrue(t *testing.T) {
	assert.False(t, isDataEmpty(AlarmData{Flag: true}))
}

// ==================== AlarmsSet ====================

func TestAlarmsSet_Success(t *testing.T) {
	defer restoreRunCmd()
	var executedCmds []string
	mockRunCmd(func(cmd string) ([]byte, error) {
		executedCmds = append(executedCmds, cmd)
		if strings.Contains(cmd, "pwd_encode") {
			return []byte("encoded_secret\n"), nil
		}
		return []byte(""), nil
	})

	data := AlarmData{
		Flag:     false,
		Sender:   "admin@example.com",
		Smtp:     "smtp.example.com",
		Password: "secret",
		Port:     465,
		Receiver: []string{"user1@example.com", "user2@example.com"},
	}

	result := AlarmsSet(data)

	assert.True(t, result["action"].(bool))
	assert.Equal(t, "Set alarm success", result["info"])
	// delete + pwd_encode + create + 2 recipients = 5 commands
	assert.Len(t, executedCmds, 5)
	assert.Contains(t, executedCmds[0], "alert delete")
	assert.Contains(t, executedCmds[1], "pwd_encode")
	assert.Contains(t, executedCmds[2], "alert create")
	assert.Contains(t, executedCmds[3], "alert recipient add")
	assert.Contains(t, executedCmds[4], "alert recipient add")
	// verify the encoded password is used in the create command, not the raw password
	assert.Contains(t, executedCmds[2], "password='encoded_secret'")
	assert.NotContains(t, executedCmds[2], "password='secret'")
}

func TestAlarmsSet_CreateAlertFails(t *testing.T) {
	defer restoreRunCmd()
	callCount := 0
	mockRunCmd(func(cmd string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			// CmdDeleteAlert — error is ignored in source
			return nil, errors.New("delete failed")
		}
		// CmdCreateAlert fails
		return nil, errors.New("permission denied")
	})

	data := AlarmData{Sender: "test@example.com"}
	result := AlarmsSet(data)

	assert.False(t, result["action"].(bool))
	assert.Equal(t, "Set alarm failed", result["error"])
}

func TestAlarmsSet_RecipientAddFails(t *testing.T) {
	defer restoreRunCmd()
	callCount := 0
	mockRunCmd(func(cmd string) ([]byte, error) {
		callCount++
		if callCount <= 2 {
			// delete + create succeed
			return []byte(""), nil
		}
		// recipient add fails
		return nil, errors.New("recipient error")
	})

	data := AlarmData{
		Sender:   "a@b.com",
		Receiver: []string{"bad@example.com"},
	}
	result := AlarmsSet(data)

	assert.False(t, result["action"].(bool))
	assert.Equal(t, "Set alarm failed", result["error"])
}

func TestAlarmsSet_EmptyReceiver(t *testing.T) {
	defer restoreRunCmd()
	var executedCmds []string
	mockRunCmd(func(cmd string) ([]byte, error) {
		executedCmds = append(executedCmds, cmd)
		return []byte(""), nil
	})

	data := AlarmData{
		Sender: "a@b.com",
		Smtp:   "smtp.b.com",
		Port:   25,
	}
	result := AlarmsSet(data)

	assert.True(t, result["action"].(bool))
	// only delete + create, no recipient commands
	assert.Len(t, executedCmds, 2)
}

func TestAlarmsSet_FlagTrue_SwitChOn(t *testing.T) {
	defer restoreRunCmd()
	var createCmd string
	callCount := 0
	mockRunCmd(func(cmd string) ([]byte, error) {
		callCount++
		if callCount == 2 {
			createCmd = cmd
		}
		return []byte(""), nil
	})

	data := AlarmData{
		Flag:   true,
		Sender: "a@b.com",
	}
	AlarmsSet(data)

	assert.Contains(t, createCmd, "switCh='on'")
}

func TestAlarmsSet_FlagFalse_SwitChOff(t *testing.T) {
	defer restoreRunCmd()
	var createCmd string
	callCount := 0
	mockRunCmd(func(cmd string) ([]byte, error) {
		callCount++
		if callCount == 2 {
			createCmd = cmd
		}
		return []byte(""), nil
	})

	data := AlarmData{
		Flag:   false,
		Sender: "a@b.com",
	}
	AlarmsSet(data)

	assert.Contains(t, createCmd, "switCh='off'")
}

func TestAlarmsSet_PasswordEncodeFails(t *testing.T) {
	defer restoreRunCmd()
	var executedCmds []string
	mockRunCmd(func(cmd string) ([]byte, error) {
		executedCmds = append(executedCmds, cmd)
		if strings.Contains(cmd, "pwd_encode") {
			return []byte("error: encode failed"), errors.New("encode failed")
		}
		return []byte(""), nil
	})

	data := AlarmData{
		Sender:   "a@b.com",
		Password: "secret",
		Port:     25,
	}
	result := AlarmsSet(data)

	assert.False(t, result["action"].(bool))
	assert.Equal(t, "Set alarm failed", result["error"])
	// only delete + pwd_encode attempted; create alert must NOT be reached
	assert.Len(t, executedCmds, 2)
	assert.Contains(t, executedCmds[0], "alert delete")
	assert.Contains(t, executedCmds[1], "pwd_encode")
}

// ==================== AlarmsTest ====================

func TestAlarmsTest_Success(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="sender@test.com"/>
        <nvpair name="email_server" value="smtp.test.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
      <recipient value="recv@test.com"/>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		return []byte(""), nil
	})

	result := AlarmsTest()

	assert.True(t, result.Action)
	assert.Equal(t, "Send alarm test success", result.Info)
}

func TestAlarmsTest_SendFails(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="sender@test.com"/>
        <nvpair name="email_server" value="smtp.test.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
      <recipient value="recv@test.com"/>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		if strings.Contains(cmd, "pwd_decode") {
			return []byte("decoded_pwd"), nil
		}
		if strings.Contains(cmd, "python_email") {
			return []byte("SMTP connection refused"), errors.New("send failed")
		}
		// echo log command
		return []byte(""), nil
	})

	result := AlarmsTest()

	assert.False(t, result.Action)
	assert.Equal(t, "Send alarm test failed", result.Error)
}

func TestAlarmsTest_NoRecipients(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		return []byte(""), nil
	})

	result := AlarmsTest()

	// 没有收件人时循环不执行，直接返回成功
	assert.True(t, result.Action)
	assert.Equal(t, "Send alarm test success", result.Info)
}

func TestAlarmsTest_MultipleRecipients_SecondFails(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="sender@test.com"/>
        <nvpair name="email_server" value="smtp.test.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
      <recipient value="ok@test.com"/>
      <recipient value="fail@test.com"/>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		if strings.Contains(cmd, "pwd_decode") {
			return []byte("decoded_pwd"), nil
		}
		if strings.Contains(cmd, "python_email") {
			if strings.Contains(cmd, "fail@test.com") {
				return []byte("timeout"), errors.New("send failed")
			}
			return []byte(""), nil
		}
		return []byte(""), nil
	})

	result := AlarmsTest()

	assert.False(t, result.Action)
	assert.Equal(t, "Send alarm test failed", result.Error)
}

func TestAlarmsTest_LogWriteFails(t *testing.T) {
	defer restoreRunCmd()
	mockRunCmd(func(cmd string) ([]byte, error) {
		if strings.Contains(cmd, "cibadmin") {
			xml := `<configuration>
  <alerts>
    <alert>
      <instance_attributes>
        <nvpair name="email_sender" value="sender@test.com"/>
        <nvpair name="email_server" value="smtp.test.com"/>
        <nvpair name="password" value="pwd"/>
        <nvpair name="port" value="25"/>
        <nvpair name="switCh" value="on"/>
      </instance_attributes>
      <recipient value="recv@test.com"/>
    </alert>
  </alerts>
</configuration>`
			return []byte(xml), nil
		}
		if strings.Contains(cmd, "pwd_decode") {
			return []byte("decoded_pwd"), nil
		}
		if strings.Contains(cmd, "python_email") {
			return []byte("SMTP connection refused"), errors.New("send failed")
		}
		// echo 日志写入也失败
		return nil, errors.New("log write failed")
	})

	result := AlarmsTest()

	// 即使日志写入失败，主流程返回值不应受影响
	assert.False(t, result.Action)
	assert.Equal(t, "Send alarm test failed", result.Error)
}

