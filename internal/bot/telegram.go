package bot

import (
	"log"
	"strings"

	"money-telegram-bot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RouteUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Handle inline keyboard callbacks
	if update.CallbackQuery != nil {
		cb := update.CallbackQuery
		data := cb.Data
		log.Printf("[INFO] Callback received: %q | userID=%d", data, cb.From.ID)

		switch {
		case data == "confirm_delete_all" || data == "cancel_delete_all":
			handlers.HandleDeleteAllCallback(bot, cb)
		case data == "confirm_delete" || data == "cancel_delete" || strings.HasPrefix(data, "confirm_delete:"):
			handlers.HandleConfirmDeleteCallback(bot, cb)
		default:
			// navigation callbacks (qnav:...) and list
			handlers.HandleQueryCallback(bot, cb)
		}
		return
	}

	msg := update.Message
	if msg == nil {
		msg = update.EditedMessage
	}

	if msg == nil {
		log.Println("[DEBUG] Update received with no message. Skipping...")
		return
	}

	log.Printf("[INFO] Message received: %q", msg.Text)

	if msg.IsCommand() {
		command := msg.Command()
		log.Printf("[INFO] Command received: /%s", command)

		switch command {
		case "start":
			handlers.HandleStart(bot, msg)
		case "help":
			handlers.HandleHelp(bot, msg)
		case "gastei":
			handlers.HandleExpense(bot, msg)
		case "consulta":
			handlers.HandleQuery(bot, msg)
		case "deletar":
			handlers.HandleDelete(bot, msg)
		case "deletartudo":
			handlers.HandleDeleteAll(bot, msg)
		default:
			log.Printf("[WARN] Unknown command received: /%s", command)
			handlers.HandleInvalidCommand(bot, msg)
		}
	}
}


func Start(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Bot authenticated successfully as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		RouteUpdate(bot, update)
		}
	return nil
}

