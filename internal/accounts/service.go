package accounts

import "context"

type service struct {
	repo Repository
}

func NewService(r Repository) Repository {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, acc *Account) error {
	return s.repo.Create(ctx, acc)
}

func (s *service) GetById(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetById(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, amount int64) error {
	return s.repo.Update(ctx, id, amount)
}
