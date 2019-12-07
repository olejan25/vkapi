package vkapi

import "github.com/prometheus/client_golang/prometheus"

var (
	promRq = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "vk",
			Name:       "requests_dur",
			Help:       "vk API requests stats",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method"},
	)
	promRqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "vk",
			Name:      "requests_total",
			Help:      "vk API requests counter",
		},
		[]string{"method"},
	)
)

// InitProm - инициализация прометея
func InitProm() {
	prometheus.MustRegister(promRq)
	prometheus.MustRegister(promRqCount)
}
