package bobzap

import "github.com/xavi-group/bconf"

const (
	// LogFieldSetKey defines the field-set key for logging configuration fields.
	LogFieldSetKey = "log"
	// LogColorKey defines the field key for log color.
	LogColorKey = "color"
	// LogConfigKey defines the field key for log config.
	LogConfigKey = "config"
	// LogFormatKey defines the field key for log format.
	LogFormatKey = "format"
	// LogLevelKey defines the field key for log level.
	LogLevelKey = "level"
)

// NewConfig provides an initialized Config struct, and sets the returned config struct as the default config used when
// calling InitializeGlobalLogger(config ...*Config) with no args.
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

// FieldSets defines the field-sets for an otelzap application logger.
func FieldSets() bconf.FieldSets {
	return bconf.FieldSets{
		LoggerFieldSet(),
	}
}

// LoggerFieldSet defines the field-set for an application logger.
func LoggerFieldSet() *bconf.FieldSet {
	return bconf.FSB(LogFieldSetKey).Fields(
		bconf.FB(LogColorKey, bconf.Bool).Default(true).
			Description("Log color defines whether console formatted logs are rendered in color.").C(),
		bconf.FB(LogConfigKey, bconf.String).Default("production").Enumeration("production", "development").
			Description(
				"Log config defines whether the Zap will be configured with development or production defaults. ",
				"Note: `development` defaults to debug log level and console format, `production` defaults to info ",
				"log level and json format.",
			).C(),
		bconf.FB(LogFormatKey, bconf.String).Enumeration("console", "json").
			Description(
				"Log format defines the format logs will be emitted in (overrides log config defaults).",
			).C(),
		bconf.FB(LogLevelKey, bconf.String).Enumeration("debug", "info", "warn", "error", "dpanic", "panic", "fatal").
			Description(
				"Log level defines the level at which logs will be emitted (overrides log config defaults).",
			).C(),
	).C()
}
