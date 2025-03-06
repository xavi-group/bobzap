package bobzap

import (
	"errors"
	"fmt"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var (
	loggerInitLock sync.RWMutex
	zapConfig      *zap.Config
	configLock     sync.RWMutex
	defaultConfig  *Config
)

// InitializeGlobalLogger defines global zap and open-telemetry zap loggers configured via the first given Config or
// the default config if it has been initialized by bobzap.NewConfig().
func InitializeGlobalLogger(config ...*Config) error {
	var c *Config

	if len(config) > 0 {
		c = config[0]
	} else {
		configLock.RLock()
		defer configLock.RUnlock()

		c = defaultConfig
	}

	if c == nil {
		return errors.New("no logger configuration provided or found")
	}

	zapConfig, err := getZapConfig(c)
	if err != nil {
		return fmt.Errorf("problem creating zap configuration: %w", err)
	}

	zapLogger := zap.Must(zapConfig.Build())
	zapLogger = zapLogger.With(zap.String("id", c.AppID))

	defer func() {
		_ = zapLogger.Sync()
	}()

	otelLogger := otelzap.New(zapLogger, otelzap.WithMinLevel(zapcore.InfoLevel))

	defer func() {
		_ = otelLogger.Sync()
	}()

	zap.ReplaceGlobals(zapLogger)
	otelzap.ReplaceGlobals(otelLogger)

	return nil
}

// NewLogger creates an open-telemetry zap logger with the given name, and attaches info+ logs to traces.
func NewLogger(name string) *otelzap.Logger {
	return otelzap.New(zap.L().Named(name), otelzap.WithMinLevel(zapcore.InfoLevel))
}

// NewObserverLogger creates an open-telemetry zap logger with the given name, and provides a struct that can be
// utilized to observe log messages created with the provided logger.
func NewObserverLogger(name string) (*otelzap.Logger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	return otelzap.New(observedLogger, otelzap.WithMinLevel(zapcore.InfoLevel)),
		observedLogs
}

// SetGlobalLogLevel updates the log level for log messages throughout the application.
func SetGlobalLogLevel(level string) error {
	zapLogLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("problem parsing log level: %w", err)
	}

	if zapConfig == nil {
		return fmt.Errorf("global logger not initialized")
	}

	zapConfig.Level.SetLevel(zapLogLevel)

	return nil
}

func getZapConfig(c *Config) (*zap.Config, error) {
	if zapConfig != nil {
		return zapConfig, nil
	}

	var newZapConfig zap.Config

	// Parse out the logging configuration
	switch c.LogConfig {
	case "production":
		newZapConfig = zap.NewProductionConfig()
	case "development":
		newZapConfig = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("unsupported log config value: '%s'", c.LogConfig)
	}

	// Parse out the logging level
	if c.LogLevel != "" {
		var err error

		newZapConfig.Level, err = zap.ParseAtomicLevel(c.LogLevel)
		if err != nil {
			return nil, fmt.Errorf("unsupported log level value: '%s'", c.LogLevel)
		}
	}

	// Parse out the logging format / encoding
	if c.LogFormat != "" && c.LogFormat != "console" && c.LogFormat != "json" {
		return nil, fmt.Errorf("unsupported log format value: '%s'", c.LogFormat)
	}

	newZapConfig.Encoding = c.LogFormat
	newZapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Handle color for console encoding
	if newZapConfig.Encoding == "console" && c.LogColor {
		newZapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig = &newZapConfig

	return &newZapConfig, nil
}
