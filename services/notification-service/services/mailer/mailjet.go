package mailer

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/config"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
)

// Mailer interface so you can swap implementations
type Mailer interface {
	SendEmail(ctx context.Context, to, subject, text, html string) error
}

type MailjetMailer struct {
	client    *mailjet.Client
	fromEmail string
	fromName  string
}

func NewMailjetMailer() *MailjetMailer {
	c := mailjet.NewMailjetClient(
		config.AppConfig.MailjetApiKeyPublic,
		config.AppConfig.MailjetApiKeyPrivate,
	)
	return &MailjetMailer{
		client:    c,
		fromEmail: config.AppConfig.MailjetSenderEmail,
		fromName:  config.AppConfig.MailjetSenderName,
	}
}

func (m *MailjetMailer) SendEmail(ctx context.Context, to, subject, text, html string) error {
	if to == "" || subject == "" {
		return errors.New("to and subject required")
	}

	message := mailjet.MessagesV31{
		Info: []mailjet.InfoMessagesV31{
			{
				From: &mailjet.RecipientV31{
					Email: m.fromEmail,
					Name:  m.fromName,
				},
				To: &mailjet.RecipientsV31{
					mailjet.RecipientV31{Email: to},
				},
				Subject:  subject,
				TextPart: text,
				HTMLPart: html,
			},
		},
	}

	_, err := m.client.SendMailV31(&message)
	if err != nil {
		return fmt.Errorf("mailjet send error: %w", err)
	}
	return nil
}
