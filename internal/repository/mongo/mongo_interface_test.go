package mongo

import (
	"context"
	"ledger/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection interface {
	InsertOne(ctx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

type LedgerRepositoryt struct {
	collection MongoCollection
}

func (r LedgerRepositoryt) SaveEntry(todo context.Context, entry *domain.LedgerEntry) interface{} {
	_, err := r.collection.InsertOne(todo, entry)
	if err != nil {
		return err
	}
	return nil

}

func NewLedgerRepositoryFromInterface(collection MongoCollection) *LedgerRepositoryt {
	return &LedgerRepositoryt{collection: collection}
}
