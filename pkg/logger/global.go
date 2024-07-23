/*
 * @Author: yujiajie
 * @Date: 2024-05-23 17:52:21
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-31 14:40:59
 * @FilePath: /manyo/pkg/logger/global.go
 * @Description:
 */
package logger

import (
	"os"
	"sync"
)

var global = &loggerApp{}

type loggerApp struct {
	mu sync.Mutex
	Logger
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerApp) SetLogger(in Logger) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Logger = in
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

func GetLogger() Logger {
	return global
}

func Info(format string, args ...interface{}) {
	global.Log(InfoLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	global.Log(DebugLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	global.Log(WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	global.Log(ErrorLevel, format, args...)
}

func Panic(format string, args ...interface{}) {
	global.Log(PanicLevel, format, args...)
}

func DPanic(format string, args ...interface{}) {
	global.Log(DPanicLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	global.Log(FatalLevel, format, args...)
	os.Exit(1)
}

func String() string {
	return global.String()
}

func Sync() {
	global.Sync()
}
