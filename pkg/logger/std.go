/*
 * @Author: yujiajie
 * @Date: 2024-05-15 18:05:20
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-24 09:16:23
 * @FilePath: /manyo/pkg/logger/std.go
 * @Description:
 */
package logger

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	w         io.Writer
	isDiscard bool
	mu        sync.Mutex
	pool      *sync.Pool
}

func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		w:         w,
		isDiscard: w == io.Discard,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (l *stdLogger) Log(level Level, format string, args ...any) error {
	if l.isDiscard || len(format) == 0 {
		return nil
	}
	buf := l.pool.Get().(*bytes.Buffer)
	defer l.pool.Put(buf)

	defer buf.Reset()
	buf.WriteString(string(level) + " ")
	if _, err := fmt.Fprintf(buf, format, args...); err != nil {
		return err
	}
	buf.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write(buf.Bytes())
	return err
}

func (l *stdLogger) Info(format string, args ...any) error {
	return l.Log(InfoLevel, format, args...)
}

func (l *stdLogger) Debug(format string, args ...any) error {
	return l.Log(DebugLevel, format, args...)
}

func (l *stdLogger) Warn(format string, args ...any) error {
	return l.Log(WarnLevel, format, args...)
}

func (l *stdLogger) Error(format string, args ...any) error {
	return l.Log(ErrorLevel, format, args...)
}

func (l *stdLogger) Panic(format string, args ...any) error {
	return l.Log(PanicLevel, format, args...)
}

func (l *stdLogger) DPanic(format string, args ...any) error {
	return l.Log(DPanicLevel, format, args...)
}

func (l *stdLogger) Fatal(format string, args ...any) error {
	return l.Log(FatalLevel, format, args...)
}

func (l *stdLogger) Close() error {
	return nil
}

func (l *stdLogger) String() string {
	return "stdout"
}

func (l *stdLogger) Sync() error {
	return nil
}
