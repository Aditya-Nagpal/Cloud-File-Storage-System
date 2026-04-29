package models

type Payload struct {
	InternalID int64  `json:"internal_id"`
	S3Key      string `json:"s3_key"`
}
