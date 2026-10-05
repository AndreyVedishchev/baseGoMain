package db

import "github.com/prometheus/client_golang/prometheus"

// newMetrics создаёт метрики БД и регистрирует их в переданном registry.
// Каждый Storage получает собственный registry, повторная регистрация при создании нескольких Storage не конфликтует.
func newMetrics(reg *prometheus.Registry) (*prometheus.CounterVec, *prometheus.HistogramVec) {
	dbOperationsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "Number of database operations by type and result",
		},
		[]string{"operation", "status"},
	)

	dbOperationDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Duration of database operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	reg.MustRegister(dbOperationsTotal, dbOperationDuration)

	return dbOperationsTotal, dbOperationDuration
}
