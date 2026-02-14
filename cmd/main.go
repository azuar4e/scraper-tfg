package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"scraper/internal/handlers"
	"scraper/internal/initializers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDynamo()
	initializers.ConnectToSQS()
	initializers.ConnectToSNS()
}

func main() {
	for {
		output, err := initializers.SQS.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(os.Getenv("SQS_QUEUE_URL")),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})

		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}

		fmt.Println("tenemos el mensaje")

		for _, m := range output.Messages {
			fmt.Println("dentro del for")
			handlers.ProcessMessageHandler(m)
			initializers.SQS.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(os.Getenv("SQS_QUEUE_URL")),
				ReceiptHandle: m.ReceiptHandle,
			})
		}
	}
}
