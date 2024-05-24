package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerField struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type zaplog struct {
	zap    *zap.Logger
	al     *zap.AtomicLevel
	fields []zapcore.Field
}

func NewLogger(cfg *config.LoggerConfig, env string) Logger {
	var writer io.Writer

	level := toZapLevel(Level(cfg.LogLevel))

	var zapOptions []zap.Option
	if env == constant.Dev.String() {
		writer = os.Stdout
		zapOptions = append(zapOptions, zap.Development())
	} else {
		writer = NewRotateWriter(cfg)
	}
	zapOptions = append(zapOptions, zap.AddCaller(),
		zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))

	zl := New(writer, level, zapOptions...)
	return zl
}

func New(out io.Writer, level zapcore.Level, opts ...zap.Option) *zaplog {
	if out == nil {
		out = os.Stdout
	}
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeName = zapcore.FullNameEncoder

	al := zap.NewAtomicLevelAt(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encodeConfig),
		zapcore.AddSync(out),
		al,
	)

	return &zaplog{zap: zap.New(core, opts...), al: &al}
}

func toZapLevel(level Level) zapcore.Level {
	var logLevel zapcore.Level
	switch level {
	case DebugLevel:
		logLevel = zap.DebugLevel
		break
	case InfoLevel:
		logLevel = zap.InfoLevel
		break
	case WarnLevel:
		logLevel = zap.WarnLevel
		break
	case ErrorLevel:
		logLevel = zap.ErrorLevel
		break
	case PanicLevel:
		logLevel = zap.PanicLevel
		break
	case DPanicLevel:
		logLevel = zap.DPanicLevel
		break
	case FatalLevel:
		logLevel = zap.FatalLevel
		break
	default:
		logLevel = zap.InfoLevel
		break
	}
	return logLevel
}

func (zl *zaplog) Log(level Level, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	switch level {
	case InfoLevel:
		zl.zap.Info(msg, zl.fields...)
	case DebugLevel:
		zl.zap.Debug(msg, zl.fields...)
	case WarnLevel:
		zl.zap.Warn(msg, zl.fields...)
	case ErrorLevel:
		zl.zap.Error(msg, zl.fields...)
	case PanicLevel:
		zl.zap.Panic(msg, zl.fields...)
	case DPanicLevel:
		zl.zap.DPanic(msg, zl.fields...)
	case FatalLevel:
		zl.zap.Fatal(msg, zl.fields...)
	}
	return nil
}

func (zl *zaplog) Info(format string, args ...any) error {
	return zl.Log(InfoLevel, format, args...)
}

func (zl *zaplog) Debug(format string, args ...any) error {
	return zl.Log(DebugLevel, format, args...)
}

func (zl *zaplog) Warn(format string, args ...any) error {
	return zl.Log(WarnLevel, format, args...)
}

func (zl *zaplog) Error(format string, args ...any) error {
	return zl.Log(ErrorLevel, format, args...)
}

func (zl *zaplog) Panic(format string, args ...any) error {
	return zl.Log(PanicLevel, format, args...)
}

func (zl *zaplog) DPanic(format string, args ...any) error {
	return zl.Log(DPanicLevel, format, args...)
}

func (zl *zaplog) Fatal(format string, args ...any) error {
	return zl.Log(FatalLevel, format, args...)
}

func (zl *zaplog) WithFields(fields ...zapcore.Field) {
	zl.fields = append(zl.fields, fields...)
}

func (zl *zaplog) Sync() error {
	return zl.zap.Sync()
}

func (zl *zaplog) String() string {
	return "zap"
}
