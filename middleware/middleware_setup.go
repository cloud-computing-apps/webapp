package middleware

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func MetricMiddleware(statsdClient *statsd.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counterMetricName := "api." + r.Method + "_" + r.URL.Path + ".count"
		timerMetricName := "api." + r.Method + "_" + r.URL.Path + ".latency"

		if err := statsdClient.Incr(counterMetricName, nil, 1); err != nil {
			log.Errorf("Error sending counter metric: %v", err)
		}

		metrics := httpsnoop.CaptureMetrics(next, w, r)
		latencyMs := metrics.Duration.Milliseconds()

		if err := statsdClient.Timing(timerMetricName, time.Duration(latencyMs)*time.Millisecond, nil, 1); err != nil {
			log.Errorf("Error sending timing metric: %v", err)
		}
	})
}

func WrapDBQuery(statsdClient *statsd.Client, queryName string, queryFunc func() error) error {
	startTime := time.Now()
	err := queryFunc()
	durationMs := time.Since(startTime).Milliseconds()
	metricName := "db." + queryName + ".latency"
	if err2 := statsdClient.Timing(metricName, time.Duration(durationMs)*time.Millisecond, nil, 1); err2 != nil {
		log.Errorf("Error sending DB timing metric: %v", err2)
	}
	return err
}

func WrapS3Call(statsdClient *statsd.Client, s3CallName string, s3CallFunc func() error) error {
	startTime := time.Now()
	err := s3CallFunc()
	durationMs := time.Since(startTime).Milliseconds()
	metricName := "s3." + s3CallName + ".latency"
	if err2 := statsdClient.Timing(metricName, time.Duration(durationMs)*time.Millisecond, nil, 1); err2 != nil {
		log.Error("Error sending S3 timing metric: %v", err2)
	}
	return err
}
