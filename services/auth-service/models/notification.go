package models

import "encoding/json"

type NotificationMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Meta    map[string]any  `json:"meta,omitempty"`
}
