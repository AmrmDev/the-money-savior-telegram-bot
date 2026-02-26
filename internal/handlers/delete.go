package handlers

import (
	"context"
	"strconv"
	"strings"
	"log"

	"money-telegram-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleDeleteAll(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userID := message.From.ID

	log.Printf("[WARN] /deletartudo invoked | userID=%d", userID)

	err := database.DeleteAllExpenses(context.Background(), userID)
	if err != nil {
		reply(bot, message, "❌ Erro ao limpar seus gastos.")
		return
	}

	reply(bot, message, "🧹 Todos os gastos foram apagados com sucesso.")
}

func HandleDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.Fields(message.Text)

	if len(args) < 2 {
		reply(bot, message, "❌ Use: /deletar AM123456")
		return
	}

	expenseID := strings.ToUpper(args[1])

	if !strings.HasPrefix(expenseID, "AM") {
		reply(bot, message, "❌ ID inválido. Exemplo válido: AM123456")
		return
	}

	userID := strconv.FormatInt(message.From.ID, 10)

	err := database.DeleteExpenseByID(context.Background(), userID, expenseID)
	if err != nil {
		reply(bot, message, "❌ Erro ao deletar o gasto.")
		return
	}

	reply(bot, message, "✅ Gasto deletado com sucesso.")
}