package cli

import (
	"testing"

	cli "github.com/urfave/cli/v3"

	"github.com/umatare5/controld-exporter/internal/config"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()

	if getVersion() != version {
		t.Errorf("getVersion = %q, want the stamped value", getVersion())
	}
}

// The flag set is the exporter's contract with an operator, so every name and
// default is pinned here rather than in the --help transcript alone.
func TestRegisterFlags(t *testing.T) {
	t.Parallel()

	byName := map[string]cli.Flag{}
	for _, f := range registerFlags() {
		byName[f.Names()[0]] = f
	}

	tests := []struct {
		name string
		want any
	}{
		{config.WebListenAddressFlagName, "0.0.0.0"},
		{config.WebListenPortFlagName, 10034},
		{config.WebTelemetryPathFlagName, "/metrics"},
		{config.ControlDAPIKeyFlagName, ""},
		{config.ControlDBusinessModeFlagName, false},
		{config.LogLevelFlagName, "info"},
	}
	if len(byName) != len(tests) {
		t.Fatalf("flags = %d, want %d", len(byName), len(tests))
	}

	for _, tt := range tests {
		f, ok := byName[tt.name]
		if !ok {
			t.Errorf("%s is not registered", tt.name)
			continue
		}
		var got any
		switch v := f.(type) {
		case *cli.StringFlag:
			got = v.Value
		case *cli.IntFlag:
			got = v.Value
		case *cli.BoolFlag:
			got = v.Value
		}
		if got != tt.want {
			t.Errorf("%s default = %v, want %v", tt.name, got, tt.want)
		}
	}

	// The key is required and reads CTRLD_API_KEY, which is how the container runs.
	key, ok := byName[config.ControlDAPIKeyFlagName].(*cli.StringFlag)
	if !ok {
		t.Fatal("the API key flag is not a string flag")
	}
	if !key.Required {
		t.Error("the API key flag is optional, want it required")
	}
}
