package providers

import (
	"context"
	provider "github.com/mailgun/mailgun-go/v3"
	"go.opencensus.io/trace"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"time"
)

type MailGunConfig struct {
	ApiKey  string
	Domain  string
	Timeout time.Duration
	RPS     int
}

type mailgun struct {
	log     *zap.Logger
	sdk     provider.Mailgun
	limiter ratelimit.Limiter
	config  MailGunConfig
}

func newMailgun(log *zap.Logger, config MailGunConfig) *mailgun {
	return &mailgun{
		log:     log,
		config:  config,
		sdk:     provider.NewMailgun(config.Domain, config.ApiKey),
		limiter: ratelimit.New(config.RPS),
	}
}

func (m *mailgun) Send(ctx context.Context, email Email) error {
	ctx, span := trace.StartSpan(ctx, "Mailgun.Send")
	defer span.End()

	m.log.Debug("start to process an email with mailgun",
		zap.Strings("to", email.To),
		zap.String("from", email.From),
		zap.String("subject", email.Subject),
		zap.String("body", email.Body))

	message := m.sdk.NewMessage(email.From, email.Subject, email.Body, email.To...)

	ctx, cancel := context.WithTimeout(ctx, m.config.Timeout)
	defer cancel()

	// Send the message	with a timeout
	m.limiter.Take()
	resp, id, err := m.sdk.Send(ctx, message)
	if err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		m.log.Error("failed to send an email", zap.Error(err))
		return err
	}

	m.log.Debug("email has been sent",
		zap.String("status", resp),
		zap.String("id", id))

	return nil
}
