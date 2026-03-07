package utils

const (
	FatalEnvironmentVariableNotSet = "[FATAL] TELEGRAM_BOT_TOKEN environment variable not configured."
	FatalTelegramInitialized       = "[FATAL] Failed to initialize Telegram bot: "
	FatalInitializeDynamoDB        = "[FATAL] Failed to initialize DynamoDB:"

	InfoLambdaAuthenticated = "[INFO] Lambda bot authenticated as @%s"

	ErrorParseUpdateTelegram = "[ERROR] Failed to parse Telegram update: %v"
)
