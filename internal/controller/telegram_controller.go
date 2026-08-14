package controller

import (
	"context"
	"log"

	"money-telegram-bot/internal/handlers"
	"money-telegram-bot/internal/repository"
	"money-telegram-bot/internal/service"
	"money-telegram-bot/internal/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandlers struct {
	Expense *handlers.ExpenseHandler
	Query   *handlers.QueryHandler
	Delete  *handlers.DeleteHandler
}

func RouteUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	h *BotHandlers,
) {
	msg := update.Message
	if msg == nil {
		msg = update.EditedMessage
	}

	if msg == nil {
		log.Println(utils.DebugUpdateWithNoMessage)
		return
	}

	log.Printf(utils.InfoMessageReceived, msg.Text)

	if msg.IsCommand() {
		command := msg.Command()
		log.Printf(utils.InfoCommandReceived, command)

		switch command {
		case "start":
			handlers.HandleStart(bot, msg)

		case "help":
			handlers.HandleHelp(bot, msg)

		case "gastei":
			h.Expense.Handle(bot, msg)

		case "consulta":
			h.Query.Handle(bot, msg)

		case "deletar":
			h.Delete.HandleDelete(bot, msg)

		case "deletartudo":
			h.Delete.HandleDeleteAll(bot, msg)

		default:
			log.Printf(utils.WarnUnknownCommand, command)
			handlers.HandleInvalidCommand(bot, msg)
		}
	}
}

func Start(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf(utils.InfoBotAuthenticated, bot.Self.UserName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	expenseRepo := repository.NewDynamoExpenseRepository(dynamoClient)

	expenseService := service.NewExpenseService(expenseRepo)

	botHandlers := &BotHandlers{
		Expense: handlers.NewExpenseHandler(expenseService),
		Query:   handlers.NewQueryHandler(expenseService),
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		RouteUpdate(bot, update, botHandlers)
	}

	return nil
}
