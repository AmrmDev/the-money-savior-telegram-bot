package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleExpense(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("────────── [HANDLER] /gastei ──────────")
	log.Printf("[HANDLER] raw message: %q", message.Text)

	parts := strings.Fields(message.Text)

	if len(parts) < 3 {
		log.Println("[ERROR] invalid /gastei format")
		reply(bot, message.Chat.ID, "Tente usar o formato: /gastei <valor> <categoria> [metodo]")
		log.Println("────────── [END] /gastei ──────────\n")
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
		log.Printf("[ERROR] invalid amount: %q", amountStr)
		reply(bot, message.Chat.ID, "Valor inválido, exemplo: /gastei 21.74 uber pix")
		log.Println("────────── [END] /gastei ──────────\n")
		return
	}

	log.Printf("[HANDLER] parsed expense | amount=%.2f | desc=%s | method=%s",
		amount, description, method,
	)

	response := fmt.Sprintf(
		"✅ Gasto registrado com sucesso!\n\n"+
			"💰 Valor: R$%.2f\n"+
			"📝 Descrição: %s\n"+
			"💳 Método: %s",
		amount,
		description,
		method,
	)

	reply(bot, message.Chat.ID, response)

	log.Println("────────── [END] /gastei ──────────\n")
}

func reply(bot *tgbotapi.BotAPI, chatID int64, text string) {
	log.Printf("[BOT] replying to chat ID %d | content=%q", chatID, text)

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}