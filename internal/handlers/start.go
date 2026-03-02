package handlers

import (
	"log"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /start command")

	utils.ReplyMarkdown(bot, message.Chat.ID, utils.MsgStart)

	user := message.From
	username := user.UserName
	if username == "" {
		username = utils.NoUsernameConst
	}

	lastName := user.LastName
	if lastName == "" {
		lastName = "-"
	}

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
