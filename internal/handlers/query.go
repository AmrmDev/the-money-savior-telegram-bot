package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"money-telegram-bot/internal/database"
	"money-telegram-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleQuery handles /consulta — lists all expenses or shows a specific one with navigation.
func HandleQuery(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[INFO] Processing /consulta command | chatID=%d | userID=%d", message.Chat.ID, message.From.ID)

	args := strings.Fields(message.CommandArguments())

	if len(args) == 1 {
		// /consulta <id> — show single expense with navigation
		seqID, err := strconv.Atoi(args[0])
		if err != nil || seqID < 1 {
			reply(bot, message, "❌ ID inválido. Use um número inteiro maior que zero.\nExemplo: /consulta 3")
			return
		}
		sendExpenseView(bot, message.Chat.ID, message.From.ID, seqID, 0)
		return
	}

	// /consulta — list all expenses
	expenses, err := database.GetUserExpenses(context.Background(), message.From.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to query expenses for userID=%d: %v", message.From.ID, err)
		reply(bot, message, "❌ Ocorreu um erro ao consultar seus gastos. Tente novamente mais tarde.")
		return
	}

	if len(expenses) == 0 {
		reply(bot, message, "📝 Você ainda não registrou nenhum gasto.")
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📋 *Seus gastos (%d registros):*\n\n", len(expenses)))

	for _, expense := range expenses {
		response.WriteString(fmt.Sprintf(
			"🆔 *#%d* | 💰 R$ %.2f | 📝 %s | 💳 %s\n",
			expense.SeqID,
			expense.Amount,
			expense.Category,
			expense.Method,
		))
	}

	response.WriteString("\n💡 Use /consulta <ID> para ver detalhes de um gasto específico.\nExemplo: /consulta 2")

	msg := tgbotapi.NewMessage(message.Chat.ID, response.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// sendExpenseView sends a single expense card with prev/next navigation buttons.
// replyToMessageID is optional (0 = new message).
func sendExpenseView(bot *tgbotapi.BotAPI, chatID int64, userID int64, seqID int, editMessageID int) {
	expense, err := database.GetExpenseBySeqID(context.Background(), userID, seqID)
	if err != nil {
		text := fmt.Sprintf("❌ Nenhum gasto encontrado com o ID *%d*.\nUse /consulta para ver a lista completa.", seqID)
		if editMessageID != 0 {
			edit := tgbotapi.NewEditMessageText(chatID, editMessageID, text)
			edit.ParseMode = "Markdown"
			bot.Send(edit)
		} else {
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
		return
	}

	total, _ := database.GetTotalExpenses(context.Background(), userID)
	text := buildExpenseCard(expense, seqID, total)
	keyboard := buildNavKeyboard(userID, seqID, total)

	if editMessageID != 0 {
		edit := tgbotapi.NewEditMessageText(chatID, editMessageID, text)
		edit.ParseMode = "Markdown"
		edit.ReplyMarkup = &keyboard
		bot.Send(edit)
	} else {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}

func buildExpenseCard(expense *models.Expense, seqID int, total int) string {
	return fmt.Sprintf(
		"📄 *Gasto %d de %d*\n\n"+
			"🆔 ID: *%d*\n"+
			"💰 Valor: *R$ %.2f*\n"+
			"📝 Categoria: *%s*\n"+
			"💳 Método: *%s*\n"+
			"🕐 Data: *%s*",
		seqID, total,
		expense.SeqID,
		expense.Amount,
		expense.Category,
		expense.Method,
		expense.CreatedAt.Format("02/01/2006 15:04"),
	)
}

func buildNavKeyboard(userID int64, seqID int, total int) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton

	if seqID > 1 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(
			"⬅️ Anterior",
			fmt.Sprintf("qnav:%d:%d", userID, seqID-1),
		))
	} else {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("⬅️ Anterior", "qnav_disabled"))
	}

	row = append(row, tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("📋 %d/%d", seqID, total),
		"qnav_info",
	))

	if seqID < total {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(
			"Próximo ➡️",
			fmt.Sprintf("qnav:%d:%d", userID, seqID+1),
		))
	} else {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Próximo ➡️", "qnav_disabled"))
	}

	deleteRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🗑️ Deletar #%d", seqID),
			fmt.Sprintf("confirm_delete:%d", seqID),
		),
		tgbotapi.NewInlineKeyboardButtonData("📋 Ver todos", "qnav_list"),
	)

	return tgbotapi.NewInlineKeyboardMarkup(row, deleteRow)
}

// HandleQueryCallback handles inline navigation callbacks for the expense viewer.
func HandleQueryCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := callback.From.ID

	switch {
	case callback.Data == "qnav_disabled" || callback.Data == "qnav_info":
		bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return

	case callback.Data == "qnav_list":
		bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		expenses, err := database.GetUserExpenses(context.Background(), userID)
		if err != nil || len(expenses) == 0 {
			edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "📝 Nenhum gasto registrado.")
			bot.Send(edit)
			return
		}
		var response strings.Builder
		response.WriteString(fmt.Sprintf("📋 *Seus gastos (%d registros):*\n\n", len(expenses)))
		for _, expense := range expenses {
			response.WriteString(fmt.Sprintf(
				"🆔 *#%d* | 💰 R$ %.2f | 📝 %s | 💳 %s\n",
				expense.SeqID,
				expense.Amount,
				expense.Category,
				expense.Method,
			))
		}
		response.WriteString("\n💡 Use /consulta <ID> para ver detalhes.\nExemplo: /consulta 2")
		edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, response.String())
		edit.ParseMode = "Markdown"
		bot.Send(edit)
		return

	default:
		// data format: "qnav:<userID>:<seqID>"
		var targetUserID int64
		var seqID int
		fmt.Sscanf(callback.Data, "qnav:%d:%d", &targetUserID, &seqID)

		bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		sendExpenseView(bot, chatID, userID, seqID, callback.Message.MessageID)
	}
}
