package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// levelFatal is a custom log level above Error, used for unrecoverable failures.
const levelFatal = slog.Level(12)

// Logger wraps an slog handler along with service-level configuration.
type Logger struct {
	handler     slog.Handler
	serviceName string
	level       Level
	exitFunc    func(int)
}

var defaultLogger *Logger

func init() {
	defaultLogger = newLogger(Config{ServiceName: "app", Level: LevelInfo}, os.Stdout)
}

// New returns a configured Logger instance for use with dependency injection.
func New(cfg Config) *Logger {
	return newLogger(cfg, os.Stdout)
}

// Setup initializes the global logger with the provided configuration.
// Must be called at application startup before any log calls.
func Setup(cfg Config) {
	defaultLogger = newLogger(cfg, os.Stdout)
}

func newLogger(cfg Config, w io.Writer) *Logger {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "app"
	}
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: toSlogLevel(cfg.Level),
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "timestamp"
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			case slog.MessageKey:
				a.Key = "message"
			case slog.LevelKey:
				if level, ok := a.Value.Any().(slog.Level); ok && level == levelFatal {
					a.Value = slog.StringValue("FATAL")
				}
			}
			return a
		},
	})
	return &Logger{
		handler:     h,
		serviceName: cfg.ServiceName,
		level:       cfg.Level,
		exitFunc:    os.Exit,
	}
}

func (l *Logger) log(ctx context.Context, level slog.Level, src, msg string, err error) {
	if !l.handler.Enabled(ctx, level) {
		return
	}

	traceID, userID := extractFromContext(ctx)

	attrs := []slog.Attr{
		slog.String("service", l.serviceName),
		slog.String("src", src),
	}
	if traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	r := slog.NewRecord(time.Now(), level, msg, 0)
	r.AddAttrs(attrs...)
	_ = l.handler.Handle(ctx, r)
}

// Debug logs a message at DEBUG level.
// src identifies the call site (e.g. "UserService.Create").
func (l *Logger) Debug(ctx context.Context, src, msg string) {
	l.log(ctx, slog.LevelDebug, src, msg, nil)
}

// Info logs a message at INFO level.
func (l *Logger) Info(ctx context.Context, src, msg string) {
	l.log(ctx, slog.LevelInfo, src, msg, nil)
}

// Warn logs a message at WARN level.
func (l *Logger) Warn(ctx context.Context, src, msg string) {
	l.log(ctx, slog.LevelWarn, src, msg, nil)
}

// Error logs a message at ERROR level. err is included as the "error" field.
func (l *Logger) Error(ctx context.Context, src, msg string, err error) {
	l.log(ctx, slog.LevelError, src, msg, err)
}

// Fatal logs a message at FATAL level and calls os.Exit(1).
func (l *Logger) Fatal(ctx context.Context, src, msg string, err error) {
	l.log(ctx, levelFatal, src, msg, err)
	l.exitFunc(1)
}

// Debugf logs a formatted message at DEBUG level.
func (l *Logger) Debugf(ctx context.Context, src, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, src, fmt.Sprintf(msg, args...), nil)
}

// Infof logs a formatted message at INFO level.
func (l *Logger) Infof(ctx context.Context, src, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, src, fmt.Sprintf(msg, args...), nil)
}

// Warnf logs a formatted message at WARN level.
func (l *Logger) Warnf(ctx context.Context, src, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, src, fmt.Sprintf(msg, args...), nil)
}

// Errorf logs a formatted message at ERROR level. err is included as the "error" field.
func (l *Logger) Errorf(ctx context.Context, src, msg string, err error, args ...any) {
	l.log(ctx, slog.LevelError, src, fmt.Sprintf(msg, args...), err)
}

// Fatalf logs a formatted message at FATAL level and calls os.Exit(1).
func (l *Logger) Fatalf(ctx context.Context, src, msg string, err error, args ...any) {
	l.log(ctx, levelFatal, src, fmt.Sprintf(msg, args...), err)
	l.exitFunc(1)
}

// The package-level functions below delegate to the global logger initialized by Setup.

// Debug logs a message at DEBUG level using the global logger.
func Debug(ctx context.Context, src, msg string) {
	defaultLogger.Debug(ctx, src, msg)
}

// Info logs a message at INFO level using the global logger.
func Info(ctx context.Context, src, msg string) {
	defaultLogger.Info(ctx, src, msg)
}

// Warn logs a message at WARN level using the global logger.
func Warn(ctx context.Context, src, msg string) {
	defaultLogger.Warn(ctx, src, msg)
}

// Error logs a message at ERROR level using the global logger.
func Error(ctx context.Context, src, msg string, err error) {
	defaultLogger.Error(ctx, src, msg, err)
}

// Fatal logs a message at FATAL level using the global logger and calls os.Exit(1).
func Fatal(ctx context.Context, src, msg string, err error) {
	defaultLogger.Fatal(ctx, src, msg, err)
}

// Debugf logs a formatted message at DEBUG level using the global logger.
func Debugf(ctx context.Context, src, msg string, args ...any) {
	defaultLogger.Debugf(ctx, src, msg, args...)
}

// Infof logs a formatted message at INFO level using the global logger.
func Infof(ctx context.Context, src, msg string, args ...any) {
	defaultLogger.Infof(ctx, src, msg, args...)
}

// Warnf logs a formatted message at WARN level using the global logger.
func Warnf(ctx context.Context, src, msg string, args ...any) {
	defaultLogger.Warnf(ctx, src, msg, args...)
}

// Errorf logs a formatted message at ERROR level using the global logger.
func Errorf(ctx context.Context, src, msg string, err error, args ...any) {
	defaultLogger.Errorf(ctx, src, msg, err, args...)
}

// Fatalf logs a formatted message at FATAL level using the global logger and calls os.Exit(1).
func Fatalf(ctx context.Context, src, msg string, err error, args ...any) {
	defaultLogger.Fatalf(ctx, src, msg, err, args...)
}
