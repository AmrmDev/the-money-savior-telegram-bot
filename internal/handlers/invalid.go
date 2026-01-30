package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInvalidCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	errorText := fmt.Sprintf(`✖︎ Comando não reconhecido: %s

O comando "%s" não existe no Money Savior.

Digite /help para ver todos os comandos disponíveis!`, message.Command(), message.Command())

	msg := tgbotapi.NewMessage(message.Chat.ID, errorText)

	log.Printf("bot sending invalid command message to chat ID %d | command: %s", message.Chat.ID, message.Command())
	bot.Send(msg)
	log.Printf("invalid command message sent", message.Chat.ID)
}
