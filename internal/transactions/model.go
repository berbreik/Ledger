package transactions

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
)

type Transaction struct {
	ID        string          `json:"id" bson:"_id,omitempty"`
	AccountID string          `json:"account_id" bson:"account_id"`
	Type      TransactionType `json:"type" bson:"type"`
	Amount    int64           `json:"amount" bson:"amount"`
	Status    string          `json:"status" bson:"status"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
}
