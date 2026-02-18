package main

import (
	"log"
	"os"

	"money-telegram-bot/internal/bot"

)
func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set | please set it to your bot's token.")
	}

	log.Println("Starting Telegram Bot | initializing...")

	if err := bot.Start(token); err != nil {
		log.Fatal(err)
	}
}

