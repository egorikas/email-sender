package providers

import (
	"context"
	"math/rand"
	"time"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

type dummy struct {
	log *zap.Logger
}

func (d *dummy) Send(ctx context.Context, email Email) error {
	d.log.Debug("start to process an email with dummy",
		zap.Strings("to", email.To),
		zap.String("from", email.From),
		zap.String("subject", email.Subject),
		zap.String("body", email.Body))

	ctx, span := trace.StartSpan(ctx, "Dummy.Send")
	defer span.End()

	sleepTime := time.Duration(random(100, 5000)) * time.Millisecond
	time.Sleep(sleepTime)

	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
