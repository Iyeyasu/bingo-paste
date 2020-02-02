package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	logger = NewLogger()
)

// Logger is used to log message.
type Logger struct {
	Out      io.Writer
	Level    Level
	ExitFunc func(int)
}

// NewLogger creates a new Logger.
func NewLogger() *Logger {
	logger := new(Logger)
	logger.Out = os.Stderr
	logger.Level = InfoLevel
	logger.ExitFunc = os.Exit
	return logger
}

// IsLevelEnabled returns true if the given logging level is enabled.
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.Level >= level
}

// Log logs the given arguments with no space between each argument.
func (logger *Logger) Log(level Level, args ...interface{}) {
	if !logger.IsLevelEnabled(level) {
		return
	}

	levelLabel := getLevelLabel(level)
	timeLabel := time.Now().Format(time.RFC850)
	msg := fmt.Sprint(args...)
	serialized := fmt.Sprintf("[%s] %s - %s\n", levelLabel, timeLabel, msg)

	if _, err := logger.Out.Write([]byte(serialized)); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to log:", err.Error())
	}

	if level == PanicLevel {
		panic(msg)
	} else if level == FatalLevel {
		logger.ExitFunc(1)
	}
}

// Logf logs the given arguments using the given format string.
func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.Log(level, fmt.Sprintf(format, args...))
	}
}

// Logln logs the given arguments with space between each argument.
func (logger *Logger) Logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.Log(level, fmt.Sprintln(args...))
	}
}

func getLevelLabel(level Level) string {
	switch level {
	case PanicLevel:
		return "PNIC"
	case FatalLevel:
		return "FATA"
	case ErrorLevel:
		return "ERRO"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DBUG"
	case TraceLevel:
		return "TRCE"
	default:
		return "?????"
	}
}
