package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ledger/internal/domain"
)

type LedgerRepository struct {
	collection *mongo.Collection
}

func NewLedgerRepository(client *mongo.Client, dbName, collectionName string) *LedgerRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &LedgerRepository{collection: collection}
}

func (r *LedgerRepository) SaveEntry(ctx context.Context, entry *domain.LedgerEntry) error {
	_, err := r.collection.InsertOne(ctx, entry)
	if err != nil {
		return err
	}
	return nil
}

func (r *LedgerRepository) GetEntriesByAccountID(ctx context.Context, accountID int64) ([]*domain.LedgerEntry, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_account_id": accountID},
			{"to_account_id": accountID},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}(cursor, ctx)

	var entries []*domain.LedgerEntry
	for cursor.Next(ctx) {
		var entry domain.LedgerEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, cursor.Err()
}
