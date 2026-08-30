// Package config is responsible for the execution of the CLI.
package config

import (
	"errors"

	"github.com/jinzhu/configor"
	cli "github.com/urfave/cli/v3"

	"github.com/umatare5/controld-exporter/internal/log"
)

const (
	WebListenAddressFlagName     = "web.listen-address"
	WebListenPortFlagName        = "web.listen-port"
	WebTelemetryPathFlagName     = "web.telemetry-path"
	ControlDAPIKeyFlagName       = "controld.api-key"
	ControlDBusinessModeFlagName = "controld.business-mode"
	LogLevelFlagName             = "log.level"
)

// Config struct holds the configuration for the exporter.
type Config struct {
	WebListenAddress     string
	WebListenPort        int
	WebTelemetryPath     string
	ControlDAPIKey       string
	ControlDBusinessMode bool
	LogLevel             string
}

// NewConfig initializes a Config struct, loads configuration values, and validates the API key.
func NewConfig(cmd *cli.Command) Config {
	config := Config{
		WebListenAddress:     cmd.String(WebListenAddressFlagName),
		WebListenPort:        cmd.Int(WebListenPortFlagName),
		WebTelemetryPath:     cmd.String(WebTelemetryPathFlagName),
		ControlDAPIKey:       cmd.String(ControlDAPIKeyFlagName),
		ControlDBusinessMode: cmd.Bool(ControlDBusinessModeFlagName),
		LogLevel:             cmd.String(LogLevelFlagName),
	}

	err := configor.New(&configor.Config{}).Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	if err := isValidControlDAPIKeyFlag(config.ControlDAPIKey); err != nil {
		log.Fatal(err)
	}

	return config
}

// isValidControlDAPIKeyFlag checks if the ControlD API key is set.
func isValidControlDAPIKeyFlag(apikey string) error {
	if apikey == "" {
		return errors.New("environment variable 'CTRLD_API_KEY' is not set")
	}

	return nil
}
