package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var db *dynamodb.Client

func InitDB(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar config AWS: %w", err)
	}

	db = dynamodb.NewFromConfig(cfg)
	return db, nil
}
