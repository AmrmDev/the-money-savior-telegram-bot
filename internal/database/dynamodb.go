package database

import (
	"context"
	"fmt"
	"log"
	"money-telegram-bot/internal/models"
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

func getNextSeqID(ctx context.Context, userID int64) (int, error) {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return 0, err
	}
	maxID := 0
	for _, e := range expenses {
		if e.SeqID > maxID {
			maxID = e.SeqID
		}
	}
	return maxID + 1, nil
}

func SaveExpense(ctx context.Context, expense *models.Expense) error {
	if dynamoClient == nil {
		return fmt.Errorf("DynamoDB client is not initialized")
	}

	nextSeq, err := getNextSeqID(ctx, expense.UserID)
	if err != nil {
		log.Printf("[ERROR] Failed to get next seq_id: %v", err)
		return err
	}
	expense.SeqID = nextSeq
	expense.ExpenseID = fmt.Sprintf("%d#%s", expense.UserID, expense.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))

	av, err := attributevalue.MarshalMap(expense)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal expense: %v", err)
		return err
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		log.Printf("[ERROR] Failed to save expense to DynamoDB: %v", err)
		return err
	}

	log.Printf(
		"[INFO] Expense saved successfully | userID=%d | seqID=%d | amount=R$%.2f | category=%s",
		expense.UserID,
		expense.SeqID,
		expense.Amount,
		expense.Category,
	)

	return nil
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

func GetExpenseBySeqID(ctx context.Context, userID int64, seqID int) (*models.Expense, error) {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range expenses {
		if expenses[i].SeqID == seqID {
			return &expenses[i], nil
		}
	}
	return nil, fmt.Errorf("nenhum gasto encontrado com ID %d", seqID)
}

func GetTotalExpenses(ctx context.Context, userID int64) (int, error) {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return 0, err
	}
	return len(expenses), nil
}

func DeleteExpenseBySeqID(ctx context.Context, userID int64, seqID int) error {
	expenses, err := GetUserExpenses(ctx, userID)
	if err != nil {
		return err
	}
	for _, expense := range expenses {
		if expense.SeqID == seqID {
			_, err := dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]types.AttributeValue{
					"user_id":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expense.UserID)},
					"expense_id": &types.AttributeValueMemberS{Value: expense.ExpenseID},
				},
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete expense seqID=%d: %v", seqID, err)
				return err
			}
			log.Printf("[INFO] Expense deleted | userID=%d | seqID=%d", userID, seqID)
			return nil
		}
	}
	return fmt.Errorf("nenhum gasto encontrado com ID %d", seqID)
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