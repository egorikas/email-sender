package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"email-sender/providers"
	"email-sender/server"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewEmail(t *testing.T) {
	controller := NewEmail(zap.NewExample(), providers.NewSenderMock(t))

	require.NotNil(t, controller)
}

func TestEmails_SendEmail(t *testing.T) {
	t.Run("invalid input returns error. invalid json", func(t *testing.T) {
		sender := providers.NewSenderMock(t)
		controller := NewEmail(zap.NewExample(), sender)

		e := echo.New()
		e.Validator = server.NewValidator()
		req := httptest.NewRequest(
			http.MethodPost,
			"/emails",
			strings.NewReader("dummy data"),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		mockCtx := e.NewContext(req, rec)

		err := controller.SendEmail(mockCtx)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid input returns error. validation failed", func(t *testing.T) {
		sender := providers.NewSenderMock(t)
		controller := NewEmail(zap.NewExample(), sender)

		json := `
			{
			  "to": [
			    "invalid_email"
			  ],
			  "from": "from@from.com",
			  "subject": "subject",
			  "body": "body"
			}
		`
		e := echo.New()
		e.Validator = server.NewValidator()
		req := httptest.NewRequest(
			http.MethodPost,
			"/emails",
			strings.NewReader(json),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		mockCtx := e.NewContext(req, rec)

		err := controller.SendEmail(mockCtx)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("works well with valid input", func(t *testing.T) {
		sender := providers.NewSenderMock(t)
		sender.SendMock.Set(func(ctx context.Context, email providers.Email) (err error) {
			require.Equal(t, "from@from.com", email.From)
			require.Equal(t, "subject", email.Subject)
			require.Equal(t, "body", email.Body)

			require.Equal(t, 2, len(email.To))
			require.Equal(t, "to@to.com", email.To[0])
			require.Equal(t, "to2@to2.com", email.To[1])

			return nil
		})
		controller := NewEmail(zap.NewExample(), sender)

		json := `
			{
			  "to": [
			    "to@to.com",
				"to2@to2.com"
			  ],
			  "from": "from@from.com",
			  "subject": "subject",
			  "body": "body"
			}
		`
		e := echo.New()
		e.Validator = server.NewValidator()
		req := httptest.NewRequest(
			http.MethodPost,
			"/emails",
			strings.NewReader(json),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		mockCtx := e.NewContext(req, rec)

		err := controller.SendEmail(mockCtx)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "{}\n", rec.Body.String())
		sender.MinimockFinish()
	})
}
