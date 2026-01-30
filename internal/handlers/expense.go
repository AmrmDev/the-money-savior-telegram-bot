package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleExpense(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("handling /gastei command: %q", message.Text)

	parts := strings.Fields(message.Text)

	if len(parts) < 3 {
		reply(bot, message.Chat.ID, "Tente usar o formato: /gastei <valor> <categoria> [metodo]")
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
		reply(bot, message.Chat.ID, "Valor inválido, exemplo: /gastei 21.74 uber pix")
		return
	}

	response := fmt.Sprintf(
		"Anotado:\nValor: R$%.2f\nDescrição: %s\nMétodo: %s",
		amount,
		description,
		method,
	)

	reply(bot, message.Chat.ID, response)
}

func reply(bot *tgbotapi.BotAPI, chatID int64, text string) {
	log.Printf("bot replying to chat ID %d: %q", chatID, text)

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
