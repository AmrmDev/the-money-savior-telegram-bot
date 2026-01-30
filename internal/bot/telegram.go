package bot

import (
	"fmt"
	"log"

	"money-telegram-bot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	fmt.Println("")
	log.Printf("authorized bot as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		fmt.Println("")
		log.Printf("update received: %+v\n", update)

		msg := update.Message
		if msg == nil {
			msg = update.EditedMessage
		}

		if msg == nil {
			log.Println("update without message, skipping")
			continue
		}

		log.Printf("message received: %s |", msg.Text)

		if msg.IsCommand() {
			fmt.Println("")
			log.Printf("command received: %s |", msg.Command())

			switch msg.Command() {
			case "start":
				fmt.Println("")
				log.Println("handling /start command")
				handlers.HandleStart(bot, msg)

			case "gastei":
				fmt.Println("")
				log.Println("handling /gastei command")
				handlers.HandleExpense(bot, msg)
			}
		}
	}

	return nil
}