package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bird-coder/manyo/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger Logger
)

type LoggerConfig struct {
	LogLevel   string                  `json:"level"`
	LogPath    string                  `json:"logpath"`
	MaxSize    int                     `json:"maxsize"`
	MaxAge     int                     `json:"age"`
	MaxBackups int                     `json:"backups"`
	Compress   string                  `json:"compress"`
	ConfigKey  *LoggerConfigKey        `json:"configKey"`
	Fields     map[string]*LoggerField `json:"fields"`
}

type LoggerConfigKey struct {
	MessageKey    string `json:"message"`
	LevelKey      string `json:"level"`
	TimeKey       string `json:"time"`
	NameKey       string `json:"name"`
	CallerKey     string `json:"caller"`
	StacktraceKey string `json:"stacktrace"`
}

type LoggerField struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type zaplog struct {
	zap *zap.Logger
	al  *zap.AtomicLevel
}

func NewLogger(logpath string, env string) Logger {
	cfg := &LoggerConfig{
		LogPath:    "",
		LogLevel:   "debug",
		MaxSize:    128,
		MaxAge:     7,
		MaxBackups: 30,
		Compress:   "false",
	}
	writer := NewRotateWriter(cfg)
	level := toZapLevel(Level(cfg.LogLevel))

	var zapOptions []zap.Option
	if env == constant.Dev.String() {
		zapOptions = append(zapOptions, zap.Development())
	}
	zapOptions = append(zapOptions, zap.AddCaller(),
		zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))

	zl := New(writer, level, zapOptions...)
	zl.zap.Info("init logger")

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

func (l *zaplog) formatFields(cfg *LoggerConfig) zap.Option {
	fields := []zapcore.Field{}
	for _, configField := range cfg.Fields {
		fields = append(fields, zap.String(configField.Key, configField.Val))
	}
	return zap.Fields(fields...)
}

func formatConfig(configMap map[string]interface{}) *LoggerConfig {
	data, err := json.Marshal(configMap)
	if err != nil {
		fmt.Fprint(os.Stderr, "load logger config error, error: json marshal failed\n")
		os.Exit(1)
	}
	config := &LoggerConfig{
		LogLevel:   "debug",
		MaxSize:    128,
		MaxAge:     7,
		MaxBackups: 30,
		Compress:   "false",
		ConfigKey: &LoggerConfigKey{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "file",
			StacktraceKey: "stacktrace",
		},
	}
	if err := json.Unmarshal(data, config); err != nil {
		fmt.Fprintf(os.Stderr, "load logger config error, error: json unmarshal failed, data: %s\n", string(data))
		os.Exit(1)
	}
	return config
}

func (zl *zaplog) Log(level Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	switch level {
	case InfoLevel:
		zl.zap.Info(msg)
	case DebugLevel:
		zl.zap.Debug(msg)
	case WarnLevel:
		zl.zap.Warn(msg)
	case ErrorLevel:
		zl.zap.Error(msg)
	case PanicLevel:
		zl.zap.Panic(msg)
	case DPanicLevel:
		zl.zap.DPanic(msg)
	case FatalLevel:
		zl.zap.Fatal(msg)
	}
}

func (zl *zaplog) Sync() {
	zl.zap.Sync()
}

func (zl *zaplog) String() string {
	return "zap"
}
