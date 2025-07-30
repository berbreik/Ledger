package domain

import "context"

// Transaction represents a fund transfer between two accounts
type Transaction struct {
	ID            string  `json:"id"` // UUID or auto-increment
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`     // e.g., "PENDING", "SUCCESS", "FAILED"
	CreatedAt     string  `json:"created_at"` // or time.Time
}

// LedgerEntry represents a transaction stored in MongoDB (audit log)
type LedgerEntry struct {
	ID            string  `json:"id"` // UUID
	TransactionID string  `json:"transaction_id"`
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"` // or time.Time
}

// TransactionService handles business logic for transfers
type TransactionService interface {
	ProcessTransaction(ctx context.Context, tx *Transaction) error
	GetTransactionHistory(ctx context.Context, accountID int64) ([]*Transaction, error)
}
