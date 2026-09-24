package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_commands_total",
			Help: "Total number of executed CLI commands, labeled by command name and status",
		},
		[]string{"command", "status"},
	)

	CommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_command_duration_seconds",
			Help:    "Duration of CLI command execution in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)

	StorageSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_storage_size",
			Help: "Current number of items in storage",
		},
	)
)
