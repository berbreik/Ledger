package domain

import "context"

// LedgerRepository defines how ledger entries are persisted and queried from MongoDB
type LedgerRepository interface {
	SaveEntry(ctx context.Context, entry *LedgerEntry) error
	GetEntriesByAccountID(ctx context.Context, accountID int64) ([]*LedgerEntry, error)
}
