package log

import (
	"io"
	"log"
	"os"
	"path"
	"sync/atomic"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	_gLogger atomic.Value
	_gProps  atomic.Value
	_gSugar  atomic.Value
)

func New(cfg *Config) {
	if cfg == nil {
		cfg = &Config{
			Level:  "info",
			Format: "console",
			File: FileLogConfig{
				MaxSize: 300,
			},
		}
	}
	l, p, _ := InitLogger(cfg)
	Reset(l, p)
}

func Reset(logger *zap.Logger, props *ZapProperties) {
	_gLogger.Store(logger)
	_gSugar.Store(logger.Sugar())
	_gProps.Store(props)
}

func GetLogger() *zap.Logger {
	return _gLogger.Load().(*zap.Logger)
}

func GetSurgar() *zap.SugaredLogger {
	return _gSugar.Load().(*zap.SugaredLogger)
}

func Sync() error {
	err := GetLogger().Sync()
	if err != nil {
		return err
	}
	return GetSurgar().Sync()
}

// InitLogger initializes a zap logger.
func InitLogger(cfg *Config, opts ...zap.Option) (*zap.Logger, *ZapProperties, error) {
	var output zapcore.WriteSyncer
	if len(cfg.File.FileName) > 0 {
		output = zapcore.AddSync(getWriter(&cfg.File))
	} else {
		stdOut, _, err := zap.Open([]string{"stdout"}...)
		if err != nil {
			return nil, nil, err
		}
		output = stdOut
	}
	return InitLoggerWithWriteSyncer(cfg, output, opts...)
}

func BuildEncoder(format string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        Datetime.ToString(),
		LevelKey:       LevelKey.ToString(),
		NameKey:        "logger",
		CallerKey:      Caller.ToString(),
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     Message.ToString(),
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	} // zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	var encoder zapcore.Encoder
	switch format {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	}
	return encoder
}

// InitLoggerWithWriteSyncer initializes a zap logger with specified  write syncer.
func InitLoggerWithWriteSyncer(cfg *Config, output zapcore.WriteSyncer, opts ...zap.Option) (*zap.Logger, *ZapProperties, error) {
	// get level
	def := DebugLevel
	lv := &def
	stdLevel := &def

	coreTree := make([]zapcore.Core, 0)

	if cfg != nil {
		lv := cfg.GetLevel()
		stdLevel = cfg.GetStdLevel()
		// build file core
		fileEncoder := BuildEncoder(cfg.Format)
		coreTree = append(coreTree, zapcore.NewCore(fileEncoder, output, lv.zapLevel()))
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	coreTree = append(coreTree, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdLevel.zapLevel()))
	cores := zapcore.NewTee(coreTree...)

	// build log
	lg := zap.New(cores, zap.AddCaller(), zap.AddCallerSkip(cfg.CallSkip))
	r := &ZapProperties{
		Core:   cores,
		Syncer: output,
		Level:  lv,
	}
	//  replace Globals log
	zap.ReplaceGlobals(lg)
	return lg, r, nil
}

// 构建日志文件路径
func initFileLogDir(cfg *FileLogConfig) (string, error) {
	if len(cfg.FileDir) > 0 {
		return cfg.FileDir, nil
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return cfg.FileDir, err
		} else {
			return pwd, err
		}
	}
}

// 构建文件名称
// Build flie log name
func initFileLogName(cfg *FileLogConfig) (string, error) {
	dir, err := initFileLogDir(cfg)
	if err != nil {
		return "", err
	}
	if len(cfg.FileName) > 0 {
		return path.Join(dir, cfg.FileName), nil
	} else {
		return "app", nil
	}
}

// 文件写入器
// initFileLog initializes file based logging options.
func initFileLog(cfg *FileLogConfig) (*lumberjack.Logger, error) {
	// 构建文件名称
	// Build flie log name
	targetPath, err := initFileLogName(cfg)
	if err != nil {
		return nil, err
	}
	// use lumberjack to logrotate
	return &lumberjack.Logger{
		Filename:   targetPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxDays,
		LocalTime:  true,
		Compress:   cfg.Compress,
	}, nil
}

func getWriter(cfg *FileLogConfig) io.Writer {
	filename, _ := initFileLogName(cfg)
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(cfg.MaxDays)), // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24),                      //切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}
