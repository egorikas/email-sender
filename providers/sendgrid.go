package providers

import (
	"context"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opencensus.io/trace"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type SendGridConfig struct {
	ApiKey   string
	Endpoint string
	Host     string
	Timeout  time.Duration
	RPS      int
}

type sendGrid struct {
	log     *zap.Logger
	limiter ratelimit.Limiter
	config  SendGridConfig
}

func newSendGrid(log *zap.Logger, config SendGridConfig) *sendGrid {
	return &sendGrid{
		log:     log,
		config:  config,
		limiter: ratelimit.New(config.RPS),
	}
}

func (s *sendGrid) Send(ctx context.Context, email Email) error {
	_, span := trace.StartSpan(ctx, "Sendgrid.Send")
	defer span.End()

	s.log.Debug("start to process an email with sendgrid",
		zap.Strings("to", email.To),
		zap.String("from", email.From),
		zap.String("subject", email.Subject),
		zap.String("body", email.Body))

	from := mail.NewEmail("", email.From)
	to := mail.NewEmail("", email.To[0])
	content := mail.NewContent("text/plain", email.Body)
	message := mail.NewV3MailInit(from, email.Subject, to, content)
	if len(email.To) > 1 {
		for i := 1; i < len(email.To); i++ {
			newAddress := &mail.Email{Address: email.To[i]}
			message.Personalizations[0].To = append(message.Personalizations[0].To, newAddress)
		}
	}

	req := sendgrid.NewSendClient(s.config.ApiKey).Request
	req.Body = mail.GetRequestBody(message)
	var custom rest.Client
	custom.HTTPClient = &http.Client{Timeout: s.config.Timeout}

	s.limiter.Take()

	response, err := custom.Send(req)
	if err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		s.log.Error("failed to send an email", zap.Error(err))
		return err
	}

	s.log.Debug("email has been sent",
		zap.Int("status", response.StatusCode),
		zap.String("body", response.Body))
	return nil
}
