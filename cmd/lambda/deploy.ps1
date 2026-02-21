param(
    [string]$Action = "help"
)

$FunctionName = "money-savior-webhook"
$RoleName = "lambda-execution-role"
$AwsRegion = "us-east-1"

function Show-Help {
    Write-Host @"
AWS Lambda Deploy Script for Money Savior
==========================================

Usage: .\deploy.ps1 -Action <action>

Actions:
  help            Show this help message
  setup           Create IAM role and DynamoDB table
  build           Build the Lambda function
  deploy          Build and deploy to Lambda
  configure       Set environment variables
  test            Invoke Lambda function with test data
  logs            View recent Lambda logs
  cleanup         Delete Lambda function and role

Prerequisites:
  - AWS CLI configured with credentials
  - Go 1.23.0 or higher
  - PowerShell 5.0 or higher

Example:
  .\deploy.ps1 -Action setup
  .\deploy.ps1 -Action deploy
"@
}

function Test-AwsCredentials {
    try {
        aws sts get-caller-identity *> $null
        return $true
    } catch {
        Write-Host "ERROR: AWS credentials not configured!" -ForegroundColor Red
        Write-Host "Run: aws configure" -ForegroundColor Yellow
        return $false
    }
}

function Test-Go {
    $goVersion = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Go installed: $goVersion" -ForegroundColor Green
        return $true
    } else {
        Write-Host "ERROR: Go not found!" -ForegroundColor Red
        return $false
    }
}

function Invoke-Setup {
    Write-Host "=== Setting up Lambda Environment ===" -ForegroundColor Cyan
    
    if (-not (Test-AwsCredentials)) { return }
    
    $AccountId = aws sts get-caller-identity --query Account --output text
    $RoleArn = "arn:aws:iam::$($AccountId):role/$RoleName"
    
    Write-Host "[1] Creating IAM Role..." -ForegroundColor Blue
    $trustPolicy = Get-Content "cmd/lambda/trust-policy.json" -Raw
    
    try {
        aws iam create-role `
            --role-name $RoleName `
            --assume-role-policy-document $trustPolicy 2>&1 | Out-Null
        Write-Host "  ✓ Role created: $RoleArn" -ForegroundColor Green
    } catch {
        Write-Host "  ! Role already exists" -ForegroundColor Yellow
    }
    
    Write-Host "[2] Attaching DynamoDB policy..." -ForegroundColor Blue
    aws iam attach-role-policy `
        --role-name $RoleName `
        --policy-arn arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess
    
    Write-Host "[3] Attaching CloudWatch Logs policy..." -ForegroundColor Blue
    aws iam attach-role-policy `
        --role-name $RoleName `
        --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
    
    Write-Host "[4] Checking DynamoDB table..." -ForegroundColor Blue
    $tableExists = aws dynamodb describe-table --table-name expenses --region $AwsRegion 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  ! You need to create DynamoDB table 'expenses'" -ForegroundColor Yellow
        Write-Host "  Run this command:" -ForegroundColor Yellow
        Write-Host "  aws dynamodb create-table --table-name expenses --attribute-definitions AttributeName=user_id,AttributeType=N AttributeName=expense_id,AttributeType=S --key-schema AttributeName=user_id,KeyType=HASH AttributeName=expense_id,KeyType=RANGE --billing-mode PAY_PER_REQUEST" -ForegroundColor Yellow
    } else {
        Write-Host "  ✓ DynamoDB table exists" -ForegroundColor Green
    }
    
    Write-Host "`n✅ Setup completed!" -ForegroundColor Green
}

function Invoke-Build {
    Write-Host "=== Building Lambda Function ===" -ForegroundColor Cyan
    
    if (-not (Test-Go)) { return }
    
    Write-Host "[1] Building for Linux..." -ForegroundColor Blue
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -o "cmd/lambda/bootstrap" "cmd/lambda/main.go"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Binary created: cmd/lambda/bootstrap" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Build failed!" -ForegroundColor Red
        return
    }
    
    Write-Host "[2] Creating ZIP package..." -ForegroundColor Blue
    Compress-Archive -Path "cmd/lambda/bootstrap" -DestinationPath "cmd/lambda/function.zip" -Force
    Write-Host "  ✓ Package created: cmd/lambda/function.zip" -ForegroundColor Green
}

function Invoke-Deploy {
    Write-Host "=== Deploying to Lambda ===" -ForegroundColor Cyan
    
    if (-not (Test-AwsCredentials)) { return }
    if (-not (Test-Go)) { return }
    
    # Build first
    Invoke-Build
    
    $AccountId = aws sts get-caller-identity --query Account --output text
    $RoleArn = "arn:aws:iam::$($AccountId):role/$RoleName"
    
    Write-Host "[3] Checking if function exists..." -ForegroundColor Blue
    $functionExists = aws lambda get-function --function-name $FunctionName --region $AwsRegion 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "[4] Updating function code..." -ForegroundColor Blue
        aws lambda update-function-code `
            --function-name $FunctionName `
            --zip-file "fileb://cmd/lambda/function.zip" `
            --region $AwsRegion | Out-Null
        Write-Host "  ✓ Function updated" -ForegroundColor Green
    } else {
        Write-Host "[4] Creating new function..." -ForegroundColor Blue
        aws lambda create-function `
            --function-name $FunctionName `
            --runtime "provided.al2" `
            --role $RoleArn `
            --handler "bootstrap" `
            --zip-file "fileb://cmd/lambda/function.zip" `
            --timeout 30 `
            --memory-size 256 `
            --region $AwsRegion | Out-Null
        Write-Host "  ✓ Function created" -ForegroundColor Green
    }
    
    Write-Host "`n✅ Deploy completed!" -ForegroundColor Green
}

function Invoke-Configure {
    Write-Host "=== Configuring Lambda Environment ===" -ForegroundColor Cyan
    
    if (-not (Test-AwsCredentials)) { return }
    
    $TelegramToken = Read-Host "Enter TELEGRAM_BOT_TOKEN"
    if ([string]::IsNullOrEmpty($TelegramToken)) {
        Write-Host "ERROR: Token cannot be empty!" -ForegroundColor Red
        return
    }
    
    Write-Host "[1] Updating environment variables..." -ForegroundColor Blue
    aws lambda update-function-configuration `
        --function-name $FunctionName `
        --environment Variables="{TELEGRAM_BOT_TOKEN=$TelegramToken,TABLE_NAME=expenses,AWS_REGION=$AwsRegion}" `
        --region $AwsRegion | Out-Null
    
    Write-Host "  ✓ Environment variables updated" -ForegroundColor Green
    Write-Host "`n✅ Configuration completed!" -ForegroundColor Green
}

function Invoke-Logs {
    Write-Host "=== Lambda Logs ===" -ForegroundColor Cyan
    
    if (-not (Test-AwsCredentials)) { return }
    
    Write-Host "[1] Fetching recent logs..." -ForegroundColor Blue
    aws logs tail "/aws/lambda/$FunctionName" --follow --region $AwsRegion
}

function Invoke-Test {
    Write-Host "=== Testing Lambda Function ===" -ForegroundColor Cyan
    
    if (-not (Test-AwsCredentials)) { return }
    
    $testPayload = @{
        Body = @{
            message = "test"
        } | ConvertTo-Json
    } | ConvertTo-Json
    
    Write-Host "[1] Invoking function..." -ForegroundColor Blue
    aws lambda invoke `
        --function-name $FunctionName `
        --payload $testPayload `
        --log-type Tail `
        --region $AwsRegion `
        response.json | Out-Null
    
    Write-Host "  ✓ Response saved to response.json" -ForegroundColor Green
    Get-Content "response.json" | Out-Host
}

function Invoke-Cleanup {
    Write-Host "=== Cleanup ===" -ForegroundColor Cyan
    Write-Host "This will DELETE the Lambda function and IAM role!" -ForegroundColor Red
    $confirm = Read-Host "Are you sure? (yes/no)"
    
    if ($confirm -ne "yes") { return }
    
    if (-not (Test-AwsCredentials)) { return }
    
    Write-Host "[1] Deleting Lambda function..." -ForegroundColor Blue
    aws lambda delete-function --function-name $FunctionName --region $AwsRegion 2>&1 | Out-Null
    
    Write-Host "[2] Removing IAM role policies..." -ForegroundColor Blue
    aws iam detach-role-policy `
        --role-name $RoleName `
        --policy-arn arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess 2>&1 | Out-Null
    aws iam detach-role-policy `
        --role-name $RoleName `
        --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole 2>&1 | Out-Null
    
    Write-Host "[3] Deleting IAM role..." -ForegroundColor Blue
    aws iam delete-role --role-name $RoleName 2>&1 | Out-Null
    
    Write-Host "✅ Cleanup completed!" -ForegroundColor Green
}

# Main
switch ($Action.ToLower()) {
    "help" { Show-Help }
    "setup" { Invoke-Setup }
    "build" { Invoke-Build }
    "deploy" { Invoke-Deploy; Invoke-Configure }
    "configure" { Invoke-Configure }
    "logs" { Invoke-Logs }
    "test" { Invoke-Test }
    "cleanup" { Invoke-Cleanup }
    default { 
        Write-Host "Unknown action: $Action" -ForegroundColor Red
        Show-Help
    }
}
