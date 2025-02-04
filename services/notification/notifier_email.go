package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"e2clicker.app/services/notification/openapi"
	"go.uber.org/fx"

	// Choose to use a deprecated library. The only other choices (xhit,
	// wneessen) are all inferior in its code quality.
	"gopkg.in/mail.v2"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
)

// EmailNotificationConfig is a user configuration for the email service.
type EmailNotificationConfig = openapi.EmailSubscription

// EmailService is a service for sending notifications via email.
type EmailService struct {
	dialer *mail.Dialer
	config *e2clickermodule.EmailSubmodule
	logger *slog.Logger
}

// NewEmailService creates a new email service.
func NewEmailService(config e2clickermodule.Notification, logger *slog.Logger, lc fx.Lifecycle) (*EmailService, error) {
	if config.Email == nil {
		return nil, nil
	}

	var mailConfig *e2clickermodule.EmailSubmodule

	switch value := config.Email.Value.(type) {
	case e2clickermodule.EmailSubmodule:
		mailConfig = &value
	case e2clickermodule.EmailPath:
		b, err := os.ReadFile(string(value))
		if err != nil {
			return nil, fmt.Errorf("cannot read email config file at %s: %w", value, err)
		}
		mailConfig = new(e2clickermodule.EmailSubmodule)
		if err := json.Unmarshal(b, mailConfig); err != nil {
			return nil, fmt.Errorf("cannot unmarshal email config at %s: %w", value, err)
		}
	default:
		panic("unreachable")
	}

	logger = logger.With(
		"notifier", "email",
		"email.host", mailConfig.SMTP.Host,
		"email.port", mailConfig.SMTP.Port,
		"email.secure", mailConfig.SMTP.Secure,
	)

	dialer := mail.NewDialer(
		mailConfig.SMTP.Host,
		mailConfig.SMTP.Port,
		mailConfig.SMTP.Auth.Username,
		mailConfig.SMTP.Auth.Password,
	)
	if mailConfig.SMTP.Secure {
		dialer.StartTLSPolicy = mail.MandatoryStartTLS
	}
	dialer.LocalName = "e2clicker"
	dialer.Timeout = 30 * time.Second
	dialer.RetryFailure = true

	return &EmailService{
		dialer: dialer,
		config: mailConfig,
		logger: logger,
	}, nil
}

func (s EmailService) Notify(ctx context.Context, n Notification, config EmailNotificationConfig) error {
	optstr := func(opt *string) string {
		if opt == nil {
			return ""
		}
		return *opt
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", s.config.From)
	msg.SetAddressHeader("To", string(config.Address), optstr(config.Name))
	msg.SetHeader("Subject", n.Message.Title)
	msg.SetBody("text/plain", n.Message.Message)
	// TODO: text/html support

	s.logger.Debug(
		"sending email",
		"from", s.config.From,
		"notification", n.Type)

	if err := s.dialer.DialAndSend(msg); err != nil {
		s.logger.Error(
			"failed to send email",
			"from", s.config.From,
			"notification", n.Type,
			"err", err)

		return fmt.Errorf("cannot send email: %w", err)
	}

	return nil
}

func stringifyEmails[T ~string](emails []T) []string {
	result := make([]string, 0, len(emails))
	for _, email := range emails {
		result = append(result, string(email))
	}
	return result
}
