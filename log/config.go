package log

// FileLogConfig serializes file log related config in toml/json.
type FileLogConfig struct {
	FileDir string `toml:"filedir" json:"filedir"`
	// Log filename, leave empty to disable file log.
	Filename string `toml:"filename" json:"filename"`
	// Max size for a single file, in MB.
	MaxSize int `toml:"max-size" json:"max-size"`
	// Max log keep days, default is never deleting.
	MaxDays int `toml:"max-days" json:"max-days"`
	// Maximum number of old log files to retain.
	MaxBackups int  `toml:"max-backups" json:"max-backups"`
	Compress   bool `toml:"compress" json:"compress"`
}

// Config serializes log related config in toml/json.
type Config struct {
	// Log CallSkip
	CallSkip     int  `toml:"callSkip" json:"callSkip"`
	IsProduction bool `toml:"isProduction" json:"isProduction"`
	// Log level.
	Level string `toml:"level" json:"level"`
	// Log format. one of json, text, or console.
	Format string `toml:"format" json:"format"`
	// File log config.
	File FileLogConfig `toml:"file" json:"file"`
}
