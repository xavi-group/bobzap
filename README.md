# `bobzap`: base observability / zap logger

`bobzap` is a support package that provides bconf configuration for open-telemetry zap logging, which is a great way to
improve application observability via log output.

This package additionally provides helper functions for initializing a global logging configuration, and creating new
loggers.

```sh
go get github.com/xavi-group/bobzap
```

## Configuration

```
Optional Configuration:
        app.description string
                Default value: 'Example application showcasing bobzap logging'
                Environment key: 'EXAMPLE_APP_DESCRIPTION'
                Flag argument: '--app_description'
        app.id string
                Default value: <generated-at-run-time>
                Environment key: 'EXAMPLE_APP_ID'
                Flag argument: '--app_id'
        app.name string
                Default value: 'bobzapexample'
                Environment key: 'EXAMPLE_APP_NAME'
                Flag argument: '--app_name'
        app.version string
                Default value: '1.0.0'
                Environment key: 'EXAMPLE_APP_VERSION'
                Flag argument: '--app_version'
        log.color bool
                Default value: 'true'
                Environment key: 'EXAMPLE_LOG_COLOR'
                Flag argument: '--log_color'
        log.config string
                Accepted values: ['production', 'development']
                Default value: 'production'
                Environment key: 'EXAMPLE_LOG_CONFIG'
                Flag argument: '--log_config'
        log.format string
                Accepted values: ['console', 'json']
                Default value: 'json'
                Environment key: 'EXAMPLE_LOG_FORMAT'
                Flag argument: '--log_format'
        log.level string
                Accepted values: ['debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal']
                Default value: 'info'
                Environment key: 'EXAMPLE_LOG_LEVEL'
                Flag argument: '--log_level'
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/segmentio/ksuid"
	"github.com/xavi-group/bconf"
	"github.com/xavi-group/bobzap"
	"go.uber.org/zap"
)

func main() {
	config := bconf.NewAppConfig(
		"bobzapexample",
		"Example application showcasing bobzap logging",
		bconf.WithAppIDFunc(func() string { return ksuid.New().String() }),
		bconf.WithAppVersion("1.0.0"),
		bconf.WithEnvironmentLoader("example"),
		bconf.WithFlagLoader(),
	)

	config.AddFieldSetGroup("bobzap", bobzap.FieldSets())

	config.AttachConfigStructs(
		bobzap.NewConfig(),
	)

	// Load when called without any options will also handle the help flag (--help or -h)
	if errs := config.Load(); len(errs) > 0 {
		fmt.Printf("problem(s) loading application configuration: %v\n", errs)
		os.Exit(1)
	}

	// -- Initialize application observability --
	if err := bobzap.InitializeGlobalLogger(); err != nil {
		fmt.Printf("problem initializing application logger: %s\n", err)
		os.Exit(1)
	}

	log := bobzap.NewLogger("main")

	log.Info(
		fmt.Sprintf("%s initialized successfully", config.AppName()),
		zap.Any("app_config", config.ConfigMap()),
		zap.Strings("warnings", config.Warnings()),
	)
}
```

## Support

For more information on Zap, check out and support the project at
[github.com/uber-go/zap](https://github.com/uber-go/zap)

For more information on Otelzap, check out and support the project at
[github.com/uptrace/opentelemetry-go-extra](https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap)

For more information on Open Telemetry, check out and support the project at
[opentelemetry.io](https://opentelemetry.io/)

For more information on bconf, check out and support the project at
[github.com/xavi-group/bconf](https://github.com/xavi-group/bconf)
