package models

import (
	"strconv"
	"time"
)

// JobDynamoItem es lo que se guarda en DynamoDB. Solo claves + atributos; sin id/user_id extra (PK y SK son eso).
type JobDynamoItem struct {
	PK          int64   `dynamodbav:"PK"` //user_id
	SK          int64   `dynamodbav:"SK"` //job_id
	URL         string  `dynamodbav:"url"`
	TargetPrice float64 `dynamodbav:"target_price"`
	LastPrice   float64 `dynamodbav:"last_price"`
	Status      string  `dynamodbav:"status"`
	CreatedAt   string  `dynamodbav:"created_at"`
	UpdatedAt   string  `dynamodbav:"updated_at"`
}

// ToJob convierte el item de Dynamo (PK, SK, atributos) al struct Job que devolvemos en JSON (id, user_id, url, ...).
func (d JobDynamoItem) ToJob() Job {
	created, _ := time.Parse(time.RFC3339, d.CreatedAt)
	updated, _ := time.Parse(time.RFC3339, d.UpdatedAt)
	return Job{
		ID:          strconv.FormatInt(d.SK, 10),
		UserID:      uint(d.PK),
		URL:         d.URL,
		TargetPrice: d.TargetPrice,
		LastPrice:   d.LastPrice,
		Status:      d.Status,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}
}
