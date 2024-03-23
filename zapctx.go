package zapctx

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key struct{}

const (
	traceID = "trace_id"
)

var (
	// logLevel defaults to info and is used to control what level of messages should be logged.
	logLevel = zap.NewAtomicLevel()
	// hooks holds the functions to be added as hooks for the logger.
	hooks []func(entry zapcore.Entry) error
)

// SetLogLevel allows defining the log level to be used for log messages when using the context logger.
func SetLogLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}

// Hook defines the required signature for your log entries hooks.
type Hook func(entry zapcore.Entry) error

// AddHooks allows you to add zap hooks to be executed when a logging operation happens. Keep in mind that
// your hooks can impact performance and that this is not safe to be used concurrently, nor when already logging, so use
// this when launching the application and configuring the logger.
func AddHooks(hks ...Hook) {
	for _, h := range hks {
		hooks = append(hooks, h)
	}
}

// from returns a logger if it is present in the context, otherwise it returns the default logger.
func from(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(key{}).(*zap.Logger); ok {
		return l
	}

	l := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "@timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.NanosDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(os.Stdout),
			&logLevel,
		),
		zap.AddCaller(),
		// since we are encapsulating the logger inside the shorthand functions we need to skip one caller
		// otherwise all callers would be the log function in this file
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel))

	if len(hooks) > 0 {
		l = l.WithOptions(zap.Hooks(hooks...))
	}

	return l
}

// From publicly exposes the underlying *zap.Logger.
func From(ctx context.Context) *zap.Logger {
	return from(ctx)
}

// WithFields will return a context with a new logger with fields.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	return With(ctx, from(ctx).With(fields...))
}

// WithTraceID will add the Trace ID field to logger in the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return WithFields(ctx, zap.String("trace_id", traceID))
}

// With will add logger to the context. Use this if you want to configure a different logger rather than the default one.
func With(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

// Error will log at the error level using the logger associated with the context.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	from(ctx).Error(msg, fields...)
}

// Info will log at the info level using the logger associated with the context.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	from(ctx).Info(msg, fields...)
}

// Debug will log at the info level using the logger associated with the context.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	from(ctx).Debug(msg, fields...)
}

// Warn will log at the info level using the logger associated with the context.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	from(ctx).Warn(msg, fields...)
}

// Fatal will log at the info level using the logger associated with the context.
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	from(ctx).Fatal(msg, fields...)
}
