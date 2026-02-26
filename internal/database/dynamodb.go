package database

import (
	"context"
	"fmt"
	"log"
	"money-telegram-bot/internal/models"
	"money-telegram-bot/internal/utils"
	"os"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamoClient *dynamodb.Client
var tableName = os.Getenv("TABLE_NAME")

func InitDB(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to load AWS config: %v", err)
		return err
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	log.Println("[INFO] DynamoDB client initialized successfully")
	return nil
}

func SaveExpense(ctx context.Context, expense *models.Expense) error {
	if dynamoClient == nil {
		return fmt.Errorf("DynamoDB client is not initialized")
	}

	expense.ExpenseID = utils.GenerateExpenseID()

	av, err := attributevalue.MarshalMap(expense)
	if err != nil {
		return err
	}

	log.Println("Table being used:", tableName)

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	return err
}

func GetUserExpenses(ctx context.Context, userID int64) ([]models.Expense, error) {
	if dynamoClient == nil {
		return nil, fmt.Errorf("DynamoDB client not initialized")
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", userID)},
		},
		ScanIndexForward: aws.Bool(true),
	}

	result, err := dynamoClient.Query(ctx, input)
	if err != nil {
		log.Printf("[ERROR] Failed to query expenses: %v", err)
		return nil, err
	}

	var expenses []models.Expense
	err = attributevalue.UnmarshalListOfMaps(result.Items, &expenses)
	if err != nil {
		log.Printf("[ERROR] Failed to unmarshal expenses: %v", err)
		return nil, err
	}

	return expenses, nil
}

func GetTotalExpenses(ctx context.Context, userID int64) (int, error) {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return 0, err
	}
	return len(expenses), nil
}

func DeleteAllExpenses(ctx context.Context, userID int64) error {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return err
	}
	for _, expense := range expenses {
		_, err := dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"user_id":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expense.UserID)},
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

