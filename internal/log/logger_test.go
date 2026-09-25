package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// The level lives on a package-level logger, so these tests cannot run beside
// each other or beside anything else that logs.
//
//nolint:paralleltest // shared logger state
func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name string
		want logrus.Level
	}{
		{name: "debug", want: logrus.DebugLevel},
		{name: "warn", want: logrus.WarnLevel},
		{name: "error", want: logrus.ErrorLevel},
		{name: "info", want: logrus.InfoLevel},
		{name: "unknown", want: logrus.InfoLevel}, // anything unnamed falls back to info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //nolint:paralleltest // shared logger state
			SetLogLevel(tt.name)
			if got := logger.GetLevel(); got != tt.want {
				t.Errorf("level = %v, want %v", got, tt.want)
			}
		})
	}
}

// The wrappers only forward, so one debug-level pass proves every level reaches
// logrus without panicking.
//
//nolint:paralleltest // shared logger state
func TestWrappersForward(t *testing.T) {
	SetLogLevel("debug")
	t.Cleanup(func() { SetLogLevel("info") })

	Info("info")
	Infof("infof %d", 1)
	Warnf("warnf %d", 2)
	Errorf("errorf %d", 3)
	Debugf("debugf %d", 4)
}
