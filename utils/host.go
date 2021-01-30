package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetNodeList() ([]map[string]string, error) {
	config, err := getCorosyncConfig()
	if err != nil {
		return nil, errors.New("read config from /etc/corosync/corosync.conf failed")
	}
	return config.NodeList, nil
}

func IsClusterExist() bool {
	_, err := os.Lstat("/etc/corosync/corosync.conf")
	return !os.IsNotExist(err)
}

type CorosyncConfig struct {
	Totem    map[string]string
	NodeList []map[string]string
	Quorum   map[string]string
	Logging  map[string]string
}

func getCorosyncConfig() (CorosyncConfig, error) {
	var result CorosyncConfig
	f, err := os.Open("/etc/corosync/corosync.conf")
	if err != nil {
		return result, err
	}
	defer f.Close()

	const (
		StateRoot       = 0
		StateInTotem    = 1
		StateInNodeList = 2
		StateInNode     = 3
		StateInQuorum   = 4
		StateInLogging  = 5
	)
	var state = StateRoot
	bf := bufio.NewReader(f)
	currentNode := map[string]string{}
	for {
		l, _, err := bf.ReadLine()
		line := strings.Trim(string(l), " ")
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return result, err
			}
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
				fmt.Println("found node index")
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
			if line != "}" {
				words := strings.Split(line, ":")
				key := strings.Trim(words[0], " ")
				value := strings.Trim(words[1], " ")
				if result.Logging == nil {
					result.Logging = make(map[string]string)
				}
				result.Logging[key] = value
			} else {
				state = StateRoot
			}
		default:
			return result, errors.New("parse corosync.conf failed, invalid state")
		}
	}

	return result, nil
}
