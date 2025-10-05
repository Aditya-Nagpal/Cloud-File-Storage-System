package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/services/mailer"
)

// Processor holds channel implementations
type Processor struct {
	Mailer mailer.Mailer
}

func NewProcessor(m mailer.Mailer) *Processor {
	return &Processor{Mailer: m}
}

// ProcessMessage - receives raw SQS body (already unmarshalled to models.NotificationMessage)
func (p *Processor) ProcessMessage(ctx context.Context, msg models.NotificationMessage) error {
	switch msg.Type {
	case "EMAIL":
		var ep models.EmailPayload
		if err := json.Unmarshal(msg.Payload, &ep); err != nil {
			return err
		}
		return p.handleEmail(ctx, ep)
	// case "SMS": add sms handler
	// case "IN_APP": add in-app handler
	default:
		return errors.New("unsupported notification type: " + msg.Type)
	}
}

func (p *Processor) handleEmail(ctx context.Context, ep models.EmailPayload) error {
	// If template is present you could render it here
	// For now, prefer HTML over text if provided
	to := ep.To
	if to == "" {
		return errors.New("empty email recipient")
	}

	text := ep.Text
	html := ep.HTML
	if html == "" && text == "" {
		return errors.New("empty email body")
	}
	subject := ep.Subject

	// send via mailer
	if err := p.Mailer.SendEmail(ctx, to, subject, text, html); err != nil {
		log.Printf("email send error: %v", err)
		return err
	}
	return nil
}
