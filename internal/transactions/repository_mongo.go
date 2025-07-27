package transactions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LedgerRepository interface {
	Insert(ctx context.Context, tx *Transaction) error
}

type ledgerRepo struct {
	coll *mongo.Collection
}

func NewLedgerRepo(db *mongo.Database) LedgerRepository {
	return &ledgerRepo{
		coll: db.Collection("transactions"),
	}
}

func (r *ledgerRepo) Insert(ctx context.Context, tx *Transaction) error {
	tx.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, tx)
	return err
}
