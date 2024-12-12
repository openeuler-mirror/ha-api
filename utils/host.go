/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"github.com/beego/beego/v2/core/logs"
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

func setCorosyncConfig(conf CorosyncConfig, confFile string) error {
	var sb strings.Builder
	// Write totem section
	sb.WriteString("totem {\n")
	for key, value := range conf.Totem {
		sb.WriteString(fmt.Sprintf("    %s: %s\n", key, value))
	}
	sb.WriteString("}\n\n")

	// Write nodelist section
	sb.WriteString("nodelist {\n")
	for _, node := range conf.NodeList {
		sb.WriteString("    node {\n")
		for key, value := range node {
			sb.WriteString(fmt.Sprintf("        %s: %s\n", key, value))
		}
		sb.WriteString("    }\n")
	}
	sb.WriteString("}\n\n")

	// Write quorum section
	sb.WriteString("quorum {\n")
	for key, value := range conf.Quorum {
		sb.WriteString(fmt.Sprintf("    %s: %s\n", key, value))
	}
	sb.WriteString("}\n\n")

	// Write logging section
	sb.WriteString("logging {\n")
	for key, value := range conf.Logging {
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case bool:
			valueStr = fmt.Sprintf("%t", v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}
		sb.WriteString(fmt.Sprintf("    %s: %s\n", key, valueStr))
	}
	sb.WriteString("}\n")

	file, err := os.OpenFile(confFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logs.Error("Failed to open file: %s", err)
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString(sb.String())
	if err != nil {
		logs.Error("Failed to write to file: %s", err)
	}
	return nil
}
