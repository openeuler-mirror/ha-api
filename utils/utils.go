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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/ha-api/settings"
	"github.com/spf13/viper"
)

var port, _ = ReadPortFromConfig()

var IsInSlice = func(str string, sli []string) bool {
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
func isDevEnvironment() bool {
	env := os.Getenv("HA_API")
	return env != "production"
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

var globalHTTPClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: isDevEnvironment(),
		},
		// 连接池优化配置
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 100,              // 对每个目标host的最大空闲连接数
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时关闭时间
	},
	// 可以设置全局请求超时
	Timeout: 90 * time.Second,
}

var SendRequest = func(url string, method string, data interface{}) (resp *http.Response, err error) {
	var httpResp *http.Response
	var req *http.Request

	// 复用全局的 client
	client := globalHTTPClient

	switch method {
	case "POST":
		// 使用 http.NewRequest 更灵活，可以统一处理Header等
		req, err = http.NewRequest("POST", url, bytes.NewReader(data.([]byte)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	case "DELETE":
		req, err = http.NewRequest("DELETE", url, bytes.NewReader(data.([]byte)))
		if err != nil {
			return nil, err
		}
	case "PUT":
		req, err = http.NewRequest("PUT", url, bytes.NewReader(data.([]byte)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	default:
		return nil, errors.New("unsupported method")
	}

	if err != nil {
		return nil, err
	}

	// 统一使用 Do 方法执行请求
	httpResp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return httpResp, nil
}

// Read Port from config file
func ReadPortFromConfig() (string, error) {
	defaultPort := "8088"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn(fmt.Sprintf("Error reading config file, %s, using default port %s", err, defaultPort))
		return defaultPort, fmt.Errorf("Error reading config file, %s, using default port %s", err, defaultPort)
	} else {
		portStr := viper.GetString("port.ha-api")
		_, err := strconv.Atoi(portStr)
		if err != nil {
			slog.Warn("Port in config is not a number, using default port 8088")
			return defaultPort, fmt.Errorf("port in config is not a number, using default port %s", defaultPort)
		}
		return portStr, nil
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

// 实现python的splitlines方法
func SplitLinesScanner(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func HostExists(hostName string) bool {
	_, err := RunCommand(fmt.Sprintf(CmdHostExists, hostName))
	return err == nil
}

func GenerateRemoteRequestURL(node string, uri string) string {
	if strings.HasPrefix(uri, "/remote") {
		return "https://" + node + ":" + port + uri
	}
	return "https://" + node + ":" + port + "/remote" + uri
}

func CopyFile(src string, dest string) error {
	cmd := exec.Command("/usr/bin/cp", src, dest)
	return cmd.Run()
}
