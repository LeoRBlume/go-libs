package logger

import "log/slog"

// Level represents the minimum severity for log output.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Config holds the configuration used to initialize a Logger.
type Config struct {
	ServiceName string // included in every log entry as the "service" field
	Level       Level  // minimum level; entries below this are discarded
}

func toSlogLevel(l Level) slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
