package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	commonUtil "github.com/bird-coder/manyo/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

var (
	DEV = "development"
	PRO = "production"
)

type LoggerConfigMap struct {
	LogLevel   string                  `json:"level"`
	Filename   string                  `json:"logpath"`
	MaxSize    string                  `json:"maxsize"`
	MaxAge     string                  `json:"age"`
	MaxBackups string                  `json:"backups"`
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

func NewLogger(filename string, env string) {
	configMap := commonUtil.LoadXmlConfig(filename)
	// checkLoggerConfig(configMap)
	config := formatConfig(configMap)
	maxSize, age, backups := config.formatIntParams()
	hook := lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    maxSize,
		MaxAge:     age,
		MaxBackups: backups,
		Compress:   strings.ToLower(config.Compress) == "true",
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     config.ConfigKey.MessageKey,
		LevelKey:       config.ConfigKey.LevelKey,
		TimeKey:        config.ConfigKey.TimeKey,
		NameKey:        config.ConfigKey.NameKey,
		CallerKey:      config.ConfigKey.CallerKey,
		StacktraceKey:  config.ConfigKey.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	atomicLevel := config.setLoggerLevel()
	var writes []zapcore.WriteSyncer
	var encoder zapcore.Encoder
	if env == DEV {
		writes = append(writes, zapcore.AddSync(os.Stdout))
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		writes = append(writes, zapcore.AddSync(&hook))
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)
	caller := zap.AddCaller()
	skip := zap.AddCallerSkip(1)
	stacktrace := zap.AddStacktrace(zap.WarnLevel)
	field := config.formatFields()

	var zapOptions []zap.Option
	if env == DEV {
		development := zap.Development()
		zapOptions = append(zapOptions, development)
	}
	zapOptions = append(zapOptions, caller, skip, stacktrace, field)

	logger = zap.New(core, zapOptions...)
	logger.Info("init logger")
}

func (config *LoggerConfigMap) setLoggerLevel() zap.AtomicLevel {
	var logLevel zapcore.Level
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		logLevel = zap.DebugLevel
		break
	case "info":
		logLevel = zap.InfoLevel
		break
	case "warn":
		logLevel = zap.WarnLevel
		break
	case "error":
		logLevel = zap.ErrorLevel
		break
	case "panic":
		logLevel = zap.PanicLevel
		break
	case "dpanic":
		logLevel = zap.DPanicLevel
		break
	case "fatal":
		logLevel = zap.FatalLevel
		break
	default:
		logLevel = zap.InfoLevel
		break
	}
	return zap.NewAtomicLevelAt(logLevel)
}

func (config *LoggerConfigMap) formatIntParams() (int, int, int) {
	errStr := "ParamError: wrong param value of logger config file, param: %s, data: %s\n"
	maxSize, err := strconv.Atoi(config.MaxSize)
	if err != nil {
		fmt.Printf(errStr, "maxsize", maxSize)
		os.Exit(1)
	}
	age, err := strconv.Atoi(config.MaxAge)
	if err != nil {
		fmt.Printf(errStr, "age", age)
		os.Exit(1)
	}
	backups, err := strconv.Atoi(config.MaxBackups)
	if err != nil {
		fmt.Printf(errStr, "backups", backups)
		os.Exit(1)
	}
	return maxSize, age, backups
}

func (config *LoggerConfigMap) formatFields() zap.Option {
	fields := []zapcore.Field{}
	for _, configField := range config.Fields {
		fields = append(fields, zap.String(configField.Key, configField.Val))
	}
	return zap.Fields(fields...)
}

func formatConfig(configMap map[string]interface{}) *LoggerConfigMap {
	data, err := json.Marshal(configMap)
	if err != nil {
		fmt.Fprint(os.Stderr, "load logger config error, error: json marshal failed\n")
		os.Exit(1)
	}
	config := &LoggerConfigMap{
		LogLevel:   "debug",
		MaxSize:    "128",
		MaxAge:     "7",
		MaxBackups: "30",
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

//该方法已废弃
func checkLoggerConfig(configMap map[string]interface{}) {
	isValid := true
Loop:
	for key, val := range configMap {
		if key == "fields" {
			if _, ok := val.(map[string]interface{}); !ok {
				isValid = false
				break Loop
			}
			for _, v := range val.(map[string]interface{}) {
				if _, ok := v.(map[string]interface{}); !ok {
					isValid = false
					break Loop
				}
				for _, v1 := range v.(map[string]interface{}) {
					if _, ok := v1.(string); !ok {
						isValid = false
						break Loop
					}
				}
			}
		} else if key == "configKey" {
			if _, ok := val.(map[string]interface{}); !ok {
				isValid = false
				break Loop
			}
			for _, v := range val.(map[string]interface{}) {
				if _, ok := v.(string); !ok {
					isValid = false
					break Loop
				}
			}
		} else {
			if _, ok := val.(string); !ok {
				isValid = false
				break Loop
			}
		}
	}
	if !isValid {
		fmt.Fprint(os.Stderr, "CheckLoggerConfig: Error: invalid params type of logger config file\n")
		os.Exit(1)
	}
}

func Sync() {
	logger.Sync()
}

func Info(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	logger.Panic(fmt.Sprintf(format, args...))
}

func DPanic(format string, args ...interface{}) {
	logger.DPanic(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, args...))
}
