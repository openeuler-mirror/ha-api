/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Fri Nov 15 16:03:14 2024 +0800
 */

package models

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gitee.com/openeuler/ha-api/settings"
)

// TestIsClusterExist 测试集群配置文件存在性检查
func TestIsClusterExist(t *testing.T) {
	// Mock配置文件路径
	originalPath := settings.CorosyncConfFile
	defer func() { settings.CorosyncConfFile = originalPath }()

	t.Run("file_exists", func(t *testing.T) {
		// 创建临时文件
		tmpDir := t.TempDir()
		confPath := filepath.Join(tmpDir, "corosync.conf")
		if err := os.WriteFile(confPath, []byte(""), 0644); err != nil {
			t.Fatal(err)
		}
		settings.CorosyncConfFile = confPath

		if !IsClusterExist() {
			t.Error("Expected true when file exists")
		}
	})

	t.Run("file_not_exists", func(t *testing.T) {
		settings.CorosyncConfFile = "/path/not/exist"
		if IsClusterExist() {
			t.Error("Expected false when file missing")
		}
	})
}

// TestGetNodeList 测试节点列表解析功能
func TestGetNodeList(t *testing.T) {
	// Mock配置文件路径
	originalPath := settings.CorosyncConfFile
	defer func() { settings.CorosyncConfFile = originalPath }()

	testCases := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "missing_file",
			content:  "",
			expected: []string{"error File /invalid/path doesn't exist!"},
		},
		{
			name: "no_nodelist_section",
			content: `quorum {
                provider: corosync_votequorum
            }`,
			expected: []string{"No nodelist found in the file."},
		},
		{
			name: "empty_nodelist",
			content: `nodelist {
            }`,
			expected: []string{},
		},
		{
			name: "valid_nodes",
			content: `nodelist {
                node { ring0_addr: node1 }
                node { ring0_addr: node2 }
            }
            quorum { }`,
			expected: []string{
				"node { ring0_addr: node1 }",
				"node { ring0_addr: node2 }",
			},
		},
		{
			name: "malformed_lines",
			content: `nodelist {
                # 注释行
                node { ring0_addr: node1 }
                invalid_line
                }
            quorum { }`,
			expected: []string{
				"# 注释行",
				"node { ring0_addr: node1 }",
				"invalid_line",
			},
		},
		{
			name: "early_quorum_break",
			content: `nodelist {
                node1
                quorum {
                    node2
                }
            }`,
			expected: []string{"node1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 特殊处理缺失文件用例
			if tc.name == "missing_file" {
				settings.CorosyncConfFile = "/invalid/path"
			} else {
				// 创建临时配置文件
				tmpDir := t.TempDir()
				confPath := filepath.Join(tmpDir, "corosync.conf")
				if err := os.WriteFile(confPath, []byte(tc.content), 0644); err != nil {
					t.Fatal(err)
				}
				settings.CorosyncConfFile = confPath
			}

			// 执行解析
			result := getNodeList()

			// 精确结果匹配
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("\nExpected: %#v\nGot     : %#v", tc.expected, result)
			}
		})
	}
}
