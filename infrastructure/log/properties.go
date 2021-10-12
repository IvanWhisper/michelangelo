package log

import (
	"go.uber.org/zap/zapcore"
)

type ZapProperties struct {
	Core   zapcore.Core
	Syncer zapcore.WriteSyncer
	Level  *Level
}
