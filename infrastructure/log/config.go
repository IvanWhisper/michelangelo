package log

// FileLogConfig
/**
 * @Description: serializes file log related config in toml/json.
 */
type FileLogConfig struct {
	FileDir    string `toml:"fileDir" json:"fileDir"`
	FileName   string `toml:"fileName" json:"fileName"`     // Log filename, leave empty to disable file log.
	MaxSize    int    `toml:"maxSize" json:"maxSize"`       // Max size for a single file, in MB.
	MaxDays    int    `toml:"maxDays" json:"maxDays"`       // Max log keep days, default is never deleting.
	MaxBackups int    `toml:"maxBackups" json:"maxBackups"` // Maximum number of old log files to retain.
	Compress   bool   `toml:"compress" json:"compress"`     // Compress
}

// Config
/**
 * @Description: serializes log related config in toml/json.
 */
type Config struct {
	CallSkip int           `toml:"callSkip" json:"callSkip"` // Log CallSkip
	Level    string        `toml:"level" json:"level"`       // Log level.
	StdLevel string        `toml:"stdLevel" json:"stdLevel"` // console level
	Format   string        `toml:"format" json:"format"`     // Log format. one of json, text, or console.
	File     FileLogConfig `toml:"file" json:"file"`         // File log config.
}

// GetLevel
/**
 * @Description:
 * @receiver c
 * @return *Level
 */
func (c *Config) GetLevel() *Level {
	l := new(Level)
	l.Unpack(c.Level)
	return l
}

// GetStdLevel
/**
 * @Description:
 * @receiver c
 * @return *Level
 */
func (c *Config) GetStdLevel() *Level {
	l := new(Level)
	l.Unpack(c.StdLevel)
	return l
}
