package models

type EmailPayload struct {
	To       string         `json:"to"`
	Subject  string         `json:"subject"`
	Template string         `json:"template,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
}
