package transactions

import "context"

type UseCase interface {
	Process(ctx context.Context, tx *Transaction) error
}
