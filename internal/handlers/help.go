package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Println("────────── [HANDLER] /help ──────────")

	helpText := `**Ajuda - Comandos Disponíveis**

**/start**
Inicia o bot e exibe uma mensagem de boas-vindas com um resumo dos comandos principais.

**/gastei**
Registra uma nova despesa em seu histórico.

**Como usar:**
/gastei <valor> <categoria> [método]

**Exemplos:**
• /gastei 45.50 supermercado débito
• /gastei 12.99 uber pix
• /gastei 150.00 academia

**Parâmetros:**
• <valor>: Valor da despesa em reais (ex: 45.50)
• <categoria>: Tipo/categoria da despesa (ex: supermercado, uber, etc)
• [método]: Opcional - Método de pagamento (ex: débito, crédito, pix)

Se tiver dúvidas, digite /start para voltar ao menu inicial.`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"

	log.Printf("[BOT] sending help message | chatID=%d", message.Chat.ID)
	bot.Send(msg)

	log.Println("────────── [END] /help ──────────\n")
}

