package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scraper/internal/handlers"
	"scraper/internal/initializers"
	"scraper/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	typeDynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		setJobStatus(m, "processing")
		handlers.ProcessMessageHandler(m)
		initializers.SQS.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(os.Getenv("SQS_QUEUE_URL")),
			ReceiptHandle: m.ReceiptHandle,
		})
	}
}

func setJobStatus(m types.Message, status string) {
	var job models.Job
	json.Unmarshal([]byte(*m.Body), &job)
	id, _ := job.ID.Int64()

	initializers.DY.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]typeDynamo.AttributeValue{
			"PK": &typeDynamo.AttributeValueMemberN{Value: fmt.Sprintf("%d", job.UserID)},
			"SK": &typeDynamo.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
		UpdateExpression:         aws.String("SET #s = :s"),
		ExpressionAttributeNames: map[string]string{"#s": "status"},
		ExpressionAttributeValues: map[string]typeDynamo.AttributeValue{
			":s": &typeDynamo.AttributeValueMemberS{Value: status},
		},
	})
}
