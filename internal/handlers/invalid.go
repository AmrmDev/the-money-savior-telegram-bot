package handlers

import (
	"log"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInvalidCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[WARN] Invalid command received: /%s", message.Command())

	utils.ReplyMarkdown(bot, message.Chat.ID, utils.MsgInvalidCommand(message.Command()))

	user := message.From
	username := user.UserName
	if username == "" {
		username = utils.NoUsernameConst
	}

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
