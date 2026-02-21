package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"money-telegram-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleDelete handles /deletar <id> — shows inline confirmation before deleting.
func HandleDelete(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[INFO] Processing /deletar command | chatID=%d | userID=%d", message.Chat.ID, message.From.ID)
	args := strings.Fields(message.CommandArguments())
	if len(args) != 1 {
		reply(bot, message, "❌ Uso incorreto. Use: /deletar <ID do gasto>\nExemplo: /deletar 3\n\nUse /consulta para ver os IDs dos seus gastos.")
		return
	}

	seqID, err := strconv.Atoi(args[0])
	if err != nil || seqID < 1 {
		reply(bot, message, "❌ ID inválido. Use um número inteiro maior que zero.\nExemplo: /deletar 3")
		return
	}

	expense, err := database.GetExpenseBySeqID(context.Background(), message.From.ID, seqID)
	if err != nil {
		reply(bot, message, fmt.Sprintf("❌ Nenhum gasto encontrado com o ID %d.\nUse /consulta para ver os IDs disponíveis.", seqID))
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✅ Sim, deletar",
				fmt.Sprintf("confirm_delete:%d", seqID),
			),
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancelar", "cancel_delete"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
		"⚠️ Tem certeza que deseja deletar este gasto?\n\n"+
			"🆔 ID: %d\n"+
			"💰 Valor: R$ %.2f\n"+
			"📝 Categoria: %s\n"+
			"💳 Método: %s",
		expense.SeqID,
		expense.Amount,
		expense.Category,
		expense.Method,
	))
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// HandleConfirmDeleteCallback handles inline button confirmation for single delete.
func HandleConfirmDeleteCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := callback.From.ID

	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	if callback.Data == "cancel_delete" {
		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "❌ Operação cancelada.")
		edit.ReplyMarkup = nil
		bot.Send(edit)
		return
	}

	// data format: "confirm_delete:<seqID>"
	seqID := 0
	fmt.Sscanf(callback.Data, "confirm_delete:%d", &seqID)

	err := database.DeleteExpenseBySeqID(context.Background(), userID, seqID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete expense for userID=%d with seqID=%d: %v", userID, seqID, err)
		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID,
			fmt.Sprintf("❌ Nenhum gasto encontrado com o ID %d.", seqID))
		bot.Send(edit)
		return
	}

	edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID,
		fmt.Sprintf("✅ Gasto #%d deletado com sucesso!", seqID))
	edit.ReplyMarkup = nil
	bot.Send(edit)
	log.Printf("[INFO] Expense deleted | userID=%d | seqID=%d", userID, seqID)
}

// HandleDeleteAll handles /deletartudo — shows inline confirmation before deleting everything.
func HandleDeleteAll(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[INFO] Processing /deletartudo command | chatID=%d | userID=%d", message.Chat.ID, message.From.ID)

	total, err := database.GetTotalExpenses(context.Background(), message.From.ID)
	if err != nil {
		reply(bot, message, "❌ Ocorreu um erro ao consultar seus gastos. Tente novamente mais tarde.")
		return
	}

	if total == 0 {
		reply(bot, message, "📝 Você não possui nenhum gasto registrado.")
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🗑️ Sim, deletar todos (%d)", total),
				"confirm_delete_all",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancelar", "cancel_delete_all"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("⚠️ *Atenção!* Você está prestes a deletar *todos os %d gastos* registrados.\n\nEssa ação é irreversível. Deseja continuar?", total),
	)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// HandleDeleteAllCallback handles inline button confirmation for delete-all.
func HandleDeleteAllCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := callback.From.ID

	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	if callback.Data == "cancel_delete_all" {
		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "❌ Operação cancelada.")
		edit.ReplyMarkup = nil
		bot.Send(edit)
		return
	}

	err := database.DeleteAllExpenses(context.Background(), userID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete all expenses for userID=%d: %v", userID, err)
		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID,
			"❌ Ocorreu um erro ao deletar os gastos. Tente novamente mais tarde.")
		bot.Send(edit)
		return
	}

	edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "✅ Todos os gastos foram deletados com sucesso!")
	edit.ReplyMarkup = nil
	bot.Send(edit)
	log.Printf("[INFO] All expenses deleted | userID=%d", userID)
}