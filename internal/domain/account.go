package domain

import "context"

// Account represents a bank account entity
type Account struct {
	ID        string  `json:"id"`
	OwnerName string  `json:"owner_name"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"` // or time.Time if you prefer
}

// AccountRepository defines DB operations related to accounts
type AccountRepository interface {
	Create(ctx context.Context, acc *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetAll(ctx context.Context) ([]*Account, error)
	UpdateBalance(ctx context.Context, id string, newBalance float64) error
	Delete(ctx context.Context, id string) error
}

// AccountService defines business logic operations
type AccountService interface {
	CreateAccount(ctx context.Context, ownerName string, initialBalance float64) error
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAllAccounts(ctx context.Context) ([]*Account, error)
	UpdateAccountBalance(ctx context.Context, id string, newBalance float64) error
	DeleteAccount(ctx context.Context, id string) error
}

var ErrAccountNotFound = "account not found"
