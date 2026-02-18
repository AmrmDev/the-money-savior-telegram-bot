package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	welcomeText := ` Bem-vindo ao Money Savior!

Seu assistente pessoal para rastreamento de despesas.


Comandos Disponíveis:

/gastei - Registre suas despesas
/help - Veja todos os comandos disponíveis

Digite /help para mais informações!`

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)

	log.Printf("bot sending start message to chat ID %d", message.Chat.ID)
	bot.Send(msg)
	log.Printf("start message sent to chat ID %d", message.Chat.ID)
	fmt.Println("")
}