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
	log.Printf("[WARN] /deletartudo invoked | userID=%d", message.From.ID)

	err := h.service.DeleteAllExpenses(context.Background(), message.From.ID)
	if err != nil {
		utils.Reply(bot, message.Chat.ID, utils.ErrDeletingExpenses)
		return
	}

	utils.Reply(bot, message.Chat.ID, utils.SuccessDeleteAll)
}

func (h *DeleteHandler) HandleDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
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
		utils.Reply(bot, message.Chat.ID, utils.ErrDeletingExpense)
		return
	}

	utils.Reply(bot, message.Chat.ID, utils.SuccessDeleteExpense)
}
