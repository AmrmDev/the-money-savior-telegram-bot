package utils

import (
	"money-telegram-bot/internal/models"
	"strings"
	"testing"
	"time"
)

func TestSuccessExpenseMessage(t *testing.T) {
	tests := []struct {
		name       string
		expenseID  string
		amount     float64
		category   string
		method     string
		wantFields []string
	}{
		{
			name:       "basic expense",
			expenseID:  "AM123456",
			amount:     50.99,
			category:   "Food",
			method:     "Débito",
			wantFields: []string{"✅", "AM123456", "50.99", "Food", "Débito"},
		},
		{
			name:       "high amount",
			expenseID:  "AM999999",
			amount:     1500.00,
			category:   "Electronics",
			method:     "Crédito",
			wantFields: []string{"✅", "AM999999", "1500.00", "Electronics", "Crédito"},
		},
		{
			name:       "zero amount",
			expenseID:  "AM000001",
			amount:     0.00,
			category:   "Test",
			method:     "Pix",
			wantFields: []string{"✅", "AM000001", "0.00", "Test", "Pix"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SuccessExpenseMessage(tt.expenseID, tt.amount, tt.category, tt.method)

			for _, field := range tt.wantFields {
				if !strings.Contains(result, field) {
					t.Errorf("SuccessExpenseMessage() result does not contain %q: %s", field, result)
				}
			}
		})
	}
}

func TestExpenseDetailsMessage(t *testing.T) {
	tests := []struct {
		name     string
		expense  *models.Expense
		wantText []string
	}{
		{
			name: "full expense details",
			expense: &models.Expense{
				ExpenseID: "AM123456",
				Amount:    100.50,
				Category:  "Food",
				Method:    "Débito",
				CreatedAt: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			},
			wantText: []string{"📄", "AM123456", "100.50", "Food", "Débito", "15/01/2024", "14:30"},
		},
		{
			name: "expense with zero amount",
			expense: &models.Expense{
				ExpenseID: "AM000000",
				Amount:    0.00,
				Category:  "Test",
				Method:    "Pix",
				CreatedAt: time.Date(2024, 12, 31, 23, 59, 0, 0, time.UTC),
			},
			wantText: []string{"📄", "AM000000", "0.00", "Test", "Pix", "31/12/2024", "23:59"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpenseDetailsMessage(tt.expense)

			for _, text := range tt.wantText {
				if !strings.Contains(result, text) {
					t.Errorf("ExpenseDetailsMessage() result does not contain %q: %s", text, result)
				}
			}
		})
	}
}

func TestExpenseListMessage(t *testing.T) {
	tests := []struct {
		name     string
		expenses []models.Expense
		wantText []string
		wantRows int
	}{
		{
			name: "single expense",
			expenses: []models.Expense{
				{
					ExpenseID: "AM111111",
					Amount:    50.00,
					Category:  "Food",
					Method:    "Débito",
				},
			},
			wantText: []string{"📋", "1 registros", "AM111111", "50.00", "Food", "Débito"},
			wantRows: 1,
		},
		{
			name: "multiple expenses",
			expenses: []models.Expense{
				{
					ExpenseID: "AM111111",
					Amount:    50.00,
					Category:  "Food",
					Method:    "Débito",
				},
				{
					ExpenseID: "AM222222",
					Amount:    100.00,
					Category:  "Transport",
					Method:    "Pix",
				},
				{
					ExpenseID: "AM333333",
					Amount:    75.50,
					Category:  "Shopping",
					Method:    "Crédito",
				},
			},
			wantText: []string{"📋", "3 registros", "AM111111", "AM222222", "AM333333"},
			wantRows: 3,
		},
		{
			name:     "empty expenses",
			expenses: []models.Expense{},
			wantText: []string{"📋", "0 registros"},
			wantRows: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpenseListMessage(tt.expenses)

			for _, text := range tt.wantText {
				if !strings.Contains(result, text) {
					t.Errorf("ExpenseListMessage() result does not contain %q: %s", text, result)
				}
			}

			// Count the number of expense rows
			expenseRowCount := strings.Count(result, "🆔")
			if expenseRowCount != tt.wantRows {
				t.Errorf("ExpenseListMessage() expected %d expense rows, got %d", tt.wantRows, expenseRowCount)
			}
		})
	}
}

func TestMsgInvalidCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    []string
	}{
		{
			name:    "help command",
			command: "help",
			want:    []string{"❌", "help", "não existe"},
		},
		{
			name:    "unknown command",
			command: "unknown",
			want:    []string{"❌", "unknown", "não existe"},
		},
		{
			name:    "typo command",
			command: "gastei123",
			want:    []string{"❌", "gastei123", "não existe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MsgInvalidCommand(tt.command)

			for _, text := range tt.want {
				if !strings.Contains(result, text) {
					t.Errorf("MsgInvalidCommand(%q) does not contain %q: %s", tt.command, text, result)
				}
			}
		})
	}
}

func TestConstMessages(t *testing.T) {
	// Test that important message constants are not empty
	tests := []struct {
		name    string
		message string
		minLen  int
	}{
		{"ErrInvalidFormat", ErrInvalidFormat, 10},
		{"ErrInvalidAmount", ErrInvalidAmount, 5},
		{"SuccessDeleteAll", SuccessDeleteAll, 5},
		{"SuccessDeleteExpense", SuccessDeleteExpense, 5},
		{"MsgHelp", MsgHelp, 50},
		{"MsgStart", MsgStart, 50},
		{"NoUsernameConst", NoUsernameConst, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.message) < tt.minLen {
				t.Errorf("%s is too short or empty: %q (len=%d, want >= %d)", tt.name, tt.message, len(tt.message), tt.minLen)
			}
		})
	}
}
