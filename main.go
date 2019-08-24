package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	// injected during build
	buildVersion = "unknown"
	buildDate    = "unknown"
)

func initLogger(config cfgApp) *logrus.Logger {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// Set logger level
	switch level := config.LogLevel; level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("verbose logging enabled")
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

func main() {
	var (
		config = initConfig()
		logger = initLogger(config.App)
	)
	// Initialize hub which contains initializations of app level objects
	hub := &Hub{
		config:  config,
		logger:  logger,
		version: buildVersion,
	}
	hub.logger.Infof("booting store-exporter-version:%v", buildVersion)
	// Initialize prometheus registry to register metrics and collectors.
	r := prometheus.NewRegistry()
	// Fetch all jobs listed in config and register with the registry.
	for _, job := range hub.config.App.Jobs {
		// This is to avoid all copies of `exporter` getting updated by the last `job` memory address
		// you instantiate with, since we pass `job` as a pointer to the struct.
		j := job
		// Initialize the exporter. Exporter is a collection of metrics to be exported.
		exporter, err := hub.NewExporter(hub.config.App.Namespace, &j, hub.config.App.QueryFile)
		if err != nil {
			hub.logger.Panicf("exporter initialization failed for %s : %s", job.Name, err)
		}
		// Register the exporters with our custom registry. Panics in case of failure.
		r.MustRegister(exporter)
		hub.logger.Debugf("registration of metrics for job %s success", job.Name)
	}
	// Default index handler.
	handleIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to store-exporter. Visit /metrics to scrape prometheus metrics."))
	})
	// Initialize router and define all endpoints.
	router := http.NewServeMux()
	router.Handle("/", handleIndex)
	router.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	// Initialize server.
	server := &http.Server{
		Addr:         hub.config.Server.Address,
		Handler:      router,
		ReadTimeout:  hub.config.Server.ReadTimeout * time.Millisecond,
		WriteTimeout: hub.config.Server.WriteTimeout * time.Millisecond,
	}
	// Start the server. Blocks the main thread.
	hub.logger.Infof("starting server listening on %v", hub.config.Server.Address)
	if err := server.ListenAndServe(); err != nil {
		hub.logger.Fatalf("error starting server: %v", err)
	}
}
