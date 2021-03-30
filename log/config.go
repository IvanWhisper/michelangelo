package log

// FileLogConfig serializes file log related config in toml/json.
type FileLogConfig struct {
	FileDir    string `toml:"filedir" json:"filedir"`
	Filename   string `toml:"filename" json:"filename"`       // Log filename, leave empty to disable file log.
	MaxSize    int    `toml:"max-size" json:"max-size"`       // Max size for a single file, in MB.
	MaxDays    int    `toml:"max-days" json:"max-days"`       // Max log keep days, default is never deleting.
	MaxBackups int    `toml:"max-backups" json:"max-backups"` // Maximum number of old log files to retain.
	Compress   bool   `toml:"compress" json:"compress"`       // Compress
}

// Config serializes log related config in toml/json.
type Config struct {
	CallSkip int           `toml:"callSkip" json:"callSkip"`   // Log CallSkip
	Level    string        `toml:"level" json:"level"`         // Log level.
	StdLevel string        `toml:"std-level" json:"std-level"` // console level
	Format   string        `toml:"format" json:"format"`       // Log format. one of json, text, or console.
	File     FileLogConfig `toml:"file" json:"file"`           // File log config.
}

func (c *Config) GetLevel() *Level {
	l := new(Level)
	l.Unpack(c.Level)
	return l
}

func (c *Config) GetStdLevel() *Level {
	l := new(Level)
	l.Unpack(c.StdLevel)
	return l
}
