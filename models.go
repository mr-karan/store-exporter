package main

import (
	"sync"
	"time"

	"github.com/mr-karan/store-exporter/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Hub represents the structure for all app wide functions and structs
type Hub struct {
	logger  *logrus.Logger
	config  config
	version string
}

// cfgApp represents the structure to hold App specific configuration.
type cfgApp struct {
	LogLevel string `koanf:"log_level"`
	Jobs     []Job  `koanf:"jobs"`
}

// cfgServer represents the structure to hold Server specific configuration
type cfgServer struct {
	Name         string        `koanf:"name"`
	Address      string        `koanf:"address"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
	MaxBodySize  int           `koanf:"max_body_size"`
}

// config represents the structure to hold configuration loaded from an external data source.
type config struct {
	App    cfgApp    `koanf:"app"`
	Server cfgServer `koanf:"server"`
}

// Job represents a list of scrape jobs with additional config for the target.
type Job struct {
	Name      string   `koanf:"name"`
	DB        string   `koanf:"db"`
	DSN       string   `koanf:"dsn"`
	QueryFile string   `koanf:"query"`
	Metrics   []Metric `koanf:"metrics"`
}

// Exporter represents the structure to hold Prometheus Descriptors. It implements prometheus.Collector
type Exporter struct {
	sync.Mutex                  // Lock exporter to protect from concurrent scrapes.
	hub        *Hub             // To access logger and other app wide config.
	job        *Job             // Holds the Job metadata.
	manager    store.Manager    // Implements Manager interface which is a set of methods to interact with dataset.
	up         *prometheus.Desc // Represents if a scrape was successful or not.
	version    *prometheus.Desc // Represents verion of the exporter.
}

// Metric represents the structure to hold details about constructing a Prometheus.Metric
type Metric struct {
	Query   string   `koanf:"query"`
	Columns []string `koanf:"columns"`
	Help    string   `koanf:"help"`
	Labels  []string `koanf:"labels"`
}
