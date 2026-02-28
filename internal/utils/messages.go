package utils

import "fmt"

const (
	ErrInvalidFormat      = "⚠️ Formato incorreto — use:\n/gastei <valor> <categoria> [método]\n\nExemplo:\n/gastei 21.90 uber pix"
	ErrInvalidAmount      = "❌ Valor inválido."
	ErrSaveExpense        = "❌ Erro ao salvar gasto."
	ErrExpenseNotFound    = "❌ Nenhum gasto encontrado com esse ID."
	ErrQueryExpenses      = "❌ Erro ao consultar gastos."
	ErrNullExpensesReturn = "📝 Você ainda não registrou nenhum gasto."
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
