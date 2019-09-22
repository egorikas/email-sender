package providers

import (
	"context"
	"go.opencensus.io/trace"
	"math/rand"
	"time"
)

type dummy struct {
}

func (d *dummy) Send(ctx context.Context, email Email) error {
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
