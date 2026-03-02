package handlers

import (
	"log"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /help command")

	utils.ReplyMarkdown(bot, message.Chat.ID, utils.MsgHelp)

	user := message.From
	username := user.UserName
	if username == "" {
		username = utils.NoUsernameConst
	}

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | command=/help | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
