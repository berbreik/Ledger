package transactions

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"ledger/internal/accounts"
)

type service struct {
	accRepo    accounts.Repository
	ledgerRepo LedgerRepository
}

func NewService(acc accounts.Repository, ledger LedgerRepository) UseCase {
	return &service{
		accRepo:    acc,
		ledgerRepo: ledger,
	}
}

func (s *service) Process(ctx context.Context, tx *Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("invalid transaction amount")
	}

	tx.ID = uuid.NewString()
	tx.CreatedAt = time.Now()
	tx.Status = "PENDING"

	delta := tx.Amount
	if tx.Type == Withdrawal {
		delta = -tx.Amount
	}

	err := s.accRepo.UpdateBalance(ctx, tx.AccountID, delta)
	if err != nil {
		tx.Status = "FAILED"
	} else {
		tx.Status = "SUCCESS"
	}

	logErr := s.ledgerRepo.Insert(ctx, tx)
	if logErr != nil {
		return logErr
	}

	return err
}
