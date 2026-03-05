package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"money-telegram-bot/internal/models"
	"money-telegram-bot/internal/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoExpenseRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoExpenseRepository(client *dynamodb.Client) *DynamoExpenseRepository {
	return &DynamoExpenseRepository{
		client:    client,
		tableName: os.Getenv("TABLE_NAME"),
	}
}

func (r *DynamoExpenseRepository) Save(
	ctx context.Context,
	expense *models.Expense,
) error {

	expense.ExpenseID = utils.GenerateExpenseID()

	av, err := attributevalue.MarshalMap(expense)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	return err
}

func (r *DynamoExpenseRepository) FindByUser(
	ctx context.Context,
	userID string,
) ([]models.Expense, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberN{Value: userID},
		},
		ScanIndexForward: aws.Bool(true),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}

	var expenses []models.Expense
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *DynamoExpenseRepository) FindByID(
	ctx context.Context,
	userID,
	expenseID string,
) (*models.Expense, error) {

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":    &types.AttributeValueMemberN{Value: userID},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	})

	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, fmt.Errorf("expense not found")
	}

	var expense models.Expense
	if err := attributevalue.UnmarshalMap(out.Item, &expense); err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *DynamoExpenseRepository) DeleteAllExpenses(ctx context.Context, userID int64) error {
	userIDStr := fmt.Sprintf("%d", userID)

	expenses, err := r.FindByUser(ctx, userIDStr)
	if err != nil {
		return err
	}

	for _, expense := range expenses {
		_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"user_id":    &types.AttributeValueMemberN{Value: userIDStr},
				"expense_id": &types.AttributeValueMemberS{Value: expense.ExpenseID},
			},
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete expense: %v", err)
			return err
		}
	}

	log.Printf("[INFO] All expenses deleted | userID=%d | count=%d", userID, len(expenses))
	return nil
}

func (r *DynamoExpenseRepository) DeleteByID(ctx context.Context, userID, expenseID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":    &types.AttributeValueMemberN{Value: userID},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	})
	if err != nil {
		log.Printf("[ERROR] Failed to delete expense %s: %v", expenseID, err)
		return err
	}
	return nil
}

var _ ExpenseRepository = (*DynamoExpenseRepository)(nil)
