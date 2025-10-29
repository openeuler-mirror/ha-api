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

// 自定义处理器添加错误堆栈
type customHandler struct {
	handler slog.Handler
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	// 错误级别时添加堆栈跟踪
	// if r.Level == slog.LevelError {
	// 	buf := make([]byte, 4096)
	// 	n := runtime.Stack(buf, false)
	// 	r.AddAttrs(slog.String("stack", string(buf[:n])))
	// }
	return h.handler.Handle(ctx, r)
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return &customHandler{handler: h.handler.WithGroup(name)}
}

// 获取调用位置信息
func CallerInfo() slog.Attr {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return slog.String("caller", "unknown")
	}

	// 简化文件路径
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		file = file[idx+1:]
	}

	return slog.String("caller", file+":"+strconv.Itoa(line))
}

// 简化日志记录函数
func Info(msg string, args ...any) {
	slog.Info(msg, append([]any{CallerInfo()}, args...)...)
	// slog.Info(msg, args)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, append([]any{CallerInfo()}, args...)...)
}

func Error(err error, msg string, args ...any) {
	slog.Error(msg, append([]any{CallerInfo(), "error", err}, args...)...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, append([]any{CallerInfo()}, args...)...)
}
