/*
 * Copyright (c) KylinSoft  Co., Ltd. 2027.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: xuxiaojuan <xuxiaojuan@kylinos.cn>
 * Date: Wed July 8 13:56:40 2026 +0800
 */

package models

import (
        "bufio"
        "fmt"
        "regexp"
        "strconv"
        "strings"
        "sync"

        "gitee.com/openeuler/ha-api/utils"
)

var allowedServiceNames = map[string]bool{
        "corosync": true,
        "pacemaker": true,
        "pcsd":     true,
        "ha-api":   true,
}

var (
        activeRegex   = regexp.MustCompile(`Active:\s+(.*?)\s+since\s+(.*?);\s+(.*)`)
        activeFallback = regexp.MustCompile(`Active:\s+(.*)`)
        pidRegex      = regexp.MustCompile(`Main PID:\s+(\d+)`)
        memoryRegex   = regexp.MustCompile(`Memory:\s+(.*)`)
        logRegex      = regexp.MustCompile(`^\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}`)
)

func getRunResult(cmd string) (string, error) {
        out, err := utils.RunCommand(cmd)
        if err != nil {
                return "", err
        }
        result := strings.TrimSpace(string(out))
        if idx := strings.IndexByte(result, '\n'); idx != -1 {
                result = result[:idx]
        }
        return result, nil
}

var constMap = map[string]string{
        "Kernel_release":          "uname -r",
        "Architecture":            "uname -m",
        "Distribution":            "cat /etc/os-release 2>/dev/null | grep '^PRETTY_NAME=' | cut -d= -f2 | tr -d '\"'",
        "Uptime":                  "uptime -p",
        "Pacemaker_version":       "cibadmin --version 2>&1 | head -1",
        "Corosync_version":        "/usr/sbin/corosync -v 2>&1 | head -1",
        "resource-agents_version": "grep 'Build version:' /usr/lib/ocf/resource.d/heartbeat/.ocf-shellfuncs",
        "Knet_version":            "rpm -q libknet1",
        "Ha-api_version":          "rpm -q ha-api",
        "Ha-web_version":          "rpm -q ha-web",
        "Glibc_version":           "rpm -q glibc",
}

type VmStat struct {
        Procs struct {
                R int // 运行队列进程数
                B int // 阻塞进程数
        }
        Memory struct { // 单位: KB
                Swpd  int // 已用交换空间
                Free  int // 空闲内存
                Buff  int // 缓冲区内存
                Cache int // 缓存内存
        }
        Swap struct { // 单位: KB/s
                Si int // 调入交换区
                So int // 调出交换区
        }
        IO struct {
                Bi int // 块设备读入块数/秒
                Bo int // 块设备写入块数/秒
                In int // 每秒中断次数
                Cs int // 每秒上下文切换次数
        }
        CPU struct { // 百分比
                Us float64 // 用户态CPU
                Sy float64 // 内核态CPU
                Id float64 // 空闲CPU
                Wa float64 // IO等待CPU
        }
}

func ParseVmstat(output string) (*VmStat, error) {
        scanner := bufio.NewScanner(strings.NewReader(output))
        stat := &VmStat{}

        for scanner.Scan() {
                line := strings.TrimSpace(scanner.Text())
                if line == "" || strings.HasPrefix(line, "procs") {
                        continue // 跳过标题行和空行
                }

                fields := strings.Fields(line)
                if len(fields) < 16 { // 基础模式字段数验证
                        continue
                }

                // 核心字段映射，任一字段解析失败则跳过该行
                var err error
                stat.Procs.R, err = strconv.Atoi(fields[0])
                if err != nil {
                        continue
                }
                stat.Procs.B, err = strconv.Atoi(fields[1])
                if err != nil {
                        continue
                }
                stat.Memory.Swpd, err = strconv.Atoi(fields[2])
                if err != nil {
                        continue
                }
                stat.Memory.Free, err = strconv.Atoi(fields[3])
                if err != nil {
                        continue
                }
                stat.Memory.Buff, err = strconv.Atoi(fields[4])
                if err != nil {
                        continue
                }
                stat.Memory.Cache, err = strconv.Atoi(fields[5])
                if err != nil {
                        continue
                }
                stat.Swap.Si, err = strconv.Atoi(fields[6])
                if err != nil {
                        continue
                }
                stat.Swap.So, err = strconv.Atoi(fields[7])
                if err != nil {
                        continue
                }
                stat.IO.Bi, err = strconv.Atoi(fields[8])
                if err != nil {
                        continue
                }
                stat.IO.Bo, err = strconv.Atoi(fields[9])
                if err != nil {
                        continue
                }
                stat.IO.In, err = strconv.Atoi(fields[10])
                if err != nil {
                        continue
                }
                stat.IO.Cs, err = strconv.Atoi(fields[11])
                if err != nil {
                        continue
                }

                stat.CPU.Us, err = strconv.ParseFloat(fields[12], 64)
                if err != nil {
                        continue
                }
                stat.CPU.Sy, err = strconv.ParseFloat(fields[13], 64)
                if err != nil {
                        continue
                }
                stat.CPU.Id, err = strconv.ParseFloat(fields[14], 64)
                if err != nil {
                        continue
                }
                stat.CPU.Wa, err = strconv.ParseFloat(fields[15], 64)
                if err != nil {
                        continue
                }
        }

        if err := scanner.Err(); err != nil {
                return nil, err
        }

        return stat, nil
}

func GetSystemInfo() map[string]interface{} {
        result := make(map[string]interface{})
        data := map[string]interface{}{}
        sysInfo := make(map[string]string)
        versionInfo := make(map[string]string)

        satcmd, err := utils.RunCommand("vmstat")
        if err != nil {
                result["action"] = false
                result["error"] = err.Error()
                return result
        }
        stat, err := ParseVmstat(string(satcmd))
        if err != nil {
                result["action"] = false
                result["error"] = err.Error()
                return result
        }
        data["stats"] = stat

        errors := make(map[string]string)
        for k, v := range constMap {
                output, err := getRunResult(v)
                if err != nil {
                        errors[k] = err.Error()
                        continue
                }
                if strings.Contains(k, "_version") {
                        versionInfo[k] = output
                } else {
                        sysInfo[k] = output
                }
        }

        data["sysInfo"] = sysInfo
        data["versionInfo"] = versionInfo
        result["action"] = true
        result["data"] = data
        if len(errors) > 0 {
                result["error"] = errors
        }
        return result
}

type SystemdUnitStatus struct {
        Name         string   `json:"name"`
        Active       string   `json:"active"`       // 运行状态
        RunningSince string   `json:"runningSince"` // 启动时间
        ProcessID    string   `json:"processID"`    // 主进程PID
        MemoryUsage  string   `json:"memoryUsage"`  // 内存占用
        Logs         []string `json:"logs"`         // 最近日志
}
func ParseSystemctlStatus(serviceName string) (*SystemdUnitStatus, error) {
        if !allowedServiceNames[serviceName] {
                return nil, fmt.Errorf("service name not allowed: %s", serviceName)
        }
        out, err := utils.RunCommandWithArgs("systemctl", "status", serviceName, "--no-pager")
        if err != nil {
                return nil, fmt.Errorf("systemctl status %s: %w", serviceName, err)
        }

        return parseStatusOutput(serviceName, string(out))
}
func parseStatusOutput(name, output string) (*SystemdUnitStatus, error) {
        status := &SystemdUnitStatus{Name: name}
        lines := strings.Split(output, "\n")

        var logs []string

        for _, line := range lines {
                line = strings.TrimSpace(line)
                switch {
                case strings.HasPrefix(line, "Active:"):
                        matches := activeRegex.FindStringSubmatch(line)
                        if len(matches) > 3 {
                                status.Active = matches[1]
                                status.RunningSince = matches[2]
                        } else {
                                matches = activeFallback.FindStringSubmatch(line)
                                if len(matches) > 1 {
                                        status.Active = matches[1]
                                }
                        }

                case strings.HasPrefix(line, "Main PID:"):
                        matches := pidRegex.FindStringSubmatch(line)
                        if len(matches) > 1 {
                                status.ProcessID = matches[1]
                        }

                case strings.HasPrefix(line, "Memory:"):
                        matches := memoryRegex.FindStringSubmatch(line)
                        if len(matches) > 1 {
                                status.MemoryUsage = matches[1]
                        }
                case logRegex.MatchString(line):
                        if strings.Contains(strings.ToLower(line), "error") || strings.Contains(strings.ToLower(line), "fail") {
                                logs = append(logs, line)
                        }
                }
        }

        status.Logs = logs
        if status.Active == "" {
                return nil, fmt.Errorf("failed to parse systemctl status for %s: no Active field found", name)
        }
        return status, nil
}
func GetServiceStatus() map[string]interface{} {
        result := make(map[string]interface{})
        services := []string{"corosync", "pacemaker", "pcsd", "ha-api"}
        data := make([]SystemdUnitStatus, len(services))

        var wg sync.WaitGroup
        for i, r := range services {
                wg.Add(1)
                go func(idx int, name string) {
                        defer wg.Done()
                        status, err := ParseSystemctlStatus(name)
                        if err != nil {
                                data[idx] = SystemdUnitStatus{Name: name, Active: "Unknown"}
                        } else {
                                data[idx] = *status
                        }
                }(i, r)
        }
        wg.Wait()

        result["action"] = true
        result["data"] = data
        return result
}
