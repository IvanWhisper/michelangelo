package log

import (
	"context"
	"go.uber.org/zap"
)

// Debug
/**
 * @Description:
 * @param msg
 * @param fields
 */
func Debug(msg string, fields ...zap.Field) {
	DebugCtx(context.TODO(), msg, fields...)
}

// DebugCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param fields
 */
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, DebugLevel, GetLogger().Debug, fields...)
}

// Info
/**
 * @Description:
 * @param msg
 * @param fields
 */
func Info(msg string, fields ...zap.Field) {
	InfoCtx(context.TODO(), msg, fields...)
}

// InfoCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param fields
 */
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, InfoLevel, GetLogger().Info, fields...)
}

// Warn
/**
 * @Description:
 * @param msg
 * @param fields
 */
func Warn(msg string, fields ...zap.Field) {
	WarnCtx(context.TODO(), msg, fields...)
}
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, WarnLevel, GetLogger().Warn, fields...)
}

// Error
/**
 * @Description:
 * @param msg
 * @param fields
 */
func Error(msg string, fields ...zap.Field) {
	ErrorCtx(context.TODO(), msg, fields...)
}

// ErrorCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param fields
 */
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, ErrorLevel, GetLogger().Error, fields...)
}

// Panic
/**
 * @Description:
 * @param msg
 * @param fields
 */
func Panic(msg string, fields ...zap.Field) {
	PanicCtx(context.TODO(), msg, fields...)
}

// PanicCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param fields
 */
func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, CriticalLevel, GetLogger().Panic, fields...)
}

// Fatal
/**
 * @Description:
 * @param msg
 * @param fields
 */func Fatal(msg string, fields ...zap.Field) {
	FatalCtx(nil, msg, fields...)
}

// FatalCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param fields
 */
func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logOutputCtx(ctx, msg, CriticalLevel, GetLogger().Panic, fields...)
}

// logOutputCtx
/**
 * @Description:
 * @param ctx
 * @param msg
 * @param level
 * @param logFunc
 * @param fields
 */
func logOutputCtx(ctx context.Context, msg string, level Level, logFunc func(string, ...zap.Field), fields ...zap.Field) {
	if GetLevel().Unabled(level) {
		return
	}
	fields = append(fields, Ctx2Fields(ctx)...)
	logFunc(msg, fields...)
}

// Ctx2Fields
/**
 * @Description:
 * @param ctx
 * @return []zap.Field
 */
func Ctx2Fields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	if ctx != nil && ctx != context.TODO() {
		if v := ctx.Value(RequestIdKey); v != nil {
			rid := v.(string)
			fields = append(fields, zap.String(SessionId.ToString(), rid))
			fields = append(fields, zap.String(TraceId.ToString(), rid))
		}
		if v := ctx.Value(BusinessKeyword); v != nil {
			kw := v.(string)
			fields = append(fields, zap.String(BusinessKeyword.ToString(), kw))
		}
		if v := ctx.Value(BusinessOperation); v != nil {
			kw := v.(string)
			fields = append(fields, zap.String(BusinessOperation.ToString(), kw))
		}
		if v := ctx.Value(BusinessTitle); v != nil {
			kw := v.(string)
			fields = append(fields, zap.String(BusinessTitle.ToString(), kw))
		}
	}
	return fields
}

// With
/**
 * @Description: creates a child logger and adds structured context to it. Fields added to the child don't affect the parent, and vice versa.
 * @param fields
 * @return *zap.Logger
 */
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// SetLevel alters the logging level.
/**
 * @Description:
 * @param levelStr
 */
func SetLevel(levelStr string) {
	l := new(Level)
	l.Unpack(levelStr)
	_gProps.Load().(*ZapProperties).Level = l
}

// GetLevel
/**
 * @Description:
 * @return *Level
 */
func GetLevel() *Level {
	return _gProps.Load().(*ZapProperties).Level
}
