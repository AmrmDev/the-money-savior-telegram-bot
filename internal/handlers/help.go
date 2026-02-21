package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("[INFO] Processing /help command")

	helpText := `🆘 *Ajuda — Comandos Disponíveis*

🚀 *Comandos principais:*

▶️ */start*  
Inicia o bot e exibe a mensagem de boas-vindas.

💸 */gastei <valor> <categoria> [método]*  
Registra uma nova despesa.  
Exemplo: /gastei 45.50 supermercado débito

📋 */consulta*  
Exibe todos os seus gastos com IDs em ordem (1, 2, 3...).

📌 */consulta <ID>*  
Ver detalhes de um gasto específico com navegação ⬅️ ➡️ entre registros.  
Exemplo: /consulta 3

🗑️ */deletar <ID>*  
Deleta um gasto específico pelo ID (com confirmação).  
Exemplo: /deletar 2

❌ */deletartudo*  
Deleta *todos* os gastos registrados (com confirmação).

💡 *Dica:* Os IDs são sequenciais (1, 2, 3...). Use /consulta para ver os IDs antes de deletar.

🔙 Digite */start* para voltar ao menu inicial.`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"

	user := message.From

	username := user.UserName
	if username == "" {
		username = "sem_username"
	}

	lastName := user.LastName
	if lastName == "" {
		lastName = "-"
	}

	bot.Send(msg)

	log.Printf(
		"[INFO] Response sent | chatID=%d | userID=%d | userName=%s | command=/help | status=success",
		message.Chat.ID,
		user.ID,
		username,
	)
}
