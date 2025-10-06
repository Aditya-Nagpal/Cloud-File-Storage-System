package models

import "encoding/json"

// NotificationMessage is the top-level SQS message format
type NotificationMessage struct {
	Type    string          `json:"type"`    // e.g. "EMAIL", "SMS", "IN_APP"
	Payload json.RawMessage `json:"payload"` // raw payload per type
	Meta    map[string]any  `json:"meta,omitempty"`
}

// EmailPayload represents the email-specific payload
type EmailPayload struct {
	To       string         `json:"to" binding:"required,email"`
	Subject  string         `json:"subject"`
	Text     string         `json:"text,omitempty"`
	HTML     string         `json:"html,omitempty"`
	Template string         `json:"template,omitempty"`
	Data     map[string]any `json:"data,omitempty"` // optional template data
}
