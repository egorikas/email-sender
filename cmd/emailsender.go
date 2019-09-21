package main

import (
	"email-sender/configuration"
	"email-sender/controllers"
	"email-sender/instrumentation/metrics"
	"email-sender/instrumentation/trace"
	"email-sender/logger"
	"email-sender/providers"
	"email-sender/server"
	"flag"
	"go.uber.org/zap"
)

func main() {
	params := parseInputParams()
	config := configuration.LoadConfig(params.configPath)
	log := logger.InitLogger(config.Logger)

	trace.InitJaeger(config.Jaeger)
	metrics.InitPrometheus(config.Prometheus)

	e := server.New(config.Server)
	v1 := e.Group("/api/v1")

	emailController := controllers.NewEmail(log, providers.NewSender(log, config.Providers))
	emailController.Register(v1)

	err := e.Start(config.Server.Port)
	if err != nil {
		log.Fatal("echo start returns error", zap.Error(err))
	}
}

type inputParams struct {
	configPath string
}

func parseInputParams() inputParams {
	configPath := flag.String("config", "", "path to the config file")

	flag.Parse()
	if configPath == nil || len(*configPath) == 0 {
		panic("config path isn't provided")
	}
	return inputParams{
		configPath: *configPath,
	}
}
