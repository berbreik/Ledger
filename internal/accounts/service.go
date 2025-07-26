package accounts

import "context"

type service struct {
	repo Repository
}

func NewService(r Repository) UseCase {
	return &service{repo: r}
}

func (s *service) CreateAccount(ctx context.Context, name string, initialBalance int64) (*Account, error) {
	account := &Account{
		Name:    name,
		Balance: initialBalance,
	}
	err := s.repo.Insert(ctx, account)
	return account, err
}

func (s *service) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) AdjustBalance(ctx context.Context, id string, delta int64) error {
	return s.repo.UpdateBalance(ctx, id, delta)
}
