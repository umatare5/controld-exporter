// Package cli handles the execution of the CLI application.
package cli

import (
	"context"
	"os"

	"github.com/umatare5/controld-exporter/internal/config"
	"github.com/umatare5/controld-exporter/internal/log"
	"github.com/umatare5/controld-exporter/internal/server"
	cli "github.com/urfave/cli/v3"
)

// Run initializes and starts the CLI application.
func Run() {
	cmd := &cli.Command{
		Name:      "controld-exporter",
		Usage:     "A Prometheus exporter for metrics from the Control D",
		UsageText: "controld-exporter [options...]",
		Version:   getVersion(),
		Flags:     registerFlags(),
		Action: func(ctx context.Context, cli *cli.Command) error {
			config := config.NewConfig(cli)
			exporter, _ := server.NewServer(&config)

			exporter.Start()

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// registerFlags defines and returns the global CLI flags.
func registerFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerWebListenAddressFlag()...)
	flags = append(flags, registerWebListenPortFlag()...)
	flags = append(flags, registerWebTelemetryPathFlag()...)
	flags = append(flags, registerAPIKeyFlag()...)
	flags = append(flags, registerBusinessModeFlag()...)
	flags = append(flags, registerLogLevelFlag()...)
	return flags
}

// registerWebListenAddressFlag defines the flag for the server's listen address.
func registerWebListenAddressFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  config.WebListenAddressFlagName,
			Usage: "Address to bind the HTTP server to.",
			Value: "0.0.0.0",
		},
	}
}

// registerWebListenPortFlag defines the flag for the server's listen port.
func registerWebListenPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  config.WebListenPortFlagName,
			Usage: "Port number to bind the HTTP server to.",
			Value: 10034,
		},
	}
}

// registerWebTelemetryPathFlag defines the flag for the metrics telemetry path.
func registerWebTelemetryPathFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.WebTelemetryPathFlagName,
			Usage:   "Set the path to expose metrics",
			Aliases: []string{"p"},
			Value:   "/metrics",
		},
	}
}

// registerAPIKeyFlag defines the flag for the Control D API key.
func registerAPIKeyFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     config.ControlDAPIKeyFlagName,
			Usage:    "API key for authenticating with the Control D API.",
			Aliases:  []string{"k"},
			Sources:  cli.EnvVars("CTRLD_API_KEY"),
			Required: true,
		},
	}
}

// registerBusinessModeFlag defines the flag for enabling business mode.
func registerBusinessModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  config.ControlDBusinessModeFlagName,
			Usage: "Enable the metrics collection available in the business subscription.",
			Value: false,
		},
	}
}

// registerLogLevelFlag defines the flag for setting the logging level.
func registerLogLevelFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  config.LogLevelFlagName,
			Usage: "Set the logging level. One of: [debug, info, warn, error]",
			Value: "info",
		},
	}
}
