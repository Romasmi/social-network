package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	WorkerProcessedEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "social_network_worker_processed_events_total",
			Help: "Total number of events processed by the worker.",
		},
		[]string{"event_type", "status"},
	)

	WorkerProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "social_network_worker_processing_duration_seconds",
			Help:    "Duration of event processing in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)
)
