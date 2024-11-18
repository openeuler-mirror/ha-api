/*
 * Copyright (c) KylinSoft Co., Ltd.2024. All Rights Reserved.
 *  ha-api is licensed under the Mulan PSL v2.
 *  You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 * 	http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: bizhiyuan
 * @Date: 2024-03-19 11:19:33
 * @LastEditTime: 2024-03-27 09:27:48
 * @Description:
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
