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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleExpense(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /gastei command")
	log.Printf("[DEBUG] Raw input: %q", message.Text)

	parts := strings.Fields(message.Text)

	if len(parts) < 3 {
		log.Println("[ERROR] Invalid command format. Expected: /gastei <amount> <category> [method]")
		reply(bot, message, "⚠️ Formato incorreto — use: /gastei <valor> <categoria> [método] | Exemplo: /gastei 21.90 uber pix")
		return
	}

	amountStr := parts[1]
	description := parts[2]

	method := "desconhecido"
	if len(parts) >= 4 {
		method = parts[3]
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Printf("[ERROR] Failed to parse amount value: %q | error=%v", amountStr, err)
		reply(bot, message, "Valor inválido, exemplo: /gastei 21.74 uber pix")
		return
	}

	log.Printf(
		"[INFO] Expense parsed successfully | amount=R$%.2f | category=%s | method=%s",
		amount,
		description,
		method,
	)

	user := message.From
	expense := &models.Expense{
		UserID:    user.ID,
		ChatID:    message.Chat.ID,
		Username:  user.UserName,
		Amount:    amount,
		Category:  description,
		Method:    method,
		CreatedAt: time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.SaveExpense(ctx, expense); err != nil {
		log.Printf("[ERROR] Failed to save expense: %v", err)
		reply(bot, message, "❌ Erro ao salvar gasto. Tente novamente.")
		return
	}

	awaitingServerResponse := "⏳ Registrando seu gasto..."
	response := fmt.Sprintf(
		"✅ Gasto registrado com sucesso!\n\n"+
			"💰 Valor: R$%.2f\n"+
			"📝 Descrição: %s\n"+
			"💳 Método: %s",
		amount,
		description,
		method,
	)

	reply(bot, message, awaitingServerResponse)
	time.Sleep(1 * time.Second)
	reply(bot, message, response)
}

func reply(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	user := message.From

	username := user.UserName
	if username == "" {
		username = "sem_username"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("[ERROR] Failed to send message | chatID=%d | userID=%d | error=%v", message.Chat.ID, user.ID, err)
		return
	}

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}