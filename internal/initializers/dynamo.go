package initializers

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DY *dynamodb.Client

func ConnectToDynamo() {
	// cambiar config a la default cuando se despliegue en aws
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}

	endpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	opts := []func(*dynamodb.Options){}
	if endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	DY = dynamodb.NewFromConfig(cfg, opts...)
	log.Println("[+] Dynamo connection successful")
}
