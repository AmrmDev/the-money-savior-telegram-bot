package repository

import (
	"context"
	"money-telegram-bot/internal/models"
)

type (
	ExpenseRepository interface {
		Save(ctx context.Context, expense *models.Expense) error
		FindByID(ctx context.Context, userID, expenseID string) (*models.Expense, error)
		FindByUser(ctx context.Context, userID string) ([]models.Expense, error)
		DeleteAllExpenses(ctx context.Context, userID int64) error
		DeleteByID(ctx context.Context, userID, expenseID string) error
	}
)
