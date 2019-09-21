package configuration

import (
	"email-sender/instrumentation/metrics"
	"email-sender/instrumentation/trace"
	"email-sender/logger"
	"email-sender/providers"
	"email-sender/server"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type EmailSender struct {
	Logger     logger.Config
	Server     server.Config
	Providers  providers.Config
	Jaeger     trace.JaegerConfig
	Prometheus metrics.PrometheusConfig
}

func LoadConfig(path string) EmailSender {
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := &EmailSender{}
	err := viper.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode into config struct, %v", err))
	}

	return *conf
}
