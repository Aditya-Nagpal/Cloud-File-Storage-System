package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

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

	template := ep.Template
	if template == "" {
		return errors.New("empty email template")
	}

	subject := ep.Subject

	data := ep.Data
	for key, val := range data {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("validation error: data map contains an empty key")
		}

		if strVal, ok := val.(string); ok {
			if strings.TrimSpace(strVal) == "" {
				return fmt.Errorf("validation error: data map key '%s' has an empty value", key)
			}
		}
	}

	// send via mailer
	if err := p.Mailer.SendEmail(ctx, to, subject, template, data); err != nil {
		log.Printf("email send error: %v", err)
		return err
	}
	return nil
}
