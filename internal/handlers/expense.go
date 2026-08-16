package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"money-telegram-bot/internal/service"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

func (h *ExpenseHandler) Handle(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.From == nil {
		log.Printf(utils.ErrorNoUserIdentified, message.Chat.ID)
		utils.Reply(bot, message.Chat.ID, utils.ErrNoUserIdentifiedMessage)
		return
	}
	log.Printf(utils.InfoProcessingExpense, message.From.ID)

	parts := strings.Fields(message.Text)

	if len(parts) < 3 {
		log.Printf(utils.ErrorInvalidExpenseFormat, message.From.ID)
		utils.Reply(bot, message.Chat.ID, utils.ErrInvalidFormat)
		return
	}

	amount, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf(utils.ErrorInvalidExpenseAmount, message.From.ID, parts[1])
		utils.Reply(bot, message.Chat.ID, utils.ErrInvalidAmount)
		return
	}

	if amount <= 0 {
		log.Printf(utils.ErrorInvalidExpenseAmount, message.From.ID, parts[1])
		utils.Reply(bot, message.Chat.ID, utils.ErrInvalidAmount)
		return
	}

	categoryInput := parts[2]
	methodInput := ""
	if len(parts) >= 4 {
		methodInput = parts[3]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf(utils.InfoPersistingDataIntoDynamodb, message.From.ID)
	expense, err := h.expenseService.CreateExpense(
		ctx,
		message.From.ID,
		message.Chat.ID,
		message.From.UserName,
		amount,
		categoryInput,
		methodInput,
	)

	if err != nil {
		log.Printf(utils.ErrorSavingExpense, message.From.ID, err)
		utils.Reply(bot, message.Chat.ID, utils.ErrSaveExpense)
		return
	}

	log.Printf(utils.InfoExpenseCreated, message.From.ID, expense.ExpenseID, expense.Amount)
	utils.Reply(
		bot,
		message.Chat.ID,
		utils.SuccessExpenseMessage(
			expense.ExpenseID,
			expense.Amount,
			expense.Category,
			expense.Method,
		),
	)
}
