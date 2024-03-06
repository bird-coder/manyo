/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-05 22:25:45
 * @LastEditTime: 2024-03-05 22:28:26
 * @LastEditors: yujiajie
 */
package logger

import (
	"io"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewRotateWriter(cfg *LoggerConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   strings.ToLower(cfg.Compress) == "true",
	}
}
