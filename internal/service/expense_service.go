package services

import (
    "context"
    "time"

    "money-telegram-bot/internal/models"
    "money-telegram-bot/internal/repository"
    "money-telegram-bot/internal/utils"
)

type ExpenseService struct {
    repo repository.ExpenseRepository
}

func NewExpenseService(repo repository.ExpenseRepository) *ExpenseService {
    return &ExpenseService{repo: repo}
}

func (s *ExpenseService) GetByID(ctx context.Context, userID, expenseID string) (*models.Expense, error)
func (s *ExpenseService) ListByUser(ctx context.Context, userID string) ([]models.Expense, error)

func (s *ExpenseService) CreateExpense(
    ctx context.Context,
    userID int64,
    chatID int64,
    username string,
    amount float64,
    categoryInput string,
    methodInput string,
) (*models.Expense, error) {

    expense := &models.Expense{
        UserID:    userID,
        ChatID:   chatID,
        Username: username,
        Amount:   amount,
        Category: utils.FormatTitle(categoryInput),
        Method:   utils.NormalizeMethod(methodInput),
        CreatedAt: time.Now().UTC(),
    }

    if err := s.repo.Save(ctx, expense); err != nil {
        return nil, err
    }

    return expense, nil
}