# Guia de Deploy para AWS Lambda

## Visão Geral

Este projeto pode ser executado de duas formas:

1. **Bot Server Local** - `cmd/bot/main.go` (polling)
2. **AWS Lambda** - `cmd/lambda/main.go` (webhook via API Gateway)

## Arquitetura Lambda + API Gateway

```
Telegram Server
     |
     | (HTTP POST com update)
     v
API Gateway
     |
     | (invoca)
     v
AWS Lambda (money-savior-webhook)
     |
     | (valida/processa)
     v
DynamoDB (armazena gastos)
```

## Pré-requisitos

- AWS Account com acesso a Lambda e DynamoDB
- AWS CLI v2 instalada (`aws --version`)
- Go 1.23.0+ 
- Credenciais AWS configuradas (`aws configure`)
- Telegram Bot Token

## Passo 1: Configurar IAM Role

### Opção A: Usar PowerShell (Windows)

```powershell
cd cmd/lambda
.\deploy.ps1 -Action setup
```

### Opção B: Manual (Linux/macOS)

```bash
# Crie a role
aws iam create-role \
  --role-name lambda-execution-role \
  --assume-role-policy-document file://trust-policy.json

# Adicione permissão DynamoDB
aws iam attach-role-policy \
  --role-name lambda-execution-role \
  --policy-arn arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess

# Adicione permissão CloudWatch Logs
aws iam attach-role-policy \
  --role-name lambda-execution-role \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

## Passo 2: Criar Tabela DynamoDB

```bash
aws dynamodb create-table \
  --table-name expenses \
  --attribute-definitions \
    AttributeName=user_id,AttributeType=N \
    AttributeName=expense_id,AttributeType=S \
  --key-schema \
    AttributeName=user_id,KeyType=HASH \
    AttributeName=expense_id,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

## Passo 3: Build e Deploy

### Opção A: Usar PowerShell (Windows)

```powershell
cd cmd/lambda
.\deploy.ps1 -Action deploy
# Será solicitado o TELEGRAM_BOT_TOKEN
```

### Opção B: Usar Shell Script (Linux/macOS)

```bash
cd cmd/lambda
chmod +x deploy.sh
./deploy.sh
```

## Passo 4: Configurar API Gateway

### Via AWS Console:

1. Acesse **API Gateway** > **Create API** > **REST API**
2. Nome: `money-savior-api`
3. Criar
4. **Resources** > Criar novo recurso `/webhook`
5. **POST** method
6. **Lambda Function**: `money-savior-webhook`
7. **Deploy API** para stage `prod`
8. Copiar o Invoke URL (exemplo: `https://xxxxx.execute-api.us-east-1.amazonaws.com/prod/webhook`)

### Via AWS CLI:

```bash
# Obter Account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Criar API
API_ID=$(aws apigateway create-rest-api \
  --name money-savior-api \
  --query 'id' --output text)

# Obter root resource ID
ROOT_RESOURCE_ID=$(aws apigateway get-resources \
  --rest-api-id $API_ID \
  --query 'items[0].id' --output text)

# Criar /webhook resource
RESOURCE_ID=$(aws apigateway create-resource \
  --rest-api-id $API_ID \
  --parent-id $ROOT_RESOURCE_ID \
  --path-part webhook \
  --query 'id' --output text)

# Criar POST method
aws apigateway put-method \
  --rest-api-id $API_ID \
  --resource-id $RESOURCE_ID \
  --http-method POST \
  --authorization-type NONE

# Integrar com Lambda
aws apigateway put-integration \
  --rest-api-id $API_ID \
  --resource-id $RESOURCE_ID \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:$ACCOUNT_ID:function:money-savior-webhook/invocations"

# Deploy
aws apigateway create-deployment \
  --rest-api-id $API_ID \
  --stage-name prod

# Obter URL
echo "https://$API_ID.execute-api.us-east-1.amazonaws.com/prod/webhook"
```

## Passo 5: Configurar Telegram Webhook

```bash
WEBHOOK_URL="https://xxxxx.execute-api.us-east-1.amazonaws.com/prod/webhook"
TELEGRAM_TOKEN="seu_bot_token"

curl -X POST \
  https://api.telegram.org/bot$TELEGRAM_TOKEN/setWebhook \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$WEBHOOK_URL\"}"
```

## Verificar Status

### Via PowerShell:

```powershell
# Ver logs
.\deploy.ps1 -Action logs

# Testar função
.\deploy.ps1 -Action test
```

### Via AWS CLI:

```bash
# Ver informações da função
aws lambda get-function --function-name money-savior-webhook

# Ver últimas invocações
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=money-savior-webhook \
  --start-time 2024-01-01T00:00:00Z \
  --end-time 2024-01-02T00:00:00Z \
  --period 3600 \
  --statistics Sum

# Ver logs
aws logs tail /aws/lambda/money-savior-webhook --follow
```

## Troubleshooting

### Erro: "Role not found"

```bash
# Aguarde 10 segundos para a role ser criada
sleep 10
.\deploy.ps1 -Action deploy
```

### Erro: "DynamoDB table not found"

```bash
# Verificar tabela
aws dynamodb describe-table --table-name expenses

# Se não existir, criar:
aws dynamodb create-table --table-name expenses ...
```

### Erro: "Lambda timeout"

Aumentar timeout na Lambda:

```bash
aws lambda update-function-configuration \
  --function-name money-savior-webhook \
  --timeout 60
```

### Bot não está respondendo

1. Verificar se webhook está configurado:
```bash
curl https://api.telegram.org/bot$TOKEN/getWebhookInfo
```

2. Ver logs da Lambda:
```bash
aws logs tail /aws/lambda/money-savior-webhook --follow
```

3. Verificar se DynamoDB tem permissões:
```bash
aws iam get-role-policy \
  --role-name lambda-execution-role \
  --policy-name inline-policy-name
```

## Custos Estimados

- **Lambda**: ~$0.20/mês (1M requests grátis ao mês)
- **DynamoDB**: ~$1/mês (On-Demand, para uso baixo)
- **API Gateway**: ~$3.50/1M requests (100k grátis ao mês)
- **CloudWatch Logs**: ~$0.50/mês

**Total**: Aproximadamente **$5-10/mês** para uso moderado

## Limpeza

Se quiser remover tudo:

```powershell
.\deploy.ps1 -Action cleanup
```

Ou manualmente:

```bash
aws lambda delete-function --function-name money-savior-webhook
aws iam delete-role --role-name lambda-execution-role
aws apigateway delete-rest-api --rest-api-id $API_ID
aws dynamodb delete-table --table-name expenses
```

## Próximos Passos

1. Configurar CI/CD com GitHub Actions para deploy automático
2. Adicionar X-Telegram-Bot-Api-Secret-Token validation
3. Implementar DLQ para failed messages
4. Adicionar CloudWatch alarms para erros

---

**Questões?** Veja [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
