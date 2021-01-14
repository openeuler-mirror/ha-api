package utils

import "os"

func GetNodeList() []string {
	// TODO:
	nodeList := []string{}

	return nodeList
}

func IsClusterExist() bool {
	_, err := os.Lstat("/etc/corosync/corosync.conf")
	return !os.IsNotExist(err)
}
