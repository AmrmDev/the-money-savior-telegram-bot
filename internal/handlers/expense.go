package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"money-telegram-bot/internal/database"
	"money-telegram-bot/internal/models"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleExpense(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /gastei command")
	log.Printf("[DEBUG] Raw input: %q", message.Text)

	parts := strings.Fields(message.Text)

	if len(parts) < 3 {
		utils.Reply(bot, message.Chat.ID, "⚠️ Formato incorreto — use:\n/gastei <valor> <categoria> [método]\n\nExemplo:\n/gastei 21.90 uber pix")
		return
	}

	amountStr := parts[1]
	categoryInput := parts[2]

	methodInput := "desconhecido"
	if len(parts) >= 4 {
		methodInput = parts[3]
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		utils.Reply(bot, message.Chat.ID, "❌ Valor inválido.\nExemplo correto:\n/gastei 21.74 uber pix")
		return
	}

	// 🔹 NORMALIZAÇÃO AQUI
	category := utils.FormatTitle(categoryInput)
	method := utils.NormalizeMethod(methodInput)

	log.Printf(
		"[INFO] Expense parsed | amount=R$%.2f | category=%s | method=%s",
		amount,
		category,
		method,
	)

	user := message.From
	expense := &models.Expense{
		UserID:    user.ID,
		ChatID:    message.Chat.ID,
		Username:  user.UserName,
		Amount:    amount,
		Category:  category,
		Method:    method,
		CreatedAt: time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.SaveExpense(ctx, expense); err != nil {
		log.Printf("[ERROR] Failed to save expense: %v", err)
		utils.Reply(bot, message.Chat.ID, "❌ Erro ao salvar gasto. Tente novamente.")
		return
	}

	response := fmt.Sprintf(
		"✅ *Gasto registrado com sucesso!*\n\n"+
			"🆔 %s\n"+
			"💰 R$%.2f\n"+
			"📝 %s\n"+
			"💳 %s",
		expense.ExpenseID,
		amount,
		category,
		method,
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

	if _, err := bot.Send(msg); err != nil {
		log.Printf("[ERROR] Failed to send confirmation message: %v", err)
	}
}