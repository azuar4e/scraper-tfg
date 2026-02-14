package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"scraper/internal/initializers"
	"scraper/internal/models"
	"scraper/internal/selectors"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	typeDynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	types2 "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/playwright-community/playwright-go"
)

func ProcessMessageHandler(message types.Message) {

	fmt.Println("processing message")
	var job models.Job
	err := json.Unmarshal([]byte(*message.Body), &job)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Starting scraper...")

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	defer pw.Stop()

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})

	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	defer browser.Close()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Locale:    playwright.String("es-ES"),
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/122.0"),
	})
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	//logica dependiendo del tipo de url
	titleTag, priceTag := selectors.Selector(job.URL)

	if titleTag == "" || priceTag == "" {
		log.Fatalf("unsupported url: %s", job.URL)
	}

	if _, err := page.Goto(job.URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	}); err != nil {
		log.Fatalf("could not navigate to page: %v", err)
	}

	_ = page.Locator(titleTag).First().WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(15000),
	})

	title := ""
	if t, err := page.Locator(titleTag).InnerText(); err == nil {
		title = strings.TrimSpace(t)
	}

	_ = page.Locator(priceTag).First().WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(15000),
	})

	price := ""
	p, err := page.Locator(priceTag).First().InnerText()
	if err != nil {
		log.Println("price err:", err)
	} else {
		price = strings.TrimSpace(p)
	}

	fmt.Println(title, price)

	// añadir el precio al registro de dynamo cambiar estado a active
	// si el precio es menor lo dejamos en notified y ya si eso hare
	//otra logica para reencolar los q ya han sido notificados

	status := "active"
	lastPrice, _ := parsePrice(price)

	if job.TargetPrice >= lastPrice {
		status = "notified"
		_, err = initializers.SNS.Publish(context.TODO(), &sns.PublishInput{
			TopicArn: aws.String(os.Getenv("SNS_TOPIC_ARN")),
			Message:  aws.String(fmt.Sprintf("The product in the following url now costs %v, which is below your target price of %v.\n\n%v", price, job.TargetPrice, job.URL)),
			Subject:  aws.String(fmt.Sprintf("¡Price Alert! The product price has dropped below %v.", job.TargetPrice)),
			MessageAttributes: map[string]types2.MessageAttributeValue{
				"user_id": {
					DataType:    aws.String("Number"),
					StringValue: aws.String(fmt.Sprintf("%d", job.UserID)),
				},
			},
		})

		if err != nil {
			log.Println(err)
			return
		}
	}

	_, err = initializers.DY.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("PriceAlerts"),
		Key: map[string]typeDynamo.AttributeValue{
			"PK": &typeDynamo.AttributeValueMemberN{Value: fmt.Sprintf("%d", job.ID)},
			"SK": &typeDynamo.AttributeValueMemberN{Value: fmt.Sprintf("%d", job.UserID)},
		},
		UpdateExpression:         aws.String("SET #s = :s, last_price = :p, updated_at = :t"),
		ExpressionAttributeNames: map[string]string{"#status": "status"},
		ExpressionAttributeValues: map[string]typeDynamo.AttributeValue{
			":s": &typeDynamo.AttributeValueMemberS{Value: status},
			":p": &typeDynamo.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", lastPrice)},
			":t": &typeDynamo.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	})

	if err != nil {
		log.Println(err)
		return
	}
}

func parsePrice(price string) (float64, error) {
	re := regexp.MustCompile(`[\d.,]+`)

	match := re.FindString(price)

	if match == "" {
		return 0, fmt.Errorf("No price found")
	}

	if strings.Contains(match, ",") && strings.Contains(match, ".") {
		match = strings.ReplaceAll(match, ",", "")
	}

	if strings.Count(match, ",") == 1 && !strings.Contains(match, ".") {
		match = strings.ReplaceAll(match, ",", ".")
	}

	return strconv.ParseFloat(match, 64)
}
