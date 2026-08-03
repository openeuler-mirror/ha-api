/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package models

import (
	"errors"
	"fmt"
	"testing"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllResourceMetas_Normal(t *testing.T) {
	utils.RunCommand = func(cmd string) ([]byte, error) {
		switch cmd {
		case utils.CmdListResourceStandards:
			return []byte("ocf\n\nlsb\nsystemd\n"), nil
		case utils.CmdListOcfProviders:
			return []byte("\npacemaker\nopenstack\n"), nil // 头部空行
		case fmt.Sprintf(utils.CmdListOcfResourceAgent, "pacemaker"):
			return []byte("\nDummy\n\nVirtualIP\n"), nil // 多重空行
		case fmt.Sprintf(utils.CmdListOcfResourceAgent, "openstack"):
			return []byte("Nova\n\nGlance\n"), nil
		case fmt.Sprintf(utils.CmdListResourceAgent, "systemd"):
			return []byte("nginx@\n\nredis\n"), nil
		default:
			return nil, fmt.Errorf("unexpected command: %s", cmd)
		}
	}

	result := GetAllResourceMetas()
	assert.True(t, result["action"].(bool))
}

func TestGetAllResourceMetas_CommandError(t *testing.T) {
	utils.RunCommand = func(cmd string) ([]byte, error) {
		return nil, errors.New("permission denied")
	}

	result := GetAllResourceMetas()

	assert.False(t, result["action"].(bool))
	assert.Contains(t, result["error"].(string), "permission denied")
}

func TestGetResourceMetas_Normal(t *testing.T) {
	// Mock 正常XML响应
	mockXML := `
    <resource name="DummyResource" >
        <parameters>
            <parameter name="ip" required="1">
                <content type="string" default="127.0.0.1"/>
                <shortdesc>IP address</shortdesc>
                <longdesc>Full description\nwith newline</longdesc>
            </parameter>
        </parameters>
        <actions>
            <action name="start" timeout="20s"/>
        </actions>
        <version>1.0.0</version>
        <shortdesc>Short description</shortdesc>
        <longdesc>Long description\nwith newline</longdesc>
    </resource>`

	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte(mockXML), nil
	}

	// 执行测试
	result := GetResourceMetas("ocf", "Dummy", "pacemaker")

	// 验证结果结构
	assert.True(t, result["action"].(bool))

}

func TestGetResourceMetas_EdgeCases(t *testing.T) {
	t.Run("EmptyMetadata", func(t *testing.T) {
		utils.RunCommand = func(cmd string) ([]byte, error) {
			return []byte("<resource/>"), nil
		}

		result := GetResourceMetas("ocf", "Empty", "")
		data := result["data"].(map[string]interface{})

		// 验证空值处理
		assert.Equal(t, "", data["name"])
		assert.NotEmpty(t, data["parameters"])
		// assert.Empty(t, data["actions"])
	})

	t.Run("MissingElements", func(t *testing.T) {
		mockXML := `<resource>
            <parameters>
                <parameter name="port">
                    <content type="string"/>
                </parameter>
            </parameters>
        </resource>`

		utils.RunCommand = func(cmd string) ([]byte, error) {
			return []byte(mockXML), nil
		}

		result := GetResourceMetas("ocf", "MissingElements", "")
		data := result["data"].(map[string]interface{})

		// 验证缺失字段的默认值
		assert.Equal(t, "", data["version"])
		assert.Equal(t, "", data["shortdesc"])
		assert.Equal(t, "", data["longdesc"])
	})
}

func TestGetResourceMetasStonith(t *testing.T) {
	// Mock 正常XML响应
	mockXML := `
    <resource-agent name="fence_sbd" shortdesc="Fence agent for sbd">
  <longdesc>
    fence_sbd is an I/O Fencing agent which can be used in environments where sbd can be used (shared storage).
  </longdesc>
  <vendor-url/>
  <parameters>
    <parameter name="action" unique="0" required="0">
      <getopt mixed="-o, --action=[action]"/>
      <content type="string" default="reboot"/>
      <shortdesc lang="en">
        Fencing action
      </shortdesc>
    </parameter>
    <parameter name="devices" unique="0" required="1">
      <getopt mixed="--devices=[device_a,device_b]"/>
      <content type="string"/>
      <shortdesc lang="en">
        SBD Device
      </shortdesc>
    </parameter>
    <parameter name="method" unique="0" required="0">
      <getopt mixed="-m, --method=[method]"/>
      <content type="select" default="cycle">
        <option value="onoff"/>
        <option value="cycle"/>
      </content>
      <shortdesc lang="en">
        Method to fence
      </shortdesc>
    </parameter>
    <parameter name="plug" unique="0" required="0" obsoletes="port">
      <getopt mixed="-n, --plug=[id]"/>
      <content type="string"/>
      <shortdesc lang="en">
        Physical plug number on device, UUID or identification of machine
      </shortdesc>
    </parameter>
    <parameter name="port" unique="0" required="0" deprecated="1">
      <getopt mixed="-n, --plug=[id]"/>
      <content type="string"/>
      <shortdesc lang="en">
        Physical plug number on device, UUID or identification of machine
      </shortdesc>
    </parameter>
    <parameter name="quiet" unique="0" required="0">
      <getopt mixed="-q, --quiet"/>
      <content type="boolean"/>
      <shortdesc lang="en">
        Disable logging to stderr. Does not affect --verbose or --debug-file or logging to syslog.
      </shortdesc>
    </parameter>
    <parameter name="verbose" unique="0" required="0">
      <getopt mixed="-v, --verbose"/>
      <content type="boolean"/>
      <shortdesc lang="en">
        Verbose mode. Multiple -v flags can be stacked on the command line (e.g., -vvv) to increase verbosity.
      </shortdesc>
    </parameter>
    <parameter name="verbose_level" unique="0" required="0">
      <getopt mixed="--verbose-level"/>
      <content type="integer"/>
      <shortdesc lang="en">
        Level of debugging detail in output. Defaults to the number of --verbose flags specified on the command line, or to 1 if verbose=1 in a stonith device configuration (i.e., on stdin).
      </shortdesc>
    </parameter>
    <parameter name="debug" unique="0" required="0" deprecated="1">
      <getopt mixed="-D, --debug-file=[debugfile]"/>
      <content type="string"/>
      <shortdesc lang="en">
        Write debug information to given file
      </shortdesc>
    </parameter>
    <parameter name="debug_file" unique="0" required="0" obsoletes="debug">
      <getopt mixed="-D, --debug-file=[debugfile]"/>
      <content type="string"/>
      <shortdesc lang="en">
        Write debug information to given file
      </shortdesc>
    </parameter>
    <parameter name="version" unique="0" required="0">
      <getopt mixed="-V, --version"/>
      <content type="boolean"/>
      <shortdesc lang="en">
        Display version information and exit
      </shortdesc>
    </parameter>
    <parameter name="help" unique="0" required="0">
      <getopt mixed="-h, --help"/>
      <content type="boolean"/>
      <shortdesc lang="en">
        Display help and exit
      </shortdesc>
    </parameter>
    <parameter name="plug_separator" unique="0" required="0">
      <getopt mixed="--plug-separator=[char]"/>
      <content type="string" default=","/>
      <shortdesc lang="en">
        Separator for plug parameter when specifying more than 1 plug
      </shortdesc>
    </parameter>
    <parameter name="separator" unique="0" required="0">
      <getopt mixed="-C, --separator=[char]"/>
      <content type="string" default=","/>
      <shortdesc lang="en">
        Separator for CSV created by 'list' operation
      </shortdesc>
    </parameter>
    <parameter name="delay" unique="0" required="0">
      <getopt mixed="--delay=[seconds]"/>
      <content type="second" default="0"/>
      <shortdesc lang="en">
        Wait X seconds before fencing is started
      </shortdesc>
    </parameter>
    <parameter name="disable_timeout" unique="0" required="0">
      <getopt mixed="--disable-timeout=[true/false]"/>
      <content type="string"/>
      <shortdesc lang="en">
        Disable timeout (true/false) (default: true when run from Pacemaker 2.0+)
      </shortdesc>
    </parameter>
    <parameter name="login_timeout" unique="0" required="0">
      <getopt mixed="--login-timeout=[seconds]"/>
      <content type="second" default="5"/>
      <shortdesc lang="en">
        Wait X seconds for cmd prompt after login
      </shortdesc>
    </parameter>
    <parameter name="power_timeout" unique="0" required="0">
      <getopt mixed="--power-timeout=[seconds]"/>
      <content type="second" default="30"/>
      <shortdesc lang="en">
        Test X seconds for status change after ON/OFF
      </shortdesc>
    </parameter>
    <parameter name="power_wait" unique="0" required="0">
      <getopt mixed="--power-wait=[seconds]"/>
      <content type="second" default="0"/>
      <shortdesc lang="en">
        Wait X seconds after issuing ON/OFF
      </shortdesc>
    </parameter>
    <parameter name="sbd_path" unique="0" required="0">
      <getopt mixed="--sbd-path=[path]"/>
      <content type="string" default="/usr/sbin/sbd"/>
      <shortdesc lang="en">
        Path to SBD binary
      </shortdesc>
    </parameter>
    <parameter name="shell_timeout" unique="0" required="0">
      <getopt mixed="--shell-timeout=[seconds]"/>
      <content type="second" default="3"/>
      <shortdesc lang="en">
        Wait X seconds for cmd prompt after issuing command
      </shortdesc>
    </parameter>
    <parameter name="stonith_status_sleep" unique="0" required="0">
      <getopt mixed="--stonith-status-sleep=[seconds]"/>
      <content type="second" default="1"/>
      <shortdesc lang="en">
        Sleep X seconds between status calls during a STONITH action
      </shortdesc>
    </parameter>
    <parameter name="retry_on" unique="0" required="0">
      <getopt mixed="--retry-on=[attempts]"/>
      <content type="integer" default="1"/>
      <shortdesc lang="en">
        Count of attempts to retry power on
      </shortdesc>
    </parameter>
  </parameters>
  <actions>
    <action name="on" automatic="0"/>
    <action name="off"/>
    <action name="reboot"/>
    <action name="status"/>
    <action name="list"/>
    <action name="list-status"/>
    <action name="monitor"/>
    <action name="metadata"/>
    <action name="manpage"/>
    <action name="validate-all"/>
    <action name="stop" timeout="20s"/>
    <action name="start" timeout="20s"/>
  </actions>
</resource-agent>`

	utils.RunCommand = func(cmd string) ([]byte, error) {
		return []byte(mockXML), nil
	}

	// 执行测试
	result := GetResourceMetas("stonith", "fence_sbd", "")

	// 验证结果结构
	assert.True(t, result["action"].(bool))
}
