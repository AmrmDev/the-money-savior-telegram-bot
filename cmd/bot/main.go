package main

import (
	"log"
	"os"

	"money-telegram-bot/internal/bot"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("[FATAL] TELEGRAM_BOT_TOKEN environment variable not configured. Please set this environment variable with your bot's token.")
	}

	log.Println("[INFO] Starting Money Savior Telegram Bot...")

	if err := bot.Start(token); err != nil {
		log.Fatal("[FATAL] Bot initialization failed:", err)
	}
}
