/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liupei <liupei@kylinos.cn>
 * Date: Tue Jun 23 15:54:28 2026 +0800
 */

 package models

 // ==================== 配置定义（唯一数据源） ====================
 
 // CategoryConfig 存储所有分类及其子项配置
 var CategoryConfig = map[string]Category{
	 "DISK": {
		 SubItems:   []string{"LVM"},
		 DefaultLen: 1,
		 ResourceDefaults: map[string]map[string]interface{}{
			 "LVM": {
				 "id":         "",
				 "volgrpname": "",
				 "tag":        "",
			 },
		 },
	 },
	 "FileSystem": {
		 SubItems:   []string{"Filesystem"},
		 DefaultLen: 1,
		 ResourceDefaults: map[string]map[string]interface{}{
			 "Filesystem": {
				 "id":        "",
				 "device":    "",
				 "directory": "",
				 "fstype":    "",
				 "options":   "",
			 },
		 },
	 },
	 "Database": {
		 SubItems:   []string{"DMDB8"},
		 DefaultLen: 1,
		 ResourceDefaults: map[string]map[string]interface{}{
			 "DMDB8": {
				 "id":          "",
				 "datadir":     "/opt/dmdbms/",
				 "instancedir": "/opt/dmdbms/data/DAMENG",
			 },
		 },
	 },
	 "Middleware": {
		 SubItems:   []string{"TongWeb8", "BES9.5"},
		 DefaultLen: 1,
		 ResourceDefaults: map[string]map[string]interface{}{
			 "TongWeb8": {
				 "id":           "",
				 "tongweb_path": "",
			 },
			 "BES9.5": {
				 "id":              "",
				 "BES_HOME":        "/opt/BES",
				 "iastoolUser":     "",
				 "iastoolPassword": "",
				 "JAVA_HOME":       "",
			 },
		 },
	 },
	 "VIP": {
		 SubItems:   []string{"IPaddr2", "IPaddr"},
		 DefaultLen: 1,
		 ResourceDefaults: map[string]map[string]interface{}{
			 "IPaddr2": {
				 "id":           "",
				 "ip":           "",
				 "nic":          "",
				 "cidr_netmask": "",
				 "broadcast":    "",
			 },
			 "IPaddr": {
				 "id":           "",
				 "ip":           "",
				 "nic":          "eth0",
				 "cidr_netmask": "",
				 "broadcast":    "",
			 },
		 },
	 },
 }
 
 // Category 表示一个资源分类
 type Category struct {
	 SubItems         []string                          // 该分类下的资源类型列表
	 DefaultLen       int                               // 默认部署数量
	 ResourceDefaults map[string]map[string]interface{} // 各资源类型的属性默认值
 }
 
 // 主页显示内容
 func ResourceModelHomepageGet() string {
	 content := "this is resource model"
	 return content
 }
 
 func ResourceModelDeploy() map[string]Category {
	 return CategoryConfig
 } 