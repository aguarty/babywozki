package main

import "github.com/prometheus/client_golang/prometheus"

var (
	requestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gb",
			Subsystem: "dp_terms",
			Name:      "total_requests",
			Help:      "Total number of HTTP requests",
		},
		[]string{"path"},
	)

	responsesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gb",
			Subsystem: "dp_terms",
			Name:      "total_responses",
			Help:      "Total number of HTTP responses",
		},
		[]string{"path", "code"},
	)

	latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gb",
			Subsystem: "dp_terms",
			Name:      "latency",
			Help:      "Request-response latency",
			Buckets:   []float64{0.0001, 0.001, 0.1, 0.2, 0.5, 1.000, 3.000, 5.000},
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(requestsCounter)
	prometheus.MustRegister(responsesCounter)
	prometheus.MustRegister(latency)
}
