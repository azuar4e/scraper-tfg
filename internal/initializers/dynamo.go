package initializers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DY *dynamodb.Client

func ConnectToDynamo() {
	// para ver conexion con localstack revisar historial de commits
	ctx := context.Background()
	dynamoConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}

	DY = dynamodb.NewFromConfig(dynamoConfig)
	log.Println("[+] Dynamo connection successful")
}
