package config

import (
	"context"
	"testing"

	cli "github.com/urfave/cli/v3"
)

func TestNewConfigReadsEveryFlag(t *testing.T) {
	t.Parallel()

	var got Config
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: WebListenAddressFlagName},
			&cli.IntFlag{Name: WebListenPortFlagName},
			&cli.StringFlag{Name: WebTelemetryPathFlagName},
			&cli.StringFlag{Name: ControlDAPIKeyFlagName},
			&cli.BoolFlag{Name: ControlDBusinessModeFlagName},
			&cli.StringFlag{Name: LogLevelFlagName},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			got = NewConfig(c)
			return nil
		},
	}

	args := []string{
		"controld-exporter",
		"--" + WebListenAddressFlagName, "127.0.0.1",
		"--" + WebListenPortFlagName, "19034",
		"--" + WebTelemetryPathFlagName, "/telemetry",
		"--" + ControlDAPIKeyFlagName, "api-key",
		"--" + ControlDBusinessModeFlagName,
		"--" + LogLevelFlagName, "debug",
	}
	if err := cmd.Run(context.Background(), args); err != nil {
		t.Fatalf("run: %v", err)
	}

	want := Config{
		WebListenAddress:     "127.0.0.1",
		WebListenPort:        19034,
		WebTelemetryPath:     "/telemetry",
		ControlDAPIKey:       "api-key",
		ControlDBusinessMode: true,
		LogLevel:             "debug",
	}
	if got != want {
		t.Errorf("config = %+v, want %+v", got, want)
	}
}

// An empty key is fatal in NewConfig, so the check is exercised on its own.
func TestIsValidControlDAPIKeyFlag(t *testing.T) {
	t.Parallel()

	if err := isValidControlDAPIKeyFlag("api-key"); err != nil {
		t.Errorf("a set key returned %v, want nil", err)
	}
	if err := isValidControlDAPIKeyFlag(""); err == nil {
		t.Error("an empty key returned nil, want the CTRLD_API_KEY error")
	}
}
