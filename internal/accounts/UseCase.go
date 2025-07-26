package accounts

import "context"

type UseCase interface {
	CreateAccount(ctx context.Context, name string, initialBalance int64) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	AdjustBalance(ctx context.Context, id string, delta int64) error
}
