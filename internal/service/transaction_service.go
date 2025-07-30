package service

import (
	"context"
	"errors"
	"ledger/internal/domain"
	"ledger/internal/queue"
	"strconv"
	"time"
)

type TransactionService struct {
	accountRepo  domain.AccountRepository
	ledgerRepo   domain.LedgerRepository
	transactionQ queue.TransactionPublisher
}

func NewTransactionService(accountRepo domain.AccountRepository, ledgerRepo domain.LedgerRepository, transactionQ queue.TransactionPublisher) *TransactionService {
	return &TransactionService{
		accountRepo:  accountRepo,
		ledgerRepo:   ledgerRepo,
		transactionQ: transactionQ,
	}
}

func (s *TransactionService) ProcessTransaction(ctx context.Context, tx *domain.Transaction) error {

	fromAccount, err := s.accountRepo.GetByID(ctx, strconv.FormatInt(tx.FromAccountID, 10))
	if err != nil {
		return errors.New("failed to fetch source account: " + err.Error())
	}

	// Check for sufficient balance
	if fromAccount.Balance < tx.Amount {
		return errors.New("insufficient funds")
	}

	// Deduct from source and credit to destination
	err = s.accountRepo.UpdateBalance(ctx, strconv.FormatInt(tx.FromAccountID, 10), -tx.Amount)
	if err != nil {
		return errors.New("failed to debit source account: " + err.Error())
	}

	err = s.accountRepo.UpdateBalance(ctx, strconv.FormatInt(tx.ToAccountID, 10), tx.Amount)
	if err != nil {
		return errors.New("failed to credit destination account: " + err.Error())
	}

	// Prepare ledger entry
	ledger := &domain.LedgerEntry{
		ID:            tx.ID,
		TransactionID: tx.ID,
		FromAccountID: tx.FromAccountID,
		ToAccountID:   tx.ToAccountID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Status:        "SUCCESS",
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	// Store in MongoDB
	err = s.ledgerRepo.SaveEntry(ctx, ledger)
	if err != nil {
		return errors.New("failed to log transaction: " + err.Error())
	}

	return nil
}

func (s *TransactionService) GetTransactionHistory(ctx context.Context, accountID int64) ([]*domain.Transaction, error) {
	// Fetch transaction history from the ledger repository
	transactions, err := s.ledgerRepo.GetEntriesByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	// Convert ledger entries to transactions
	var txs []*domain.Transaction
	for _, entry := range transactions {
		txs = append(txs, &domain.Transaction{
			ID:            entry.TransactionID,
			FromAccountID: entry.FromAccountID,
			ToAccountID:   entry.ToAccountID,
			Amount:        entry.Amount,
			Currency:      entry.Currency,
			Status:        entry.Status,
			CreatedAt:     entry.Timestamp, // Assuming CreatedAt is the same as Timestamp
		})
	}
	if len(txs) == 0 {

		return nil, errors.New("no transactions found for account")
	}
	return txs, nil
}

func (s *TransactionService) QueueTransaction(ctx context.Context, tx domain.Transaction) error {
	// Publish to RabbitMQ (or Kafka)
	return s.transactionQ.Publish(ctx, tx)
}
