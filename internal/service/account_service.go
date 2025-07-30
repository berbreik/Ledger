package service

import (
	"context"

	"github.com/google/uuid"
	"ledger/internal/domain"
)

type AccountService struct {
	accountRepo domain.AccountRepository
}

func NewAccountService(accountRepo domain.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(ctx context.Context, ownerName string, initialBalance float64) error {
	account := domain.Account{
		ID:        uuid.New().String(),
		OwnerName: ownerName,
		Balance:   initialBalance,
	}

	err := s.accountRepo.Create(ctx, &account)
	if err != nil {
		return err
	}

	return nil
}
func (s *AccountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAllAccounts(ctx context.Context) ([]*domain.Account, error) {
	accounts, err := s.accountRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *AccountService) UpdateAccountBalance(ctx context.Context, id string, newBalance float64) error {
	err := s.accountRepo.UpdateBalance(ctx, id, newBalance)
	if err != nil {
		return err
	}
	return nil
}
func (s *AccountService) DeleteAccount(ctx context.Context, id string) error {
	err := s.accountRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
