package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}
func DebugCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Debug(msg, fs...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}
func InfoCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Info(msg, fs...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}
func WarnCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Warn(msg, fs...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}
func ErrorCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Error(msg, fs...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	GetLogger().Panic(msg, fields...)
}
func PanicCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Panic(msg, fs...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
func FatalCtx(msg string, ctx context.Context, fields ...zap.Field) {
	fs := PickRequestId(ctx, fields)
	GetLogger().Fatal(msg, fs...)
}

func PickRequestId(ctx context.Context, fields []zap.Field) []zap.Field {
	rid := ctx.Value(REQUEST_ID_KEY).(string)
	fields = append(fields, zap.String(REQUEST_ID, rid))
	return fields
}

// With creates a child logger and adds structured context to it.
// Fields added to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// SetLevel alters the logging level.
func SetLevel(l zapcore.Level) {
	_gProps.Load().(*ZapProperties).Level.SetLevel(l)
}

// GetLevel gets the logging level.
func GetLevel() zapcore.Level {
	return _gProps.Load().(*ZapProperties).Level.Level()
}
