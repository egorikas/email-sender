package trace

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"fmt"
	"go.opencensus.io/trace"
)

type JaegerConfig struct {
	ServiceName   string
	AgentEndpoint string
}

func InitJaeger(config JaegerConfig) {
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: config.AgentEndpoint,
		Process: jaeger.Process{
			ServiceName: config.ServiceName,
		},
	})
	if err != nil {
		panic(fmt.Errorf("Failed to create the Jaeger exporter: %v", err))
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
}
