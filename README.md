# Scraper

![Go](https://img.shields.io/badge/Go-1.24-blue)
![AWS](https://img.shields.io/badge/AWS-EKS%20%7C%20DynamoDB%20%7C%20SQS%20%7C%20SNS-orange)

The source code of this repository corresponds to the Scraper microservice for my TFG (bachelor's thesis). The scraper is the responsible for processing all user messages published at the SQS queue. It uses Playwright, a library for automatization of browsers, and the AWS SDK for Go, necessary for the integration with the managed services used in the architecture.

The complete explanation of the project can be found in my [TFG](https://azuar4e.github.io/en/posts/tfg) article on my blog.

## Overview

The scraper has the following features:

- Receives the jobs messages through the Simple Queue Service (SQS) queue.
- Searches for the price and name of the product.
- Reads the jobs information from DynamoDB table and updates the last price.
- Notifies users when the price reaches the desired value by the Simple Notification Service (SNS) topic.
- Deletes the jobs from the SQS queue, which are marked again as `ready` and re-scheduled by a lambda function.

## How it works?

The scraper has the following flow:

```Go
messageChan := make(chan types.Message)

for w := 1; w <= 5; w++ {
    go worker(w, messageChan)
}
...
for _, m := range output.Messages {
    messageChan <- m
}
```

We receive the messages from SQS and process them with five concurrent workers using **goroutines**.

```Go
for m := range messages {
    setJobStatus(m, "processing")
    handlers.ProcessMessageHandler(m)
    initializers.SQS.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(os.Getenv("SQS_QUEUE_URL")),
        ReceiptHandle: m.ReceiptHandle,
    })
}
```

Each worker set the status as `processing`, to avoid the lambda function re-scheduling a job that is being processed, process the message and deletes it from the queue.

Regarding the data processing, we configure the playwright browser, make the research and apply the necessary changes as follows:

```Go
pw, err := playwright.Run()
...
defer pw.Stop()
browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
    Headless: playwright.Bool(true),
})
...
defer browser.Close()
page, err := browser.NewPage(playwright.BrowserNewPageOptions{
    Locale:    playwright.String("es-ES"),
    UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/122.0"),
})
...
```

## Requirements

- Go 1.24+
- AWS credentials configured
