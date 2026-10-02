package db

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dbOperationsTotal   *prometheus.CounterVec
	dbOperationDuration *prometheus.HistogramVec
	metricsOnce sync.Once
)

// initMetrics создаёт и регистрирует метрики БД. Безопасно вызывать повторно.
func initMetrics() {
	metricsOnce.Do(func() {
		dbOperationsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_operations_total",
				Help: "Number of database operations by type and result",
			},
			[]string{"operation", "status"},
		)

		dbOperationDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_operation_duration_seconds",
				Help:    "Duration of database operations",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		)

		prometheus.MustRegister(dbOperationsTotal, dbOperationDuration)
	})
}
