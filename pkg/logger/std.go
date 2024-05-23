/*
 * @Author: yujiajie
 * @Date: 2024-05-15 18:05:20
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:46:30
 * @FilePath: /manyo/pkg/logger/std.go
 * @Description:
 */
package logger

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap/zapcore"
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

func (l *stdLogger) Log(level Level, msg string, fields ...zapcore.Field) error {
	if l.isDiscard || len(msg) == 0 {
		return nil
	}
	buf := l.pool.Get().(*bytes.Buffer)
	defer l.pool.Put(buf)

	defer buf.Reset()
	buf.WriteString(string(level) + " ")
	buf.WriteString(msg)
	buf.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write(buf.Bytes())
	return err
}

func (l *stdLogger) Info(msg string, fields ...zapcore.Field) error {
	return l.Log(InfoLevel, msg, fields...)
}

func (l *stdLogger) Debug(msg string, fields ...zapcore.Field) error {
	return l.Log(DebugLevel, msg, fields...)
}

func (l *stdLogger) Warn(msg string, fields ...zapcore.Field) error {
	return l.Log(WarnLevel, msg, fields...)
}

func (l *stdLogger) Error(msg string, fields ...zapcore.Field) error {
	return l.Log(ErrorLevel, msg, fields...)
}

func (l *stdLogger) Panic(msg string, fields ...zapcore.Field) error {
	return l.Log(PanicLevel, msg, fields...)
}

func (l *stdLogger) DPanic(msg string, fields ...zapcore.Field) error {
	return l.Log(DPanicLevel, msg, fields...)
}

func (l *stdLogger) Fatal(msg string, fields ...zapcore.Field) error {
	return l.Log(FatalLevel, msg, fields...)
}

func (l *stdLogger) Logf(level Level, format string, args ...any) error {
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

func (l *stdLogger) Infof(format string, args ...any) error {
	return l.Logf(InfoLevel, format, args...)
}

func (l *stdLogger) Debugf(format string, args ...any) error {
	return l.Logf(DebugLevel, format, args...)
}

func (l *stdLogger) Warnf(format string, args ...any) error {
	return l.Logf(WarnLevel, format, args...)
}

func (l *stdLogger) Errorf(format string, args ...any) error {
	return l.Logf(ErrorLevel, format, args...)
}

func (l *stdLogger) Panicf(format string, args ...any) error {
	return l.Logf(PanicLevel, format, args...)
}

func (l *stdLogger) DPanicf(format string, args ...any) error {
	return l.Logf(DPanicLevel, format, args...)
}

func (l *stdLogger) Fatalf(format string, args ...any) error {
	return l.Logf(FatalLevel, format, args...)
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
