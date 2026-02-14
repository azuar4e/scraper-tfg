package initializers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var SQS *sqs.Client

func ConnectToSQS() {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	SQS = sqs.NewFromConfig(sdkConfig)
	log.Println("[+] SQS connection successful")
	//creamos la cola por cli
}
