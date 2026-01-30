package main

import (
	"log"
	"os"

	"money-telegram-bot/internal/bot"

)
func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	if err := bot.Start(token); err != nil {
		log.Fatal(err)
	}

	log.Println("starting telegram bot...")
}