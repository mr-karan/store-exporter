package main

import (
	"context"
	"sync"

	"github.com/mr-karan/store-exporter/store"
	"github.com/prometheus/client_golang/prometheus"
)

// NewExporter returns an initialized `Exporter`.
func (hub *Hub) NewExporter(namespace string, job *Job) (*Exporter, error) {
	manager, err := store.NewManager(job.DB, job.DSN, &store.DBConnOpts{
		QueryFilePath: job.QueryFile,
	})
	if err != nil {
		hub.logger.Errorf("Error initializing database manager: %s", err)
		return nil, err
	}
	return &Exporter{
		Mutex:   sync.Mutex{},
		manager: manager,
		job:     job,
		hub:     hub,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, job.Name, "up"),
			"Could the data source be reached.",
			nil,
			nil,
		),
		version: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "version"),
			"Version of store-exporter",
			[]string{"build"},
			nil,
		),
	}, nil
}

// sendSafeMetric is a concurrent safe method to send metrics to a channel. Since we are collecting metrics from AWS API, there might be possibility where
// a timeout occurs from Prometheus' collection context and the channel is closed but Goroutines running in background can still
// send metrics to this closed channel which would result in panic and crash. To solve that we use context and check if the channel is not closed
// and only send the metrics in that case. Else it logs the error and returns in a safe way.
func (hub *Hub) sendSafeMetric(ctx context.Context, ch chan<- prometheus.Metric, metric prometheus.Metric) error {
	// Check if collection context is finished
	select {
	case <-ctx.Done():
		// don't send metrics, instead return in a "safe" way
		hub.logger.Errorf("Attempted to send metrics to a closed channel after collection context had finished: %s", metric)
		return ctx.Err()
	default: // continue
	}
	// Send metrics if collection context is still open
	ch <- metric
	return nil
}

// Describe describes all the metrics ever exported by the exporter. It implements `prometheus.Collector`.
func (p *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.version
	ch <- p.up
}

// Collect is called by the Prometheus registry when collecting
// metrics. This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. It implements `prometheus.Collector`
func (p *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Initialize context to keep track of the collection.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Lock the exporter for one iteration of collection as `Collect` can be called concurrently.
	p.Lock()
	defer p.Unlock()
	for _, m := range p.job.Metrics {
		p.collectMetrics(ctx, ch, m)
	}
	// Send default metrics data.
	p.hub.sendSafeMetric(ctx, ch, prometheus.MustNewConstMetric(p.version, prometheus.GaugeValue, 1, p.hub.version))
}

// collectMetrics fetches data from external stores and sends as Prometheus metrics
func (p *Exporter) collectMetrics(ctx context.Context, ch chan<- prometheus.Metric, metric Metric) {
	data, err := p.manager.FetchResults(metric.Query)
	if err != nil {
		p.hub.logger.Errorf("Error while fetching result from DB: %v", err)
		p.hub.sendSafeMetric(ctx, ch, prometheus.MustNewConstMetric(p.up, prometheus.GaugeValue, 0))
		return
	}
	value, labelValues, err := constructMetricData(data, metric.Value, metric.Labels)
	if err != nil {
		p.hub.logger.Errorf("Error while converting results to metrics: %v", err)
		p.hub.sendSafeMetric(ctx, ch, prometheus.MustNewConstMetric(p.up, prometheus.GaugeValue, 0))
		return
	}
	// Create metrics on the fly
	metricDesc := createMetricDesc(p.job.Namespace, metric.Name, p.job.Name, metric.Help, metric.Labels)
	p.hub.sendSafeMetric(ctx, ch, prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value, labelValues...))

}

// createMetricDesc returns an intialized prometheus.Desc instance
func createMetricDesc(namespace string, metricName string, jobName string, helpText string, additionalLabels []string) *prometheus.Desc {
	// Default labels for any metric constructed with this function.
	var labels []string
	// Iterate through a slice of additional labels to be exported.
	for _, k := range additionalLabels {
		// Replace all tags with underscores if present to make it a valid Prometheus label name.
		labels = append(labels, replaceWithUnderscores(k))
	}
	return prometheus.NewDesc(
		prometheus.BuildFQName(replaceWithUnderscores(namespace), "", replaceWithUnderscores(metricName)),
		helpText,
		labels, prometheus.Labels{"job": replaceWithUnderscores(jobName)},
	)
}
