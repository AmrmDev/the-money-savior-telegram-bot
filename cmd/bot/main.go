package main

import (
	"log"
	"os"

	"money-telegram-bot/internal/controller"
	"money-telegram-bot/internal/utils"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal(utils.FatalEnvironmentVariableNotSet)
	}

	log.Println(utils.InfoStartingBot)

	if err := controller.Start(token); err != nil {
		log.Fatal(utils.FatalBotInitializationFailed, err)
	}
}
