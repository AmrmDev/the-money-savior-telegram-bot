package handlers

import (
	"log"
	"money-telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println(utils.InfoProcessingStart)

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

	log.Printf(utils.InfoCommandResponseSent, message.Chat.ID, user.ID, username, "start")
}
