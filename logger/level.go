/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-30 21:34:50
 * @LastEditTime: 2023-09-30 21:46:31
 * @LastEditors: yuanshisan
 */
package logger

type Level int8

const (
	InfoLevel Level = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	DPanicLevel
	FatalLevel
)

func Info(format string, args ...interface{}) {
	logger.Log(InfoLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Log(DebugLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Log(WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Log(ErrorLevel, format, args...)
}

func Panic(format string, args ...interface{}) {
	logger.Log(PanicLevel, format, args...)
}

func DPanic(format string, args ...interface{}) {
	logger.Log(DPanicLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Log(FatalLevel, format, args...)
}
