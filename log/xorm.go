package log

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	xlog "xorm.io/xorm/log"
)

type OrmLogger struct {
	level xlog.LogLevel
}

func (l *OrmLogger) Debug(v ...interface{}) {
	zap.L().Debug(Interfaces2String(v))
}
func (l *OrmLogger) Debugf(format string, v ...interface{}) {
	zap.L().Debug(fmt.Sprintf(format, v))
}
func (l *OrmLogger) Error(v ...interface{}) {
	zap.L().Error(Interfaces2String(v))
}
func (l *OrmLogger) Errorf(format string, v ...interface{}) {
	zap.L().Error(fmt.Sprintf(format, v))
}
func (l *OrmLogger) Info(v ...interface{}) {
	zap.L().Info(Interfaces2String(v))
}
func (l *OrmLogger) Infof(format string, v ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, v))
}
func (l *OrmLogger) Warn(v ...interface{}) {
	zap.L().Warn(Interfaces2String(v))
}
func (l *OrmLogger) Warnf(format string, v ...interface{}) {
	zap.L().Warn(fmt.Sprintf(format, v))
}

func (l *OrmLogger) Level() xlog.LogLevel {
	return l.level
}
func (l *OrmLogger) SetLevel(level xlog.LogLevel) {
	l.level = level
}

func (l *OrmLogger) ShowSQL(show ...bool) {

}
func (l *OrmLogger) IsShowSQL() bool {
	return true
}

type OrmCtxLogger struct {
	level xlog.LogLevel
}

func (l *OrmCtxLogger) BeforeSQL(context xlog.LogContext) {}
func (l *OrmCtxLogger) AfterSQL(context xlog.LogContext)  {}

func (l *OrmCtxLogger) Debugf(format string, v ...interface{}) {
	zap.L().Debug(fmt.Sprintf(format, v))
}
func (l *OrmCtxLogger) Errorf(format string, v ...interface{}) {
	zap.L().Error(fmt.Sprintf(format, v))
}
func (l *OrmCtxLogger) Infof(format string, v ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, v))
}
func (l *OrmCtxLogger) Warnf(format string, v ...interface{}) {
	zap.L().Warn(fmt.Sprintf(format, v))
}

func (l *OrmCtxLogger) Level() xlog.LogLevel {
	return l.level
}

func (l *OrmCtxLogger) SetLevel(level xlog.LogLevel) {
	l.level = level
}

func (l *OrmCtxLogger) ShowSQL(show ...bool) {}

func (l *OrmCtxLogger) IsShowSQL() bool {
	return true
}

func Interfaces2String(v ...interface{}) string {
	var sb strings.Builder
	for _, i := range v {
		sb.WriteString(i.(string))
	}
	return sb.String()
}
