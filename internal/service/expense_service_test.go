package service

import (
	"context"
	"errors"
	"money-telegram-bot/internal/models"
	"money-telegram-bot/internal/repository"
	"testing"
)

// MockExpenseRepository is a mock implementation of ExpenseRepository for testing
type MockExpenseRepository struct {
	SaveFunc              func(ctx context.Context, expense *models.Expense) error
	FindByIDFunc          func(ctx context.Context, userID, expenseID string) (*models.Expense, error)
	FindByUserFunc        func(ctx context.Context, userID string) ([]models.Expense, error)
	DeleteAllExpensesFunc func(ctx context.Context, userID int64) error
	DeleteByIDFunc        func(ctx context.Context, userID, expenseID string) error
}

func (m *MockExpenseRepository) Save(ctx context.Context, expense *models.Expense) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, expense)
	}
	return nil
}

func (m *MockExpenseRepository) FindByID(ctx context.Context, userID, expenseID string) (*models.Expense, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, userID, expenseID)
	}
	return nil, nil
}

func (m *MockExpenseRepository) FindByUser(ctx context.Context, userID string) ([]models.Expense, error) {
	if m.FindByUserFunc != nil {
		return m.FindByUserFunc(ctx, userID)
	}
	return []models.Expense{}, nil
}

func (m *MockExpenseRepository) DeleteAllExpenses(ctx context.Context, userID int64) error {
	if m.DeleteAllExpensesFunc != nil {
		return m.DeleteAllExpensesFunc(ctx, userID)
	}
	return nil
}

func (m *MockExpenseRepository) DeleteByID(ctx context.Context, userID, expenseID string) error {
	if m.DeleteByIDFunc != nil {
		return m.DeleteByIDFunc(ctx, userID, expenseID)
	}
	return nil
}

var _ repository.ExpenseRepository = (*MockExpenseRepository)(nil)

func TestNewExpenseService(t *testing.T) {
	mockRepo := &MockExpenseRepository{}
	service := NewExpenseService(mockRepo)

	if service == nil {
		t.Error("NewExpenseService returned nil")
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		expenseID string
		mockFunc  func(ctx context.Context, userID, expenseID string) (*models.Expense, error)
		wantErr   bool
		wantData  *models.Expense
	}{
		{
			name:      "successful get",
			userID:    "123",
			expenseID: "AM123456",
			mockFunc: func(ctx context.Context, userID, expenseID string) (*models.Expense, error) {
				return &models.Expense{
					UserID:    123,
					ExpenseID: "AM123456",
					Amount:    50.0,
					Category:  "Food",
				}, nil
			},
			wantErr: false,
			wantData: &models.Expense{
				UserID:    123,
				ExpenseID: "AM123456",
				Amount:    50.0,
				Category:  "Food",
			},
		},
		{
			name:      "expense not found",
			userID:    "123",
			expenseID: "AM000000",
			mockFunc: func(ctx context.Context, userID, expenseID string) (*models.Expense, error) {
				return nil, errors.New("not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockExpenseRepository{
				FindByIDFunc: tt.mockFunc,
			}
			service := NewExpenseService(mockRepo)

			result, err := service.GetByID(ctx, tt.userID, tt.expenseID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && result != nil {
				if result.ExpenseID != tt.wantData.ExpenseID {
					t.Errorf("GetByID() returned wrong ExpenseID: got %v, want %v", result.ExpenseID, tt.wantData.ExpenseID)
				}
				if result.Amount != tt.wantData.Amount {
					t.Errorf("GetByID() returned wrong Amount: got %v, want %v", result.Amount, tt.wantData.Amount)
				}
			}
		})
	}
}

func TestListByUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		mockFunc  func(ctx context.Context, userID string) ([]models.Expense, error)
		wantErr   bool
		wantCount int
	}{
		{
			name:   "successful list",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) ([]models.Expense, error) {
				return []models.Expense{
					{ExpenseID: "AM111111", Amount: 50.0},
					{ExpenseID: "AM222222", Amount: 100.0},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:   "empty list",
			userID: "456",
			mockFunc: func(ctx context.Context, userID string) ([]models.Expense, error) {
				return []models.Expense{}, nil
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:   "query error",
			userID: "789",
			mockFunc: func(ctx context.Context, userID string) ([]models.Expense, error) {
				return nil, errors.New("database error")
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockExpenseRepository{
				FindByUserFunc: tt.mockFunc,
			}
			service := NewExpenseService(mockRepo)

			result, err := service.ListByUser(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListByUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(result) != tt.wantCount {
				t.Errorf("ListByUser() returned %d expenses, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestDeleteAllExpenses(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		userID   int64
		mockFunc func(ctx context.Context, userID int64) error
		wantErr  bool
	}{
		{
			name:   "successful delete all",
			userID: 123,
			mockFunc: func(ctx context.Context, userID int64) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "delete error",
			userID: 456,
			mockFunc: func(ctx context.Context, userID int64) error {
				return errors.New("delete failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockExpenseRepository{
				DeleteAllExpensesFunc: tt.mockFunc,
			}
			service := NewExpenseService(mockRepo)

			err := service.DeleteAllExpenses(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAllExpenses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteExpense(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    int64
		expenseID string
		mockFunc  func(ctx context.Context, userID, expenseID string) error
		wantErr   bool
	}{
		{
			name:      "successful delete",
			userID:    123,
			expenseID: "AM123456",
			mockFunc: func(ctx context.Context, userID, expenseID string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:      "delete error",
			userID:    456,
			expenseID: "AM000000",
			mockFunc: func(ctx context.Context, userID, expenseID string) error {
				return errors.New("delete failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockExpenseRepository{
				DeleteByIDFunc: tt.mockFunc,
			}
			service := NewExpenseService(mockRepo)

			err := service.DeleteExpense(ctx, tt.userID, tt.expenseID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteExpense() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateExpense(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		userID       int64
		chatID       int64
		username     string
		amount       float64
		category     string
		method       string
		mockFunc     func(ctx context.Context, expense *models.Expense) error
		wantErr      bool
		wantCategory string
		wantMethod   string
	}{
		{
			name:     "successful create with formatted inputs",
			userID:   123,
			chatID:   456,
			username: "testuser",
			amount:   50.99,
			category: "food",
			method:   "debito",
			mockFunc: func(ctx context.Context, expense *models.Expense) error {
				// Assign an ID after mock save
				expense.ExpenseID = "AM123456"
				return nil
			},
			wantErr:      false,
			wantCategory: "Food",
			wantMethod:   "Débito",
		},
		{
			name:     "create with save error",
			userID:   789,
			chatID:   012,
			username: "another",
			amount:   100.0,
			category: "transport",
			method:   "pix",
			mockFunc: func(ctx context.Context, expense *models.Expense) error {
				return errors.New("save failed")
			},
			wantErr: true,
		},
		{
			name:     "create with pix method",
			userID:   111,
			chatID:   222,
			username: "user3",
			amount:   75.50,
			category: "shopping",
			method:   "PIX",
			mockFunc: func(ctx context.Context, expense *models.Expense) error {
				expense.ExpenseID = "AM999999"
				return nil
			},
			wantErr:      false,
			wantCategory: "Shopping",
			wantMethod:   "Pix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockExpenseRepository{
				SaveFunc: tt.mockFunc,
			}
			service := NewExpenseService(mockRepo)

			result, err := service.CreateExpense(ctx, tt.userID, tt.chatID, tt.username, tt.amount, tt.category, tt.method)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateExpense() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && result != nil {
				if result.UserID != tt.userID {
					t.Errorf("CreateExpense() UserID = %d, want %d", result.UserID, tt.userID)
				}
				if result.Amount != tt.amount {
					t.Errorf("CreateExpense() Amount = %f, want %f", result.Amount, tt.amount)
				}
				if result.Category != tt.wantCategory {
					t.Errorf("CreateExpense() Category = %q, want %q", result.Category, tt.wantCategory)
				}
				if result.Method != tt.wantMethod {
					t.Errorf("CreateExpense() Method = %q, want %q", result.Method, tt.wantMethod)
				}
			}
		})
	}
}
