package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"money-telegram-bot/internal/service"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QueryHandler struct {
	expenseService *service.ExpenseService
}

func NewQueryHandler(expenseService *service.ExpenseService) *QueryHandler {
	return &QueryHandler{expenseService: expenseService}
}

func (h *QueryHandler) Handle(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.From == nil {
		log.Printf(utils.ErrorNoUserIdentified, message.Chat.ID)
		utils.Reply(bot, message.Chat.ID, utils.ErrNoUserIdentifiedMessage)
		return
	}
	
	log.Printf(utils.InfoProcessingQuery, message.From.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := strings.Fields(message.CommandArguments())
	userID := fmt.Sprint(message.From.ID)

	if len(args) == 1 {
		expense, err := h.expenseService.GetByID(ctx, userID, args[0])
		if err != nil {
			log.Printf(utils.ErrorExpenseNotFound, userID, args[0])
			utils.Reply(bot, message.Chat.ID, utils.ErrExpenseNotFound)
			return
		}

		log.Printf(utils.InfoQueryResultSent, userID, 1)
		utils.ReplyMarkdown(
			bot,
			message.Chat.ID,
			utils.ExpenseDetailsMessage(expense),
		)
		return
	}

	expenses, err := h.expenseService.ListByUser(ctx, userID)
	if err != nil {
		log.Printf(utils.ErrorQueryingExpenses, userID, err)
		utils.Reply(bot, message.Chat.ID, utils.ErrQueryExpenses)
		return
	}

	if len(expenses) == 0 {
		utils.Reply(bot, message.Chat.ID, utils.ErrNullExpensesReturn)
		return
	}

	log.Printf(utils.InfoQueryResultSent, userID, len(expenses))
	utils.ReplyMarkdown(
		bot,
		message.Chat.ID,
		utils.ExpenseListMessage(expenses),
	)
}
