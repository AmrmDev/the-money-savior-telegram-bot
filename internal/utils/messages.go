package utils

import (
	"fmt"
	"log"
	"money-telegram-bot/internal/models"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ErrInvalidFormat           = "⚠️ Formato incorreto — use:\n/gastei <valor> <categoria> [método]\n\nExemplo:\n/gastei 21.90 uber pix"
	ErrInvalidAmount           = "❌ Valor inválido."
	ErrSaveExpense             = "❌ Erro ao salvar gasto."
	ErrExpenseNotFound         = "❌ Nenhum gasto encontrado com esse ID."
	ErrQueryExpenses           = "❌ Erro ao consultar gastos."
	ErrNullExpensesReturn      = "📝 Você ainda não registrou nenhum gasto."
	ErrShouldUseFormatToDelete = "❌ Use: /deletar AM123456"
	ErrDeletingExpenses        = "❌ Erro ao limpar seus gastos."
	ErrInvalidDeleteIdFormat   = "❌ ID inválido. Exemplo válido: AM123456"
	ErrDeletingExpense         = "❌ Erro ao deletar o gasto."
	ErrNoUserIdentifiedMessage = "❌ Não foi possível identificar o usuário. Tente novamente."

	SuccessDeleteAll     = "🧹 Todos os gastos foram apagados com sucesso."
	SuccessDeleteExpense = "✅ Gasto deletado com sucesso."

	NoUsernameConst = "noUserName"

	MsgHelp = `🆘 *Ajuda — Comandos Disponíveis*

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

	MsgStart = `👋 Bem-vindo ao *Money Savior*!

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
)

func SuccessExpenseMessage(
	expenseID string,
	amount float64,
	category string,
	method string,
) string {
	return fmt.Sprintf(
		"✅ *Gasto registrado com sucesso!*\n\n"+
			"🆔 %s\n"+
			"💰 R$ %.2f\n"+
			"📝 %s\n"+
			"💳 %s",
		expenseID,
		amount,
		category,
		method,
	)
}

func ExpenseDetailsMessage(expense *models.Expense) string {
	return fmt.Sprintf(
		"📄 *Detalhes do gasto*\n\n"+
			"🆔 ID: *%s*\n"+
			"💰 Valor: *R$ %.2f*\n"+
			"📝 Categoria: *%s*\n"+
			"💳 Método: *%s*\n"+
			"🕐 Data: *%s*",
		expense.ExpenseID,
		expense.Amount,
		expense.Category,
		expense.Method,
		expense.CreatedAt.Format("02/01/2006 15:04"),
	)
}

func ReplyMarkdown(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("erro ao enviar mensagem: %v", err)
	}
}

func ExpenseListMessage(expenses []models.Expense) string {
	var response strings.Builder

	response.WriteString(
		fmt.Sprintf("📋 *Seus gastos (%d registros):*\n\n", len(expenses)),
	)

	for _, e := range expenses {
		response.WriteString(fmt.Sprintf(
			"🆔 *%s* | 💰 R$ %.2f | 📝 %s | 💳 %s\n",
			e.ExpenseID,
			e.Amount,
			e.Category,
			e.Method,
		))
	}

	response.WriteString("\n💡 Use `/consulta AMxxxxxx` para ver detalhes.")

	return response.String()
}

func MsgInvalidCommand(command string) string {
	return fmt.Sprintf(
		`❌ *Comando não reconhecido*

O comando *%s* não existe no *Money Savior* 😕  

📋 Para ver a lista de comandos disponíveis, digite */help*.

💡 Dica: confira se o comando foi digitado corretamente.`,
		command,
	)
}
