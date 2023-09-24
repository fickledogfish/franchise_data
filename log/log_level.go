package log

import "fmt"

type logLevel uint

const (
	LevelDebug logLevel = iota
	LevelInfo
	LevelWarning
	LevelError
)

func (self logLevel) String() string {
	switch self {
	case LevelDebug:
		return "DEBUG"

	case LevelInfo:
		return "INFO"

	case LevelWarning:
		return "WARN"

	case LevelError:
		return "ERR"

	default:
		panic(fmt.Sprintf("Unknown log level: %d", self))
	}
}

var (
	minLevel = LevelError
)

func SetMinLevel(newLevel logLevel) {
	minLevel = newLevel
}
