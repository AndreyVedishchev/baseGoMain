package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	CommandsTotal   *prometheus.CounterVec
	CommandDuration *prometheus.HistogramVec
	StorageSize     prometheus.Gauge
}

func NewMetrics() *Metrics {
	m := &Metrics{
		CommandsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_commands_total",
				Help: "Total number of executed CLI commands, labeled by command name and status",
			},
			[]string{"command", "status"},
		),
		CommandDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "app_command_duration_seconds",
				Help:    "Duration of CLI command execution in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"command"},
		),
		StorageSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "app_storage_size",
				Help: "Current number of items in storage",
			},
		),
	}
	prometheus.MustRegister(m.CommandsTotal, m.CommandDuration, m.StorageSize)
	return m
}
