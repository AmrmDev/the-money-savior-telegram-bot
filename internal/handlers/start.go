package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		"working",
	)

	log.Printf("bot sending message: %q", msg.Text)
	bot.Send(msg)
	log.Printf("message sent to chat ID %d | content: %q", message.Chat.ID, msg.Text)
	fmt.Println("")
}