/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2023-09-30 21:16:54
 * @LastEditTime: 2024-03-05 22:45:08
 * @LastEditors: yujiajie
 */
package logger

type Logger interface {
	Log(level Level, format string, args ...interface{})
	String() string
}

func Sync() {
	if zl, ok := logger.(*zaplog); ok {
		zl.Sync()
	}
}
