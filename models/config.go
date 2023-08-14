package models

import (
	"bufio"
	"os"
	"strings"
)

func isClusterExist() bool {
	_, err := os.Stat("/etc/corosync/corosync.conf")
	return err == nil
}

func getNodeList() []string {

	filename := "/etc/corosync/corosync.conf"
	node_list := []string{} // 存储文件内数据

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
			node_list = append(node_list, line)
		}
	}

	return node_list
}
