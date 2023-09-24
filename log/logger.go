package log

import (
	"fmt"
	"time"
)

func Debug(format string, args ...any) {
	sendMessage(LevelDebug, format, args...)
}

func Info(format string, args ...any) {
	sendMessage(LevelInfo, format, args...)
}

func Warn(format string, args ...any) {
	sendMessage(LevelWarning, format, args...)
}

func Error(format string, args ...any) {
	sendMessage(LevelError, format, args...)
}

func sendMessage(level logLevel, format string, args ...any) {
	if level < minLevel {
		return
	}

	fmt.Printf(
		"[%s] %s => %s\n",
		level.String(),
		time.Now().Format(time.RFC3339),
		fmt.Sprintf(format, args...),
	)
}
