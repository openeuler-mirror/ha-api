/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bixiaoyan <bixiaoyan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gitee.com/openeuler/ha-api/settings"
)

func TestGetCorosyncConfig(t *testing.T) {
	result, err := GetCorosyncConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func TestGetCorosyncConfig2(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		configContent string
		want          CorosyncConfig
		wantErr       bool
	}{
		{
			name: "configuration with qdevice",
			configContent: `totem {
								transport: knet
								token: 8000
								crypto_cipher: aes256
								crypto_hash: sha256
								cluster_uuid: d2caae471703440cbac34541cee7b396
								version: 2
								cluster_name: mycluster
							}

							nodelist {
								node {
									nodeid: 1
									ring0_addr: node1
								}
								node {
									nodeid: 2
									ring0_addr: node2
								}
							}

							quorum {
								provider: corosync_votequorum
								device {
									model: net
									votes: 1
									net {
										algorithm: ffsplit
										host: host105
									}
								}
							}

							logging {
								to_syslog: yes
								timestamp: on
								to_logfile: yes
							}`,
			want: CorosyncConfig{
				Totem: map[string]string{
					"token":         "8000",
					"crypto_cipher": "aes256",
					"crypto_hash":   "sha256",
					"cluster_uuid":  "d2caae471703440cbac34541cee7b396",
					"version":       "2",
					"cluster_name":  "mycluster",
					"transport":     "knet",
				},
				NodeList: []map[string]string{
					{
						"nodeid":     "1",
						"ring0_addr": "node1",
					},
					{
						"nodeid":     "2",
						"ring0_addr": "node2",
					},
				},
				Quorum: map[string]interface{}{
					"provider": "corosync_votequorum",
					"device": map[string]interface{}{
						"model": "net",
						"votes": "1",
						"net": map[string]string{
							"algorithm": "ffsplit",
							"host":      "host105",
						},
					},
				},
				Logging: map[string]interface{}{
					"timestamp":  "on",
					"to_logfile": "yes",
					"to_syslog":  "yes",
				},
			},
			wantErr: false,
		},
		{
			name:          "empty",
			configContent: "",
			want:          CorosyncConfig{},
			wantErr:       false,
		},
		{
			name: "configuration withpit qdevice",
			configContent: `totem {
								transport: knet
								token: 8000
								crypto_cipher: aes256
								crypto_hash: sha256
								cluster_uuid: d2caae471703440cbac34541cee7b396
								version: 2
								cluster_name: mycluster
							}

							nodelist {
								node {
									nodeid: 1
									ring0_addr: node1
								}
								node {
									nodeid: 2
									ring0_addr: node2
								}
							}

							quorum {
								two_node: 1
    							provider: corosync_votequorum
							}

							logging {
								to_syslog: yes
								timestamp: on
								to_logfile: yes
							}`,
			want: CorosyncConfig{
				Totem: map[string]string{
					"token":         "8000",
					"crypto_cipher": "aes256",
					"crypto_hash":   "sha256",
					"cluster_uuid":  "d2caae471703440cbac34541cee7b396",
					"version":       "2",
					"cluster_name":  "mycluster",
					"transport":     "knet",
				},
				NodeList: []map[string]string{
					{
						"nodeid":     "1",
						"ring0_addr": "node1",
					},
					{
						"nodeid":     "2",
						"ring0_addr": "node2",
					},
				},
				Quorum: map[string]interface{}{
					"provider": "corosync_votequorum",
					"two_node": "1",
				},
				Logging: map[string]interface{}{
					"timestamp":  "on",
					"to_logfile": "yes",
					"to_syslog":  "yes",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时配置文件
			tmpFile := filepath.Join(tmpDir, "corosync.conf")
			err := os.WriteFile(tmpFile, []byte(tt.configContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			// 保存原始设置并临时修改
			originalFile := settings.CorosyncConfFile
			settings.CorosyncConfFile = tmpFile
			defer func() {
				settings.CorosyncConfFile = originalFile
			}()

			// 执行测试
			got, err := GetCorosyncConfig()

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCorosyncConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果不期望错误，检查结果
			if !tt.wantErr {
				// 比较 Totem
				if !reflect.DeepEqual(got.Totem, tt.want.Totem) {
					t.Errorf("GetCorosyncConfig() Totem = %v, want %v", got.Totem, tt.want.Totem)
				}

				// 比较 NodeList
				if !reflect.DeepEqual(got.NodeList, tt.want.NodeList) {
					t.Errorf("GetCorosyncConfig() NodeList = %v, want %v", got.NodeList, tt.want.NodeList)
				}

				// 比较 Quorum
				if !reflect.DeepEqual(got.Quorum, tt.want.Quorum) {
					t.Errorf("GetCorosyncConfig() Quorum = %v, want %v", got.Quorum, tt.want.Quorum)
				}

				// 比较 Logging
				if !reflect.DeepEqual(got.Logging, tt.want.Logging) {
					t.Errorf("GetCorosyncConfig() Logging = %v, want %v", got.Logging, tt.want.Logging)
				}
			}
		})
	}
}

func TestGetCorosyncpcConfig_FileNotFound(t *testing.T) {
	// 保存原始设置并临时修改为不存在的文件
	originalFile := settings.CorosyncConfFile
	settings.CorosyncConfFile = "/nonexistent/path/corosync.conf"
	defer func() {
		settings.CorosyncConfFile = originalFile
	}()

	_, err := GetCorosyncConfig()
	if err == nil {
		t.Error("GetCorosyncConfig() expected error for nonexistent file, got nil")
	}
}

// 测试边界情况
func TestGetCorosyncConfig_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	edgeCases := []struct {
		name    string
		content string
	}{
		{
			name: "nested brackets in values",
			content: `totem {
    config: value{with}brackets
}`,
		},
		{
			name: "multiple spaces",
			content: `totem    {
    cluster_name:    test    
}`,
		},
		{
			name: "tabs instead of spaces",
			content: `totem	{
	cluster_name:	test
}`,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := filepath.Join(tmpDir, "corosync.conf")
			err := os.WriteFile(tmpFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			originalFile := settings.CorosyncConfFile
			settings.CorosyncConfFile = tmpFile
			defer func() {
				settings.CorosyncConfFile = originalFile
			}()

			// 测试是否能够正常解析而不崩溃
			_, err = GetCorosyncConfig()
			if err != nil {
				t.Logf("GetCorosyncConfig() with edge case '%s' returned error (may be expected): %v", tc.name, err)
			}
		})
	}
}
