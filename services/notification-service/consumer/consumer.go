package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/utils"
)

func Start(ctx context.Context, sqsClient *utils.SQSClient, proc *handlers.Processor) {
	waitSeconds := int32(config.AppConfig.SqsWaitTimeSeconds)
	maxMsgs := int32(config.AppConfig.SqsMaxMessages)

	log.Println("SQS consumer started, listening to", config.AppConfig.SqsQueueUrl)

	for {
		select {
		case <-ctx.Done():
			log.Println("SQS consumer shutting down")
			return
		default:
			// receive messages
			msgs, err := sqsClient.ReceiveMessages(ctx, maxMsgs, waitSeconds)
			if err != nil {
				log.Printf("sqs receive error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			if len(msgs) == 0 {
				continue
			}

			for _, m := range msgs {
				// parse message body
				var nm models.NotificationMessage
				if err := json.Unmarshal([]byte(*m.Body), &nm); err != nil {
					log.Printf("invalid message body, deleting message: %v", err)
					// delete malformed message to avoid poison
					_ = sqsClient.DeleteMessage(ctx, *m.ReceiptHandle)
					continue
				}

				// process
				if err := proc.ProcessMessage(ctx, nm); err != nil {
					log.Printf("processing failed: %v - leaving message for retry", err)
					// do not delete -> SQS will retry / go to DLQ
					continue
				}

				// success -> delete message
				if err := sqsClient.DeleteMessage(ctx, *m.ReceiptHandle); err != nil {
					log.Printf("failed to delete message: %v", err)
				}
			}
		}
	}
}
