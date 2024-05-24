/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2023-09-30 21:16:54
 * @LastEditTime: 2024-05-24 09:18:04
 * @LastEditors: yujiajie
 */
package logger

import (
	"log"
)

var (
	DefaultLogger = NewStdLogger(log.Writer())
)

type Logger interface {
	Log(level Level, format string, args ...interface{}) error
	Info(format string, args ...any) error
	Debug(format string, args ...any) error
	Warn(format string, args ...any) error
	Error(format string, args ...any) error
	Panic(format string, args ...any) error
	DPanic(format string, args ...any) error
	Fatal(format string, args ...any) error
	String() string
	Sync() error
}
