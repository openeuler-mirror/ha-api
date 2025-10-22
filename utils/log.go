/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */
 package utils

 import (
	 "context"
	 "io"
	 "log/slog"
	 "os"
	 "runtime"
	 "strconv"
	 "strings"
 
	 "github.com/spf13/viper"
	 "gopkg.in/natefinch/lumberjack.v2"
 )
 
 // 日志级别映射
 var levelMap = map[string]slog.Level{
	 "debug": slog.LevelDebug,
	 "info":  slog.LevelInfo,
	 "warn":  slog.LevelWarn,
	 "error": slog.LevelError,
 }
 
 // 初始化日志系统
 func SetupLogger() *slog.Logger {
	 // 加载配置
	 viper.SetConfigName("config")
	 viper.SetConfigType("yaml")
	 viper.AddConfigPath(".")
 
	 if err := viper.ReadInConfig(); err != nil {
		 panic("加载配置文件失败: " + err.Error())
	 }
 
	 // 获取日志配置
	 logCfg := viper.Sub("logging")
	 if logCfg == nil {
		 panic("配置文件中缺少 logging 部分")
	 }
	 var writers []io.Writer
	 if viper.GetBool("logging.output.file") {
		 logWriter := &lumberjack.Logger{
			 Filename:   "/var/log/ha-api.log",
			 MaxSize:    viper.GetInt("logging.rotation.max_size"),
			 MaxBackups: viper.GetInt("logging.rotation.max_backups"),
			 MaxAge:     viper.GetInt("logging.rotation.max_age"),
			 Compress:   viper.GetBool("logging.rotation.compress"),
			 LocalTime:  viper.GetBool("logging.rotation.local_time"),
		 }
		 writers = append(writers, logWriter)
	 }
 
	 // 取消控制台输出（开发模式强制开启）
	 consoleEnabled := viper.GetBool("logging.output.console")
	 if consoleEnabled {
		 writers = append(writers, os.Stdout)
	 }
 
	 // 如果没有激活任何输出目标，则强制开启控制台
	 if len(writers) == 0 {
		 writers = append(writers, os.Stdout)
		 slog.Warn("未激活任何日志输出目标，默认输出到控制台")
	 }
 
	 multiWriter := io.MultiWriter(writers...)
	 // 获取日志级别
	 levelStr := strings.ToLower(logCfg.GetString("level"))
	 level, ok := levelMap[levelStr]
	 if !ok {
		 level = slog.LevelInfo // 默认级别
	 }
 
	 // 创建自定义日志处理器
	 handler := &customHandler{
		 handler: slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			 Level: level,
		 }),
	 }
 
	 return slog.New(handler)
 }