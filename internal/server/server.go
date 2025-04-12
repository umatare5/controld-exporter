// Package server provides the HTTP server implementation for the exporter.
package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/umatare5/controld-exporter/internal/collector"
	"github.com/umatare5/controld-exporter/internal/config"
	"github.com/umatare5/controld-exporter/internal/controld"
	"github.com/umatare5/controld-exporter/internal/log"
)

// Server represents the HTTP server for the exporter.
type Server struct {
	Client *controld.Client // ControlD API client
	Config *config.Config   // Configuration for the server
}

// NewServer initializes and returns a new Server instance.
func NewServer(config *config.Config) (Server, error) {
	return Server{
		Client: controld.NewClient(config.ControlDAPIKey),
		Config: config,
	}, nil
}

// Start configures and launches the HTTP server to serve metrics and help pages.
func (s *Server) Start() {
	log.SetLogLevel(s.Config.LogLevel)
	reg := prometheus.NewRegistry()

	// Register standard process and Go metrics.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	// Register HTTP handlers.
	http.HandleFunc("/", s.help)
	http.HandleFunc(s.Config.WebTelemetryPath, func(w http.ResponseWriter, r *http.Request) {
		s.metricsHandler(w, r)
	})

	// Print server start message.
	log.Infof(
		"Starting the %s exporter on port %d.",
		map[bool]string{true: "business mode", false: "personal mode"}[s.Config.ControlDBusinessMode],
		s.Config.WebListenPort,
	)

	srv := &http.Server{
		Addr:         s.Config.WebListenAddress + ":" + strconv.Itoa(s.Config.WebListenPort),
		Handler:      nil,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// metricsHandler registers Prometheus metrics and serves them via HTTP.
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()

	// Register the ControlD collector and metrics.
	registry.MustRegister(
		collector.NewCollector(s.Client, s.Config.ControlDBusinessMode),
	)

	// Serve metrics using Prometheus client library.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})
	h.ServeHTTP(w, r)
}

// help generates and serves an HTML help page for the root URL.
func (s *Server) help(w http.ResponseWriter, _ *http.Request) {
	listenAddrAndPort := s.Config.WebListenAddress + ":" + strconv.Itoa(s.Config.WebListenPort)

	var builder strings.Builder
	builder.WriteString("<h1>Prometheus ControlD Exporter</h1>")
	builder.WriteString("<p>To fetch metrics from ControlD, access the telemetry path:</p>")
	builder.WriteString(fmt.Sprintf("http://%s%s", listenAddrAndPort, s.Config.WebTelemetryPath))
	builder.WriteString("<p><b>Example:</b></p>")
	builder.WriteString("<ul>")
	builder.WriteString(fmt.Sprintf("<li><a href=\"http://%s%s\">http://%s%s</a></li>", listenAddrAndPort, s.Config.WebTelemetryPath, listenAddrAndPort, s.Config.WebTelemetryPath))
	builder.WriteString("</ul>")

	if _, err := w.Write([]byte(builder.String())); err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}
