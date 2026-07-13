/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
)

func GetNodeList() ([]map[string]string, error) {
	config, err := GetCorosyncConfig()
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
	Quorum   map[string]interface{}
	Logging  map[string]interface{}
}

func GetCorosyncConfig() (CorosyncConfig, error) {
	var result CorosyncConfig
	f, err := os.Open(settings.CorosyncConfFile)
	if err != nil {
		return result, err
	}
	defer f.Close()

	const (
		StateRoot              = 0
		StateInTotem           = 1
		StateInNodeList        = 2
		StateInNode            = 3
		StateInQuorum          = 4
		StateInLogging         = 5
		StateInLoggerSubsys    = 6
		StateInQuorumDevice    = 7 // 新增：处理 device 块
		StateInQuorumDeviceNet = 8 // 新增：处理 device 内的 net 块

	)
	var state = StateRoot
	bf := bufio.NewReader(f)
	currentNode := map[string]string{}
	currentLoggerSubsys := map[string]string{}
	currentDevice := map[string]interface{}{}
	currentNet := map[string]string{}
	for {
		l, _, err := bf.ReadLine()
		line := strings.TrimSpace(string(l))
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
				if result.Quorum == nil {
					result.Quorum = make(map[string]interface{})
				}
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
			if strings.HasPrefix(line, "device {") {
				currentDevice = make(map[string]interface{})
				state = StateInQuorumDevice
			} else if line == "}" {
				state = StateRoot
			} else {
				// 处理普通的键值对配置
				words := strings.SplitN(line, ":", 2) // 使用 SplitN 防止值中包含冒号
				if len(words) == 2 {
					key := strings.TrimSpace(words[0])
					value := strings.TrimSpace(words[1])
					result.Quorum[key] = value
				}
			}
		case StateInQuorumDevice:
			if strings.HasPrefix(line, "net {") {
				currentNet = make(map[string]string)
				state = StateInQuorumDeviceNet
			} else if line == "}" {
				// device 块结束
				result.Quorum["device"] = currentDevice
				state = StateInQuorum
			} else {
				// 处理 device 内的键值对
				words := strings.SplitN(line, ":", 2)
				if len(words) == 2 {
					key := strings.TrimSpace(words[0])
					value := strings.TrimSpace(words[1])
					currentDevice[key] = value
				}
			}
		case StateInQuorumDeviceNet:
			if line == "}" {
				// net 块结束
				currentDevice["net"] = currentNet
				state = StateInQuorumDevice
			} else {
				// 处理 net 内的键值对
				words := strings.SplitN(line, ":", 2)
				if len(words) == 2 {
					key := strings.TrimSpace(words[0])
					value := strings.TrimSpace(words[1])
					currentNet[key] = value
				}
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

func SetCorosyncConfig(conf CorosyncConfig, confFile string) error {
	var sb strings.Builder
	// Write totem section
	sb.WriteString("totem {\n")
	for key, value := range conf.Totem {
		sb.WriteString(fmt.Sprintf("    %s: %s\n", key, value))
	}
	sb.WriteString("}\n\n")

	// Write nodelist section
	if len(conf.NodeList) > 0 {
		sb.WriteString("nodelist {\n")
		for _, node := range conf.NodeList {
			sb.WriteString("    node {\n")
			for key, value := range node {
				sb.WriteString(fmt.Sprintf("        %s: %s\n", key, value))
			}
			sb.WriteString("    }\n")
		}
		sb.WriteString("}\n\n")
	}

	// Write quorum section
	sb.WriteString("quorum {\n")
	for key, value := range conf.Quorum {
		// 处理 device 嵌套配置
		if key == "device" {
			if device, ok := value.(map[string]interface{}); ok {
				sb.WriteString("    device {\n")
				for dkey, dvalue := range device {
					// 处理 net 子配置
					if dkey == "net" {
						if net, ok := dvalue.(map[string]string); ok {
							sb.WriteString("        net {\n")
							for nkey, nvalue := range net {
								sb.WriteString(fmt.Sprintf("            %s: %s\n", nkey, nvalue))
							}
							sb.WriteString("        }\n")
						}
					} else {
						// 处理 device 的其他键值对
						var valueStr string
						switch v := dvalue.(type) {
						case string:
							valueStr = v
						case bool:
							valueStr = fmt.Sprintf("%t", v)
						case int:
							valueStr = fmt.Sprintf("%d", v)
						default:
							valueStr = fmt.Sprintf("%v", v)
						}
						sb.WriteString(fmt.Sprintf("        %s: %s\n", dkey, valueStr))
					}
				}
				sb.WriteString("    }\n")
			}
		} else {
			// 处理 quorum 的其他键值对
			var valueStr string
			switch v := value.(type) {
			case string:
				valueStr = v
			case bool:
				valueStr = fmt.Sprintf("%t", v)
			case int:
				valueStr = fmt.Sprintf("%d", v)
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			sb.WriteString(fmt.Sprintf("    %s: %s\n", key, valueStr))
		}
	}
	sb.WriteString("}\n\n")

	// Write logging section
	if len(conf.Logging) > 0 {
		sb.WriteString("logging {\n")
		for key, value := range conf.Logging {
			// 处理 logger 子系统配置
			if subsys, ok := value.(map[string]string); ok {
				sb.WriteString(fmt.Sprintf("    %s {\n", key))
				for skey, svalue := range subsys {
					sb.WriteString(fmt.Sprintf("        %s: %s\n", skey, svalue))
				}
				sb.WriteString("    }\n")
			} else {
				// 处理普通键值对
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
		}
		sb.WriteString("}\n\n")
	}

	file, err := os.OpenFile(confFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to open corosync file: %s", err))
		return err
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString(sb.String())
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to write to corosync file: %s", err))
		return err
	}

	return nil
}
