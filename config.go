package bobzap

import "github.com/xavi-group/bconf"

const (
	// LogFieldSetKey defines the field-set key for logging configuration fields.
	LogFieldSetKey = "log"
	// LogColorKey ...
	LogColorKey = "color"
	// LogConfigKey ...
	LogConfigKey = "config"
	// LogFormatKey ...
	LogFormatKey = "format"
	// LogLevelKey ...
	LogLevelKey = "level"
)

func NewConfig() *Config {
	configLock.Lock()
	defer configLock.Unlock()

	defaultConfig = &Config{}

	return defaultConfig
}

// Config defines the expected values for configuring an application logger. It is recommended to initialize a
// Config with either bobzap.NewConfig() in order to set the default config when initializing the global logger.
type Config struct {
	bconf.ConfigStruct
	AppID     string `bconf:"app.id"`
	LogColor  bool   `bconf:"log.color"`
	LogConfig string `bconf:"log.config"`
	LogFormat string `bconf:"log.format"`
	LogLevel  string `bconf:"log.level"`
}

// LoggerFieldSet defines the field-set for an application logger.
func LoggerFieldSet() *bconf.FieldSet {
	return bconf.FSB(LogFieldSetKey).Fields(
		bconf.FB(LogColorKey, bconf.Bool).Default(true).C(),
		bconf.FB(LogConfigKey, bconf.String).Default("production").Enumeration("production", "development").C(),
		bconf.FB(LogFormatKey, bconf.String).Default("json").Enumeration("console", "json").C(),
		bconf.FB(LogLevelKey, bconf.String).Default("info").
			Enumeration("debug", "info", "warn", "error", "dpanic", "panic", "fatal").C(),
	).C()
}

// LoggerFieldSets defines the field-sets for an applicaiton logger.
func LoggerFieldSets() bconf.FieldSets {
	return bconf.FieldSets{
		LoggerFieldSet(),
	}
}
