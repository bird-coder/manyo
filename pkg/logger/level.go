/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-30 21:34:50
 * @LastEditTime: 2024-05-23 17:52:03
 * @LastEditors: yujiajie
 */
package logger

type Level string

const (
	InfoLevel   Level = "info"
	DebugLevel  Level = "debug"
	WarnLevel   Level = "warn"
	ErrorLevel  Level = "error"
	PanicLevel  Level = "panic"
	DPanicLevel Level = "dpanic"
	FatalLevel  Level = "fatal"
)
