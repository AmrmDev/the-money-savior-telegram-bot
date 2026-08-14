package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"money-telegram-bot/internal/controller"
	"money-telegram-bot/internal/handlers"
	"money-telegram-bot/internal/repository"
	"money-telegram-bot/internal/service"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var telegramBot *tgbotapi.BotAPI
var botHandlers *controller.BotHandlers // ponteiro para a struct, não o tipo

func init() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal(utils.FatalEnvironmentVariableNotSet)
	}

	var err error
	telegramBot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(utils.FatalTelegramInitialized, err)
	}
	log.Printf(utils.InfoLambdaAuthenticated, telegramBot.Self.UserName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := repository.InitDB(ctx) // captura o client retornado
	if err != nil {
		log.Fatal(utils.FatalInitializeDynamoDB, err)
	}

	expenseRepo := repository.NewDynamoExpenseRepository(dbClient)
	expenseService := service.NewExpenseService(expenseRepo)

	botHandlers = &controller.BotHandlers{
		Expense: handlers.NewExpenseHandler(expenseService),
		Query:   handlers.NewQueryHandler(expenseService),
		Delete:  handlers.NewDeleteHandler(expenseService),
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) error {
	var update tgbotapi.Update
	if err := json.Unmarshal([]byte(req.Body), &update); err != nil {
		log.Printf(utils.ErrorParseUpdateTelegram, err)
		return err
	}

	controller.RouteUpdate(telegramBot, update, botHandlers)
	return nil
}

func main() {
	lambda.Start(Handler)
}
