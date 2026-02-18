package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInvalidCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("────────── [HANDLER] invalid command ──────────")

	errorText := fmt.Sprintf(
		`✖︎ Comando não reconhecido: %s
		

O comando "%s" não existe no Money Savior.

Digite /help para ver todos os comandos disponíveis!`,
		message.Command(),
		message.Command(),
	)

	log.Printf("[HANDLER] invalid command received: %s", message.Command())

	msg := tgbotapi.NewMessage(message.Chat.ID, errorText)

	log.Printf("[BOT] sending invalid command message | chatID=%d", message.Chat.ID)
	bot.Send(msg)

	log.Println("────────── [END] invalid command ──────────\n")
}
