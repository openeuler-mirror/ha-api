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
	"path/filepath"
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

func GrepHostsFile(host string) (int, string, error) {
	cmd := fmt.Sprintf(CmdHostSearch, ShellEscape(host))
	out, err := RunCommand(cmd)
	if err != nil {
		return 0, "", err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	count := len(lines)
	if lines[0] == "" { // 处理空输出情况
		count = 0
	}

	return count, string(out), nil
}

func HasSpace(s string) bool {
	matched, _ := regexp.MatchString(`\s`, s)
	return matched
}

/*
 *检查集群pcsd offfline的节点
 */
func CheckPcsdOfflineNodes() ([]string, error) {
	out, err := RunCommand(CmdPcsdStatus)
	if err != nil {
		return nil, err
	}
	var offlineNodes []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "Offline") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				node := strings.TrimSpace(parts[0])
				offlineNodes = append(offlineNodes, node)
			}
		}
	}
	return offlineNodes, nil
}

func IsClusterStarted() bool {
	_, err := RunCommand(CmdClusterStatus)
	return err == nil
}

func FindWithStart(s, substr string, start int) int {
	if start < 0 || start > len(s) {
		return -1 // 无效的起始位置
	}

	// 从 start 位置开始截取字符串
	slice := s[start:]
	index := strings.Index(slice, substr)

	if index == -1 {
		return -1 // 未找到
	}

	// 返回全局索引（start + 子串在切片中的位置）
	return start + index
}

func FindSubstring(s string, substr string, start int, end int) int {
	// 检查 start 和 end 是否有效
	if start < 0 || end > len(s) || start > end {
		return -1
	}

	// 切片字符串，从 start 到 end
	slice := s[start:end]

	// 在切片中查找子字符串
	index := strings.Index(slice, substr)

	// 如果找到子字符串，返回全局索引
	if index != -1 {
		return start + index
	}

	// 未找到，返回 -1
	return -1
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

func CopyKeyFiles(privateKeySrc, publicKeySrc, destDir string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}

	privateKeyDest := filepath.Join(destDir, filepath.Base(settings.RSA_PRIVATE_KEY))
	if err := CopyFile(privateKeySrc, privateKeyDest); err != nil {
		return fmt.Errorf("failed to copy private key: %v", err)
	}

	publicKeyDest := filepath.Join(destDir, filepath.Base(settings.RSA_PUBLIC_KEY))
	if err := CopyFile(publicKeySrc, publicKeyDest); err != nil {
		return fmt.Errorf("failed to copy public key: %v", err)
	}

	return nil
}
