package metrics

import (
	"fmt"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
)

type PrometheusConfig struct {
	Namespace string
	Endpoint  string
	Port      string
}

func InitPrometheus(config PrometheusConfig) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: config.Namespace,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Prometheus exporter: %v", err))
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle(config.Endpoint, pe)
		if err := http.ListenAndServe(config.Port, mux); err != nil {
			panic(fmt.Errorf("failed to run Prometheus /metrics endpoint: %v", err))
		}
	}()
}
