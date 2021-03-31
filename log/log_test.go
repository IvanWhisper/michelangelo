package log

import (
	"log"
	"testing"
)

func TestUpdateLevel(t *testing.T) {
	f := FileLogConfig{
		FileDir:    "./logs",
		Filename:   "test",
		MaxSize:    10,
		MaxDays:    1,
		MaxBackups: 10,
		Compress:   false,
	}
	cfg := &Config{
		CallSkip: 1,
		Level:    "info",
		StdLevel: "debug",
		Format:   "console",
		File:     f,
	}
	New(cfg)
	log.Print(GetLevel())
	Debug("debug")
	Info("info")
	Error("error")
}
