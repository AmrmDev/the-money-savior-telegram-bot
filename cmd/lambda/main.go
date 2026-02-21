package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
	
	"money-telegram-bot/internal/bot"
	"money-telegram-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var telegramBot *tgbotapi.BotAPI

func init() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("[FATAL] TELEGRAM_BOT_TOKEN environment variable not configured.")
	}

	var err error
	telegramBot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("[FATAL] Failed to initialize Telegram bot:", err)
	}

	log.Printf("[INFO] Lambda bot authenticated as @%s", telegramBot.Self.UserName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := database.InitDB(ctx); err != nil {
		log.Fatal("[FATAL] Failed to initialize DynamoDB:", err)
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) error {
	var update tgbotapi.Update
	if err := json.Unmarshal([]byte(req.Body), &update); err != nil {
		log.Printf("[ERROR] Failed to parse Telegram update: %v", err)
		return err
	}

	bot.RouteUpdate(telegramBot, update)
	return nil
}

func main() {
	lambda.Start(Handler)
}