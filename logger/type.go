package logger

type Logger interface {
	Init()
	Log(level Level, format string, args ...interface{})
	String() string
}

func Sync() {
	if zl, ok := logger.(*zaplog); ok {
		zl.Sync()
	}
}
