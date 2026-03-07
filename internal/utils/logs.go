package utils

const (
	FatalEnvironmentVariableNotSet = "[FATAL] TELEGRAM_BOT_TOKEN environment variable not configured."
	FatalTelegramInitialized       = "[FATAL] Failed to initialize Telegram bot: "
	FatalInitializeDynamoDB        = "[FATAL] Failed to initialize DynamoDB:"
	FatalBotInitializationFailed   = "[FATAL] Bot initialization failed: "

	InfoStartingBot         = "[INFO] Starting Money Savior Telegram Bot..."
	InfoBotAuthenticated    = "[INFO] Bot authenticated as @%s"
	InfoLambdaAuthenticated = "[INFO] Lambda bot authenticated as @%s"
	InfoMessageReceived     = "[INFO] Message received: %q"
	InfoCommandReceived     = "[INFO] Command received: /%s"
	InfoCommandResponseSent = "[INFO] Response sent | chatID=%d | userID=%d | userName=%s | command=/%s | status=success"

	InfoProcessingStart   = "[INFO] Processing /start command"
	InfoProcessingHelp    = "[INFO] Processing /help command"
	InfoProcessingExpense = "[INFO] Processing /gastei | userID=%d"
	InfoProcessingQuery   = "[INFO] Processing /consulta | userID=%d"
	InfoProcessingDelete  = "[INFO] Processing /deletar | userID=%d"

	InfoPersistingDataIntoDynamodb = "[INFO] Persisting data into DynamoDB | userID=%d"
	InfoExpenseCreated             = "[INFO] Expense created | userID=%d | expenseID=%s | amount=%.2f"
	InfoQueryResultSent            = "[INFO] Query result sent | userID=%s | count=%d"
	InfoExpenseDeleted             = "[INFO] Expense deleted | userID=%d | expenseID=%s"
	InfoAllExpensesDeleted         = "[INFO] All expenses deleted | userID=%d"

	WarnUnknownCommand         = "[WARN] Unknown command received: /%s"
	WarnInvalidCommandReceived = "[WARN] Invalid command received: /%s"
	WarnDeleteAllInvoked       = "[WARN] /deletartudo invoked | userID=%d"

	ErrorParseUpdateTelegram  = "[ERROR] Failed to parse Telegram update: %v"
	ErrorInvalidExpenseFormat = "[ERROR] Invalid expense format | userID=%d"
	ErrorInvalidExpenseAmount = "[ERROR] Invalid expense amount | userID=%d | value=%s"
	ErrorSavingExpense        = "[ERROR] Failed to save expense | userID=%d | err=%v"
	ErrorExpenseNotFound      = "[ERROR] Expense not found | userID=%s | expenseID=%s"
	ErrorQueryingExpenses     = "[ERROR] Failed to query expenses | userID=%s | err=%v"
	ErrorDeletingExpense      = "[ERROR] Failed to delete expense | userID=%d | expenseID=%s | err=%v"
	ErrorDeletingAllExpenses  = "[ERROR] Failed to delete all expenses | userID=%d | err=%v"

	DebugUpdateWithNoMessage = "[DEBUG] Update received with no message. Skipping..."
)
