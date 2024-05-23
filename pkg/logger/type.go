/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2023-09-30 21:16:54
 * @LastEditTime: 2024-05-23 17:48:18
 * @LastEditors: yujiajie
 */
package logger

import (
	"log"

	"go.uber.org/zap/zapcore"
)

var (
	DefaultLogger = NewStdLogger(log.Writer())
)

type Logger interface {
	Log(level Level, msg string, fields ...zapcore.Field) error
	Info(msg string, fields ...zapcore.Field) error
	Debug(msg string, fields ...zapcore.Field) error
	Warn(msg string, fields ...zapcore.Field) error
	Error(msg string, fields ...zapcore.Field) error
	Panic(msg string, fields ...zapcore.Field) error
	DPanic(msg string, fields ...zapcore.Field) error
	Fatal(msg string, fields ...zapcore.Field) error
	Logf(level Level, format string, args ...interface{}) error
	Infof(format string, args ...any) error
	Debugf(format string, args ...any) error
	Warnf(format string, args ...any) error
	Errorf(format string, args ...any) error
	Panicf(format string, args ...any) error
	DPanicf(format string, args ...any) error
	Fatalf(format string, args ...any) error
	String() string
	Sync() error
}
