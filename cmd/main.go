package main

import (
	"context"
	"log"
	"os"
	"scraper/internal/handlers"
	"scraper/internal/initializers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDynamo()
	initializers.ConnectToSQS()
	initializers.ConnectToSNS()
}

func main() {
	messageChan := make(chan types.Message)

	for w := 1; w <= 5; w++ {
		go worker(w, messageChan)
	}

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

		for _, m := range output.Messages {
			messageChan <- m
		}
	}
}

func worker(id int, messages <-chan types.Message) {
	for m := range messages {
		handlers.ProcessMessageHandler(m)
		initializers.SQS.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(os.Getenv("SQS_QUEUE_URL")),
			ReceiptHandle: m.ReceiptHandle,
		})
	}
}
