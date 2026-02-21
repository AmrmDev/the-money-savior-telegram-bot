package models

import "time"

type Expense struct {
	UserID    int64     `dynamodbav:"user_id"`
	ChatID    int64     `dynamodbav:"chat_id"`
	Username  string    `dynamodbav:"username"`
	Amount    float64   `dynamodbav:"amount"`
	Category  string    `dynamodbav:"category"`
	Method    string    `dynamodbav:"method"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	ExpenseID string    `dynamodbav:"expense_id"` // sort key: user_id#timestamp
	SeqID     int       `dynamodbav:"seq_id"`      // sequential ID: 1, 2, 3...
}