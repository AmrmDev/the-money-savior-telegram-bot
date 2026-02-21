#!/bin/bash
set -e

echo "=== Money Savior Lambda Deploy ==="

FUNCTION_NAME="money-savior-webhook"
ROLE_NAME="lambda-execution-role"
AWS_REGION="us-east-1"

echo "[1] Building binary for Linux..."
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go

echo "[2] Creating deployment package..."
rm -f function.zip
zip function.zip bootstrap

echo "[3] Checking AWS credentials..."
aws sts get-caller-identity > /dev/null || { echo "AWS credentials not configured!"; exit 1; }

echo "[4] Getting Account ID..."
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ROLE_ARN="arn:aws:iam::$ACCOUNT_ID:role/$ROLE_NAME"

echo "[5] Checking if Lambda function exists..."
if aws lambda get-function --function-name $FUNCTION_NAME --region $AWS_REGION 2>/dev/null; then
  echo "[6] Updating Lambda function code..."
  aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://function.zip \
    --region $AWS_REGION
else
  echo "[6] Creating new Lambda function..."
  aws lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2 \
    --role $ROLE_ARN \
    --handler bootstrap \
    --zip-file fileb://function.zip \
    --timeout 30 \
    --memory-size 256 \
    --region $AWS_REGION
fi

echo "[7] Configuring environment variables..."
read -p "Enter TELEGRAM_BOT_TOKEN: " TELEGRAM_TOKEN
aws lambda update-function-configuration \
  --function-name $FUNCTION_NAME \
  --environment Variables={TELEGRAM_BOT_TOKEN=$TELEGRAM_TOKEN,TABLE_NAME=expenses,AWS_REGION=$AWS_REGION} \
  --region $AWS_REGION

echo "✅ Deploy completed successfully!"
echo ""
echo "Lambda Function Details:"
echo "  Name: $FUNCTION_NAME"
echo "  Region: $AWS_REGION"
echo "  Role ARN: $ROLE_ARN"
echo ""
echo "Next steps:"
echo "  1. Create API Gateway endpoint and link to this Lambda"
echo "  2. Configure Telegram webhook to point to API Gateway URL"
echo ""