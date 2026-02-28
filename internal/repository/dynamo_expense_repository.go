package repository

import (
	"context"

	"money-telegram-bot/internal/models"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoExpenseRepository struct {
	client *dynamodb.Client
}

func NewDynamoExpenseRepository(client *dynamodb.Client) *DynamoExpenseRepository {
	return &DynamoExpenseRepository{
		client: client,
	}
}

var _ ExpenseRepository = (*DynamoExpenseRepository)(nil)