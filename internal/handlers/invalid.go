package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInvalidCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[WARN] Invalid command received: /%s", message.Command())

	errorText := fmt.Sprintf(
		`❌ *Comando não reconhecido*

O comando *%s* não existe no *Money Savior* 😕  

📋 Para ver a lista de comandos disponíveis, digite */help*.

💡 Dica: confira se o comando foi digitado corretamente.`,
		message.Command(),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, errorText)

	user := message.From

	username := user.UserName
	if username == "" {
		username = "sem_username"
	}

	lastName := user.LastName
	if lastName == "" {
		lastName = "-"
	}

	bot.Send(msg)

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
