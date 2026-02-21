#!/bin/bash
set -e

FUNCTION_NAME="money-savior-webhook"

# Build para Linux
GOOS=linux GOARCH=amd64 go build -o bootstrap

# Zip usando PowerShell
powershell Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

# Atualiza a função na AWS
aws lambda update-function-code \
  --function-name $FUNCTION_NAME \
  --zip-file fileb://function.zip

echo "✅ Deploy concluído!"