/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Jason011125 <zic022@ucsd.edu>
 * Date: Mon Aug 14 15:53:52 2023 +0800
 */
package models

import (
	"bufio"
	"os"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
)

func IsClusterExist() bool {
	_, err := os.Stat(settings.CorosyncConfFile)
	return err == nil
}

func getNodeList() []string {

	filename := settings.CorosyncConfFile
	nodeList := []string{} // 存储文件内数据

	file, err := os.Open(filename)
	if err != nil {
		return []string{"error " + "File " + filename + " doesn't exist!"}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundNodeList := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "nodelist {") {
			foundNodeList = true
			break
		}
	}

	if !foundNodeList {
		return []string{"No nodelist found in the file."}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "quorum {") {
			break
		}
		if line != "nodelist {" && line != "}" && line != "" {
			nodeList = append(nodeList, line)
		}
	}

	return nodeList
}
