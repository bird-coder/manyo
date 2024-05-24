/*
 * @Author: yujiajie
 * @Date: 2024-05-23 17:52:21
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-24 09:18:27
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

func Infof(format string, args ...interface{}) {
	global.Log(InfoLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	global.Log(DebugLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	global.Log(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	global.Log(ErrorLevel, format, args...)
}

func Panicf(format string, args ...interface{}) {
	global.Log(PanicLevel, format, args...)
}

func DPanicf(format string, args ...interface{}) {
	global.Log(DPanicLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	global.Log(FatalLevel, format, args...)
	os.Exit(1)
}

func String() string {
	return global.String()
}

func Sync() {
	global.Sync()
}
