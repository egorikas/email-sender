package providers

import (
	"context"
	sp "github.com/SparkPost/gosparkpost"
	"go.opencensus.io/trace"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"time"
)

type SparkpostConfig struct {
	ApiKey  string
	Timeout time.Duration
	RPS     int
}

type sparkPost struct {
	log     *zap.Logger
	sdk     *sp.Client
	config  SparkpostConfig
	limiter ratelimit.Limiter
}

func newSparkPost(logger *zap.Logger, config SparkpostConfig) *sparkPost {
	var sparky sp.Client
	err := sparky.Init(&sp.Config{ApiKey: config.ApiKey})
	if err != nil {
		panic("failed to init sdk client")
	}
	sparky.Client.Timeout = config.Timeout
	return &sparkPost{
		config: config,
		log:    logger, sdk: &sparky,
		limiter: ratelimit.New(config.RPS),
	}
}

func (s *sparkPost) Send(ctx context.Context, email Email) error {
	_, span := trace.StartSpan(ctx, "Sparkpost.Send")
	defer span.End()

	isSandbox := true
	tx := &sp.Transmission{
		Recipients: email.To,
		Options:    &sp.TxOptions{Sandbox: &isSandbox},
		Content: sp.Content{
			Text:    email.Body,
			From:    email.From,
			Subject: email.Subject,
		},
	}

	s.limiter.Take()
	id, _, err := s.sdk.Send(tx)
	if err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		s.log.Error("failed to send an email", zap.Error(err))
		return err
	}
	s.log.Debug("email has been sent",
		zap.String("id", id))

	return nil
}
