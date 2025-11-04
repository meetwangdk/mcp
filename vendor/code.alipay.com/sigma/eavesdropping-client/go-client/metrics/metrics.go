package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricNamePrefix = "eavesdropping_client_"
)

var (
	// 统计查询方法耗时
	QueryMethodDurationMilliSeconds = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       MetricNamePrefix + "query_method_duration_milliseconds",
			Help:       "how long a query method taken",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001},
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(QueryMethodDurationMilliSeconds)
}

func ObserveQueryMethodDuration(method string, begin time.Time) {
	cost := float64(time.Since(begin).Nanoseconds() / time.Millisecond.Nanoseconds())
	QueryMethodDurationMilliSeconds.WithLabelValues(method).Observe(cost)
}
