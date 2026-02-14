package models

import "time"

type Job struct {
	ID          string    `json:"id" dynamodbav:"id"`
	UserID      uint      `json:"user_id" dynamodbav:"user_id"`
	URL         string    `json:"url" dynamodbav:"url" binding:"required,url"`
	TargetPrice float64   `json:"target_price" dynamodbav:"target_price" binding:"required"`
	LastPrice   float64   `json:"last_price" dynamodbav:"last_price"`
	Status      string    `json:"status" dynamodbav:"status"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
}
