package controllers

import (
	"context"
	"email-sender/providers"
	"email-sender/server"
	"github.com/labstack/echo/v4"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"net/http"
)

type Emails struct {
	log    *zap.Logger
	sender providers.Sender
}

func NewEmail(log *zap.Logger, sender providers.Sender) *Emails {
	return &Emails{log: log, sender: sender}
}

func (c *Emails) Register(v1 *echo.Group) {
	v1.POST("/emails", c.SendEmail)
}

type newEmail struct {
	To      []string `json:"to" validate:"min=1,max=1000,dive,required,email"`
	From    string   `json:"from" validate:"required,email"`
	Subject string   `json:"subject" validate:"max=78"`
	Body    string   `json:"body"`
}

func (c *Emails) SendEmail(reqCtx echo.Context) error {
	ctx, span := trace.StartSpan(context.Background(), "Emails.SendEmail")
	defer span.End()

	stats.Record(ctx, statInputRequestCount.M(1))

	newEmail := new(newEmail)
	if err := reqCtx.Bind(newEmail); err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		stats.Record(ctx, statFailedRequestCount.M(1))

		return reqCtx.JSON(http.StatusBadRequest, server.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	if err := reqCtx.Validate(newEmail); err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		stats.Record(ctx, statFailedRequestCount.M(1))

		return reqCtx.JSON(http.StatusBadRequest, server.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	err := c.sender.Send(ctx, providers.Email{
		To:      newEmail.To,
		From:    newEmail.From,
		Subject: newEmail.Subject,
		Body:    newEmail.Body,
	})
	if err != nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())

		stats.Record(ctx, statFailedRequestCount.M(1))

		c.log.Error("request failed", zap.Error(err))

		return reqCtx.JSON(http.StatusBadRequest, server.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	stats.Record(ctx, statSuccessedRequestCount.M(1))
	return reqCtx.JSON(http.StatusOK, nil)
}
