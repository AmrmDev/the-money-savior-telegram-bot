package utils

import (
	"fmt"
	"money-telegram-bot/internal/models"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ErrInvalidFormat      = "⚠️ Formato incorreto — use:\n/gastei <valor> <categoria> [método]\n\nExemplo:\n/gastei 21.90 uber pix"
	ErrInvalidAmount      = "❌ Valor inválido."
	ErrSaveExpense        = "❌ Erro ao salvar gasto."
	ErrExpenseNotFound    = "❌ Nenhum gasto encontrado com esse ID."
	ErrQueryExpenses      = "❌ Erro ao consultar gastos."
	ErrNullExpensesReturn = "📝 Você ainda não registrou nenhum gasto."
)

func SuccessExpenseMessage(
	expenseID string,
	amount float64,
	category string,
	method string,
) string {
	return fmt.Sprintf(
		"✅ *Gasto registrado com sucesso!*\n\n"+
			"🆔 %s\n"+
			"💰 R$ %.2f\n"+
			"📝 %s\n"+
			"💳 %s",
		expenseID,
		amount,
		category,
		method,
	)
}

func ExpenseDetailsMessage(expense *models.Expense) string {
	return fmt.Sprintf(
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
}

func ReplyMarkdown(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func ExpenseListMessage(expenses []models.Expense) string {
	var response strings.Builder

	response.WriteString(
		fmt.Sprintf("📋 *Seus gastos (%d registros):*\n\n", len(expenses)),
	)

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

	return response.String()
}
