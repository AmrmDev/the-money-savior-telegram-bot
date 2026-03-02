package repository

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var db *dynamodb.Client

func InitDB(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("erro ao carregar config AWS: %v", err)
	}

	db = dynamodb.NewFromConfig(cfg)
	return db
}
