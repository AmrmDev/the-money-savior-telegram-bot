package models

import (
	"testing"
	"time"
)

func TestExpenseStructure(t *testing.T) {
	// Test that Expense struct can be created and fields set correctly
	now := time.Now().UTC()
	expense := &Expense{
		UserID:    12345,
		ChatID:    67890,
		Username:  "testuser",
		Amount:    100.50,
		Category:  "Food",
		Method:    "Débito",
		CreatedAt: now,
		ExpenseID: "AM123456",
		DisplayID: "1",
	}

	// Verify all fields are set correctly
	if expense.UserID != 12345 {
		t.Errorf("UserID = %d, want 12345", expense.UserID)
	}
	if expense.ChatID != 67890 {
		t.Errorf("ChatID = %d, want 67890", expense.ChatID)
	}
	if expense.Username != "testuser" {
		t.Errorf("Username = %q, want %q", expense.Username, "testuser")
	}
	if expense.Amount != 100.50 {
		t.Errorf("Amount = %f, want 100.50", expense.Amount)
	}
	if expense.Category != "Food" {
		t.Errorf("Category = %q, want %q", expense.Category, "Food")
	}
	if expense.Method != "Débito" {
		t.Errorf("Method = %q, want %q", expense.Method, "Débito")
	}
	if expense.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", expense.CreatedAt, now)
	}
	if expense.ExpenseID != "AM123456" {
		t.Errorf("ExpenseID = %q, want %q", expense.ExpenseID, "AM123456")
	}
	if expense.DisplayID != "1" {
		t.Errorf("DisplayID = %q, want %q", expense.DisplayID, "1")
	}
}

func TestExpenseZeroValues(t *testing.T) {
	// Test Expense with zero/empty values
	expense := &Expense{}

	if expense.UserID != 0 {
		t.Errorf("UserID default = %d, want 0", expense.UserID)
	}
	if expense.ChatID != 0 {
		t.Errorf("ChatID default = %d, want 0", expense.ChatID)
	}
	if expense.Amount != 0.0 {
		t.Errorf("Amount default = %f, want 0.0", expense.Amount)
	}
	if expense.Username != "" {
		t.Errorf("Username default = %q, want empty string", expense.Username)
	}
	if expense.Category != "" {
		t.Errorf("Category default = %q, want empty string", expense.Category)
	}
	if expense.Method != "" {
		t.Errorf("Method default = %q, want empty string", expense.Method)
	}
}

func TestExpenseWithLargeValues(t *testing.T) {
	// Test with large values
	expense := &Expense{
		UserID:   9223372036854775807, // max int64
		Amount:   999999999.99,
		Username: "user_with_very_long_name_that_could_be_problematic",
		Category: "A category with many words separated by spaces and special chars!@#$%",
	}

	if expense.UserID != 9223372036854775807 {
		t.Errorf("UserID with large value failed")
	}
	if expense.Amount != 999999999.99 {
		t.Errorf("Amount with large value failed")
	}
	if expense.Username == "" {
		t.Errorf("Username should not be empty")
	}
	if expense.Category == "" {
		t.Errorf("Category should not be empty")
	}
}
