package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /start command")

	welcomeText := `👋 Bem-vindo ao *Money Savior*!

💰 Seu assistente pessoal para controle de gastos e organização financeira.

📌 *Comandos disponíveis:*

➕ /gastei — Registre um novo gasto  
Exemplo: /gastei 21.90 uber pix

📋 /consulta — Veja todos os gastos (IDs em ordem)

🔎 /consulta <ID> — Veja um gasto específico com navegação  
Exemplo: /consulta 3

🗑️ /deletar <ID> — Delete um gasto pelo ID

❌ /deletartudo — Delete todos os gastos

ℹ️ /help — Veja todos os comandos e exemplos

✨ Dica: os IDs são sequenciais (1, 2, 3...), facilitando o gerenciamento!

Digite /help para mais detalhes 🚀`

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	bot.Send(msg)
	user := message.From

	username := user.UserName
	if username == "" {
		username = "sem_username"
	}

	lastName := user.LastName
	if lastName == "" {
		lastName = "-"
	}

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
