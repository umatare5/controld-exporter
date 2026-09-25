package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/umatare5/controld-exporter/internal/config"
	"github.com/umatare5/controld-exporter/internal/controld"
)

func newTestServer() Server {
	return NewServer(&config.Config{
		WebListenAddress: "127.0.0.1",
		WebListenPort:    10034,
		WebTelemetryPath: "/metrics",
		ControlDAPIKey:   "test-key",
	})
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	s := newTestServer()
	if s.Client == nil {
		t.Error("client is nil, want one built from the configured key")
	}
	if s.Config.WebTelemetryPath != "/metrics" {
		t.Errorf("telemetry path = %q, want /metrics", s.Config.WebTelemetryPath)
	}
}

// The landing page names the listen address and the scrape path, and nothing
// an operator supplied beyond them.
func TestHelpPage(t *testing.T) {
	t.Parallel()

	s := newTestServer()
	w := httptest.NewRecorder()

	s.help(w, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "http://127.0.0.1:10034/metrics") {
		t.Errorf("body = %q, want it to carry the telemetry URL", body)
	}
	if !strings.Contains(body, "<h1>Prometheus ControlD Exporter</h1>") {
		t.Error("body lost its heading")
	}
}

// The handler serves whatever the collector gathered, so a reachable API must
// reach the exposition.
func TestMetricsHandlerServesTheCollector(t *testing.T) {
	t.Parallel()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != controld.NetworkEndpoint {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(
			`{"success":true,"body":{"network":[{"iata_code":"NRT","city_name":"Tokyo",` +
				`"country_name":"JP","status":{"api":1,"dns":1,"pxy":-1}}]}}`,
		))
	}))
	t.Cleanup(api.Close)

	s := Server{
		Client: controld.NewClientWithBaseURL(api.URL, "test-key"),
		Config: &config.Config{WebTelemetryPath: "/metrics"},
	}
	w := httptest.NewRecorder()

	s.metricsHandler(w, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `controld_network_health_code{city_name="Tokyo"`) {
		t.Errorf("body = %q, want the network series", w.Body.String())
	}
}

// Every read failing leaves no families, and ContinueOnError answers 200 with an
// empty body rather than failing the scrape.
func TestMetricsHandlerSurvivesAnUnreachableAPI(t *testing.T) {
	t.Parallel()

	api := httptest.NewServer(http.NotFoundHandler())
	api.Close()

	s := Server{
		Client: controld.NewClientWithBaseURL(api.URL, "test-key"),
		Config: &config.Config{WebTelemetryPath: "/metrics"},
	}
	w := httptest.NewRecorder()

	s.metricsHandler(w, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if strings.Contains(w.Body.String(), "controld_") {
		t.Error("series published although the API is unreachable")
	}
}
