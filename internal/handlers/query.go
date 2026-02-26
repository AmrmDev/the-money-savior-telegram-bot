package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"money-telegram-bot/internal/database"
	"money-telegram-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleQuery(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[INFO] Processing /consulta | userID=%d", message.From.ID)

	args := strings.Fields(message.CommandArguments())

	// 🔍 /consulta AMxxxxxx
	if len(args) == 1 {
		expenseID := strings.ToUpper(args[0])

		if !strings.HasPrefix(expenseID, "AM") {
			reply(bot, message, "❌ ID inválido. Exemplo: /consulta AM123456")
			return
		}

		expenses, err := database.GetUserExpenses(context.Background(), message.From.ID)
		if err != nil {
			reply(bot, message, "❌ Erro ao consultar gastos.")
			return
		}

		var expense *models.Expense
		for i := range expenses {
			if expenses[i].ExpenseID == expenseID {
				expense = &expenses[i]
				break
			}
		}

		if expense == nil {
			reply(bot, message, "❌ Nenhum gasto encontrado com esse ID.")
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

	// 📋 /consulta
	expenses, err := database.GetUserExpenses(context.Background(), message.From.ID)
	if err != nil {
		reply(bot, message, "❌ Erro ao consultar gastos.")
		return
	}

	if len(expenses) == 0 {
		reply(bot, message, "📝 Você ainda não registrou nenhum gasto.")
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