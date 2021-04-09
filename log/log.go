package log

import (
	"context"
	"go.uber.org/zap"
)

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(DebugLevel) {
		return
	}
	GetLogger().Debug(msg, fields...)
}
func DebugCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(DebugLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Debug(msg, fs...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(InfoLevel) {
		return
	}
	GetLogger().Info(msg, fields...)
}
func InfoCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(InfoLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Info(msg, fs...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(WarnLevel) {
		return
	}
	GetLogger().Warn(msg, fields...)
}
func WarnCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(WarnLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Warn(msg, fs...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(ErrorLevel) {
		return
	}
	GetLogger().Error(msg, fields...)
}
func ErrorCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(ErrorLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Error(msg, fs...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(CriticalLevel) {
		return
	}
	GetLogger().Panic(msg, fields...)
}
func PanicCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(CriticalLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Panic(msg, fs...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	if GetLevel().Unabled(CriticalLevel) {
		return
	}
	GetLogger().Fatal(msg, fields...)
}

func FatalCtx(msg string, ctx context.Context, fields ...zap.Field) {
	if GetLevel().Unabled(CriticalLevel) {
		return
	}
	fs := PickRequestId(ctx, fields)
	GetLogger().Fatal(msg, fs...)
}

func PickRequestId(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx != nil {
		if v := ctx.Value(REQUEST_ID_KEY); v != nil {
			rid := v.(string)
			fields = append(fields, zap.String(REQUEST_ID, rid))
		}
	}
	return fields
}

// With creates a child logger and adds structured context to it.
// Fields added to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// SetLevel alters the logging level.
func SetLevel(levelStr string) {
	l := new(Level)
	l.Unpack(levelStr)
	_gProps.Load().(*ZapProperties).Level = l
}

// GetLevel gets the logging level.
func GetLevel() *Level {
	return _gProps.Load().(*ZapProperties).Level
}
