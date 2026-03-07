package handlers

import (
	"context"
	"log"
	"money-telegram-bot/internal/service"
	"strings"

	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteHandler struct {
	service *service.ExpenseService
}

func NewDeleteHandler(service *service.ExpenseService) *DeleteHandler {
	return &DeleteHandler{service: service}
}

func (h *DeleteHandler) HandleDeleteAll(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf(utils.WarnDeleteAllInvoked, message.From.ID)

	err := h.service.DeleteAllExpenses(context.Background(), message.From.ID)
	if err != nil {
		log.Printf(utils.ErrorDeletingAllExpenses, message.From.ID, err)
		utils.Reply(bot, message.Chat.ID, utils.ErrDeletingExpenses)
		return
	}

	log.Printf(utils.InfoAllExpensesDeleted, message.From.ID)
	utils.Reply(bot, message.Chat.ID, utils.SuccessDeleteAll)
}

func (h *DeleteHandler) HandleDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf(utils.InfoProcessingDelete, message.From.ID)

	args := strings.Fields(message.Text)

	if len(args) < 2 {
		utils.Reply(bot, message.Chat.ID, utils.ErrShouldUseFormatToDelete)
		return
	}

	expenseID := strings.ToUpper(args[1])
	if !strings.HasPrefix(expenseID, "AM") {
		utils.Reply(bot, message.Chat.ID, utils.ErrInvalidDeleteIdFormat)
		return
	}

	err := h.service.DeleteExpense(context.Background(), message.From.ID, expenseID)
	if err != nil {
		log.Printf(utils.ErrorDeletingExpense, message.From.ID, expenseID, err)
		utils.Reply(bot, message.Chat.ID, utils.ErrDeletingExpense)
		return
	}

	log.Printf(utils.InfoExpenseDeleted, message.From.ID, expenseID)
	utils.Reply(bot, message.Chat.ID, utils.SuccessDeleteExpense)
}
