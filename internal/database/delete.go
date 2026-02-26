package database

import (
	"context"
	"fmt"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DeleteExpenseByID(ctx context.Context, userID string, expenseID string) error {
	if dynamoClient == nil {
		return fmt.Errorf("DynamoDB client is not initialized")
	}

	_, err := dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{
				Value: userID,
			},
			"expense_id": &types.AttributeValueMemberS{
				Value: expenseID,
			},
		},
	})

	return err
}