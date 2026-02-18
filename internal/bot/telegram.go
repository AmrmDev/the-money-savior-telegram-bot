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
	

	log.Printf("authorized bot as @%s | ready for use.", bot.Self.UserName)
	fmt.Println("")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		logSection("NEW UPDATE")

		msg := update.Message
		if msg == nil {
			msg = update.EditedMessage
		}

		if msg == nil {
			log.Println("no message found | skipping update")
			logSeparator()
			continue
		}

		log.Printf("message received: %q", msg.Text)

		if msg.IsCommand() {
			command := msg.Command()
			log.Printf("command received: %s", command)

			switch command {
			case "start":
				log.Println("routing to handler: start")
				fmt.Println("")
				handlers.HandleStart(bot, msg)
				
			case "help":
				log.Println("routing to handler: help")
				fmt.Println("")
				handlers.HandleHelp(bot, msg)
				
			case "gastei":
				log.Println("routing to handler: gastei")
				fmt.Println("")
				handlers.HandleExpense(bot, msg)
				
			default:
				log.Printf("unknown command: %s | routing to invalid command handler", command)
				handlers.HandleInvalidCommand(bot, msg)
			}
		} else {
			log.Println("message is not a command | ignoring")
		}
	}

	return nil
}

func logSeparator() {
	log.Println("────────────────────────────────────\n")
}

func logSection(title string) {
	log.Println("")
	log.Println("========== " + title + " ==========")
}