/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 11:22:03
 * Description: 主机corosync配置
 ******************************************************************************/
package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
)

func GetNodeList() ([]map[string]string, error) {
	config, err := getCorosyncConfig()
	if err != nil {
		return nil, errors.New("read config from /etc/corosync/corosync.conf failed")
	}
	return config.NodeList, nil
}

func IsClusterExist() bool {
	_, err := os.Lstat(settings.CorosyncConfFile)
	return !os.IsNotExist(err)
}

type CorosyncConfig struct {
	Totem    map[string]string
	NodeList []map[string]string
	Quorum   map[string]string
	Logging  map[string]interface{}
}

func getCorosyncConfig() (CorosyncConfig, error) {
	var result CorosyncConfig
	f, err := os.Open(settings.CorosyncConfFile)
	if err != nil {
		return result, err
	}
	defer f.Close()

	const (
		StateRoot           = 0
		StateInTotem        = 1
		StateInNodeList     = 2
		StateInNode         = 3
		StateInQuorum       = 4
		StateInLogging      = 5
		StateInLoggerSubsys = 6
	)
	var state = StateRoot
	bf := bufio.NewReader(f)
	currentNode := map[string]string{}
	currentLoggerSubsys := map[string]string{}
	for {
		l, _, err := bf.ReadLine()
		line := strings.Trim(string(l), " ")
		line = strings.Trim(line, "\t")
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return result, err
			}
		}
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// parse line here
		switch state {
		case StateRoot:
			if strings.HasPrefix(line, "totem {") {
				state = StateInTotem
			} else if strings.HasPrefix(line, "nodelist {") {
				state = StateInNodeList
			} else if strings.HasPrefix(line, "quorum {") {
				state = StateInQuorum
			} else if strings.HasPrefix(line, "logging {") {
				state = StateInLogging
			}
		case StateInTotem:
			if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				if result.Totem == nil {
					result.Totem = make(map[string]string)
				}
				result.Totem[key] = value
			} else {
				state = StateRoot
			}
		case StateInNodeList:
			if strings.HasPrefix(line, "node {") {
				currentNode = make(map[string]string)
				state = StateInNode
			} else if line == "}" {
				state = StateRoot
			}
		case StateInNode:
			if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				currentNode[key] = value
			} else {
				if result.NodeList == nil {
					result.NodeList = []map[string]string{}
				}
				result.NodeList = append(result.NodeList, currentNode)
				state = StateInNodeList
			}
		case StateInQuorum:
			if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				if result.Quorum == nil {
					result.Quorum = make(map[string]string)
				}
				result.Quorum[key] = value
			} else {
				state = StateRoot
			}
		case StateInLogging:
			if strings.HasSuffix(line, "{") {
				currentLoggerSubsys = make(map[string]string)
				words := strings.Split(line, " ")
				key := strings.Trim(words[0], " ")
				if result.Logging == nil {
					result.Logging = make(map[string]interface{})
				}
				result.Logging[key] = currentLoggerSubsys
				state = StateInLoggerSubsys
			} else if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				if result.Logging == nil {
					result.Logging = make(map[string]interface{})
				}
				result.Logging[key] = value
			} else {
				state = StateRoot
			}
		case StateInLoggerSubsys:
			if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				currentLoggerSubsys[key] = value
			} else {
				state = StateInLogging
			}
		default:
			return result, errors.New("parse corosync.conf failed, invalid state")
		}
	}

	return result, nil
}
