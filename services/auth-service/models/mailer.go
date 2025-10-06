package models

type EmailPayload struct {
	To       string         `json:"to"`
	Subject  string         `json:"subject"`
	Text     string         `json:"text,omitempty"`
	HTML     string         `json:"html,omitempty"`
	Template string         `json:"template,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
}
