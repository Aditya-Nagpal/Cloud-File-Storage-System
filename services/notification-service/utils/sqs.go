package utils

import (
	"context"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/config"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	Client *sqs.Client
	Queue  string
}

func NewSQSClient(ctx context.Context) (*SQSClient, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.AppConfig.AWSRegion))
	if err != nil {
		log.Printf("error loading aws config: %v", err)
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return &SQSClient{
		Client: client,
		Queue:  config.AppConfig.SqsQueueUrl,
	}, nil
}

// ReceiveMessages - long polling wrapper
func (s *SQSClient) ReceiveMessages(ctx context.Context, maxMessages, waitTime int32) ([]types.Message, error) {
	output, err := s.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              &s.Queue,
		MaxNumberOfMessages:   maxMessages,
		WaitTimeSeconds:       waitTime,
		VisibilityTimeout:     int32(60),
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		log.Printf("error receiving messages from sqs: %v", err)
		return nil, err
	}
	return output.Messages, nil
}

func (s *SQSClient) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := s.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.Queue,
		ReceiptHandle: &receiptHandle,
	})
	if err != nil {
		log.Printf("error deleting message from sqs: %v", err)
		return err
	}
	return nil
}
