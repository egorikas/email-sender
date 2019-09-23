package providers

import (
	"context"
	"go.uber.org/zap"
)

type Type int

const (
	Mailgun Type = iota
	SendGrid
	Sparkpost
	Dummy
)

type Config struct {
	Current   Type
	MailGun   MailGunConfig
	SendGrid  SendGridConfig
	Sparkpost SparkpostConfig
}

type Email struct {
	To      []string
	From    string
	Subject string
	Body    string
}

type Sender interface {
	Send(ctx context.Context, email Email) error
}

func NewSender(log *zap.Logger, config Config) Sender {
	switch config.Current {
	case Mailgun:
		return newMailgun(log, config.MailGun)
	case SendGrid:
		return newSendGrid(log, config.SendGrid)
	case Sparkpost:
		return newSparkPost(log, config.Sparkpost)
	case Dummy:
		return &dummy{log: log}
	default:
		log.Panic("unknown provider")
	}

	return nil
}
