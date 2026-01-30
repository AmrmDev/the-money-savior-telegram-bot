package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"money-telegram-bot/internal/handlers"
)

func Start(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf("authorized bot as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("update received: %+v\n", update)

		msg := update.Message
		if msg == nil {
			msg = update.EditedMessage
		}

		if msg == nil {
			log.Println("update without message, skipping")
			continue
		}

		log.Printf("message received: %s", msg.Text)

		if msg.IsCommand() {
			log.Printf("command received: %s", msg.Command())

			switch msg.Command() {
			case "start":
				log.Println("handling /start command")
				handlers.HandleStart(bot, msg)

			case "help":
				log.Println("handling /help command")
				handlers.HandleHelp(bot, msg)

			case "gastei":
				log.Println("handling /gastei command")
				handlers.HandleExpense(bot, msg)

			default:
				log.Printf("unknown command received: %s", msg.Command())
				handlers.HandleInvalidCommand(bot, msg)
			}
		}
	}

	return nil
}