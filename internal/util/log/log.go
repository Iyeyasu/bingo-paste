package log

// GetLevel returns the active logging level.
func GetLevel() Level {
	return logger.Level
}

// SetLevel sets the active logging level.
func SetLevel(level Level) {
	logger.Level = level
}

func Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

func Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
}

func Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Logf(PanicLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Logf(FatalLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

func Panicln(args ...interface{}) {
	logger.Logln(PanicLevel, args...)
}

func Fatalln(args ...interface{}) {
	logger.Logln(FatalLevel, args...)
}

func Errorln(args ...interface{}) {
	logger.Logln(ErrorLevel, args...)
}

func Warnln(args ...interface{}) {
	logger.Logln(WarnLevel, args...)
}

func Infoln(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func Debugln(args ...interface{}) {
	logger.Logln(DebugLevel, args...)
}

func Traceln(args ...interface{}) {
	logger.Logln(TraceLevel, args...)
}
