package sqs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var client *sqs.Client
var queueURL string

// Initialize SQS client
func InitSQS() {
	cfg, err := awscfg.LoadDefaultConfig(
		context.TODO(),
		awscfg.WithRegion(config.AppConfig.AWSRegion),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.AppConfig.AWSAccessKeyId,
				config.AppConfig.AWSSecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		log.Printf("error loading aws config: %v", err)
	}

	client = sqs.NewFromConfig(cfg)
	queueURL = config.AppConfig.SqsQueueUrl
	log.Println("SQS initialized for region:", config.AppConfig.AWSRegion)
}

func Publish(ctx context.Context, message models.NotificationMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		log.Printf("Failed to publish message to SQS: %v", err)
		return err
	}

	log.Println("Message published to SQS successfully")
	return nil
}

// Publish sends a message to the queue
func PublishOTP(ctx context.Context, toEmail, otp, flowId string) error {
	ep := models.EmailPayload{
		To:       toEmail,
		Subject:  "Your password reset OTP",
		Template: "forgot_password",
		Data: map[string]any{
			"OTP":    otp,
			"FlowId": flowId,
		},
	}

	payloadBytes, err := json.Marshal(ep)
	if err != nil {
		return err
	}

	msg := models.NotificationMessage{
		Type:    "EMAIL",
		Payload: json.RawMessage(payloadBytes),
		Meta: map[string]any{
			"FlowId": flowId,
			"Source": "auth-service",
		},
	}

	return Publish(ctx, msg)
}
