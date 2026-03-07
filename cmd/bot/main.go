package main

import (
	"log"
	"os"

	"money-telegram-bot/internal/bot"
	"money-telegram-bot/internal/utils"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal(utils.FatalEnvironmentVariableNotSet)
	}

	log.Println(utils.InfoStartingBot)

	if err := bot.Start(token); err != nil {
		log.Fatal(utils.FatalBotInitializationFailed, err)
	}
}
