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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/openeuler/ha-api/settings"
	"github.com/beego/beego/v2/core/logs"
	"github.com/spf13/viper"
)

func IsInSlice(str string, sli []string) bool {
	//TODO

	for _, item := range sli {
		if item == str {
			return true
		}
	}
	return false
}

// RemoveDupl remove duplicates in string array
func RemoveDupl(strs []string) []string {
	strSet := map[string]bool{}
	for _, v := range strs {
		strSet[v] = true
	}
	strsDupl := []string{}
	for k := range strSet {
		strsDupl = append(strsDupl, k)
	}
	return strsDupl
}

// GetNumAndUnitFromStr gets the first number and the unit after this number
// like "20.5min" ==> ["20.5", "min"]
func GetNumAndUnitFromStr(s string) (string, string) {
	r := regexp.MustCompile("[0-9](.*)[0-9]")
	index := r.FindStringIndex(s)
	if len(index) == 0 {
		return s[:1], s[1:]
	}
	return s[:index[1]], s[index[1]:]
}

// Send http request
func SendRequest(url string, method string, data interface{}) (resp *http.Response, err error) {
	var httpResp *http.Response
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	switch method {
	case "POST":
		httpResp, err = client.Post(url, "application/json", bytes.NewReader(data.([]byte)))
	case "GET":
		httpResp, err = client.Get(url)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", url, bytes.NewReader(data.([]byte)))
		httpResp, err = client.Do(req)
	case "PUT":

		req, _ := http.NewRequest("PUT", url, bytes.NewReader(data.([]byte)))
		httpResp, err = client.Do(req)
	default:
		return nil, errors.New("unsupported method")
	}

	if err != nil {
		return nil, err
	}

	return httpResp, nil
}

// Read Port from config file
func ReadPortFromConfig() (string, error) {
	defaultPort := "8080"
	viper.SetConfigName("port")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logs.Error("Error reading config file, %s, using default port %s", err, defaultPort)
		return defaultPort, fmt.Errorf("Error reading config file, %s, using default port %s", err, defaultPort)
	} else {
		port := viper.GetString("port.haapi_port")
		_, err := strconv.Atoi(port)
		if err != nil {
			logs.Warning("Port in config is not a number, using default port 8080")
			return defaultPort, fmt.Errorf("Port in config is not a number, using default port %s", defaultPort)
		}
		return port, nil
	}
}

func IsLocalCluster(clusterName string) bool {
	localClusterName, err := getLocalClusterName()
	if err != nil {
		return false
	}
	if localClusterName != clusterName {
		return false
	}
	return true
}

func getLocalClusterName() (string, error) {
	filename := settings.CorosyncConfFile
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("open corosync conf failed")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "cluster_name" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("not found cluster name info in corosync conf")
}

func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func RemoveByValue[T comparable](slice []T, value T) []T {
	var newSlice []T
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func Pop[T any](slice []T) (T, []T) {
	if len(slice) == 0 {
		var zero T
		return zero, slice
	}
	last := slice[len(slice)-1]
	return last, slice[:len(slice)-1]
}

func PopFirst[T any](slice []T) (T, []T) {
	if len(slice) == 0 {
		var zero T
		return zero, slice
	}
	first := slice[0]
	return first, slice[1:]
}

// 计算两个切片的差集
func DifferenceSlice[T comparable](mainSlice, subtractSlice []T) []T {
	subtractSet := make(map[T]struct{})
	for _, item := range subtractSlice {
		subtractSet[item] = struct{}{}
	}

	var diffSlice []T
	for _, item := range mainSlice {
		if _, found := subtractSet[item]; !found {
			diffSlice = append(diffSlice, item)
		}
	}

	return diffSlice
}
