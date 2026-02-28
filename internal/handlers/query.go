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
	expenseService *services.ExpenseService
}

func NewQueryHandler(expenseService *services.ExpenseService) *QueryHandler {
	return &QueryHandler{expenseService: expenseService}
}

func (h *QueryHandler) Handle(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[INFO] Processing /consulta | userID=%d", message.From.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := strings.Fields(message.CommandArguments())
	userID := fmt.Sprint(message.From.ID)

	// 🔹 Consulta por ID
	if len(args) == 1 {
		expenseID := strings.ToUpper(args[0])

		expense, err := h.expenseService.GetByID(ctx, userID, expenseID)
		if err != nil {
			utils.Reply(bot, message.Chat.ID, utils.ErrExpenseNotFound)
			return
		}

		text := fmt.Sprintf(
			"📄 *Detalhes do gasto*\n\n"+
				"🆔 ID: *%s*\n"+
				"💰 Valor: *R$ %.2f*\n"+
				"📝 Categoria: *%s*\n"+
				"💳 Método: *%s*\n"+
				"🕐 Data: *%s*",
			expense.ExpenseID,
			expense.Amount,
			expense.Category,
			expense.Method,
			expense.CreatedAt.Format("02/01/2006 15:04"),
		)

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		return
	}

	expenses, err := h.expenseService.ListByUser(ctx, userID)
	if err != nil {
		utils.Reply(bot, message.Chat.ID, utils.ErrQueryExpenses)
		return
	}

	if len(expenses) == 0 {
		utils.Reply(bot, message.Chat.ID, utils.ErrNullExpensesReturn)
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📋 *Seus gastos (%d registros):*\n\n", len(expenses)))

	for _, e := range expenses {
		response.WriteString(fmt.Sprintf(
			"🆔 *%s* | 💰 R$ %.2f | 📝 %s | 💳 %s\n",
			e.ExpenseID,
			e.Amount,
			e.Category,
			e.Method,
		))
	}

	response.WriteString("\n💡 Use `/consulta AMxxxxxx` para ver detalhes.")

	msg := tgbotapi.NewMessage(message.Chat.ID, response.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
