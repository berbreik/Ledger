// file: mongo/ledger_repository_test.go

package mongo_test

import (
	"context"
	"errors"
	"ledger/internal/domain"

	mongo1 "ledger/internal/repository/mongo"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCursor simulates a MongoDB cursor
type MockCursor struct {
	mock.Mock
	entries []domain.LedgerEntry
	index   int
}

func (m *MockCursor) Next(ctx context.Context) bool {
	if m.index >= len(m.entries) {
		return false
	}
	m.index++
	return true
}

func (m *MockCursor) Decode(val interface{}) error {
	if m.index == 0 || m.index > len(m.entries) {
		return errors.New("invalid index")
	}
	ptr := val.(*domain.LedgerEntry)
	*ptr = m.entries[m.index-1]
	return nil
}

func (m *MockCursor) Close(ctx context.Context) error {
	return nil
}

func (m *MockCursor) Err() error {
	return nil
}

// MockCollection mocks MongoCollection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, doc)
	return nil, args.Error(1)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

// SimulatedCursor for unit test (wrapped)
type SimulatedCursor struct {
	*MockCursor
}

func TestSaveEntry_Success(t *testing.T) {
	mockCol := new(MockCollection)
	repo := mongo1.NewLedgerRepositoryFromInterface(mockCol)

	entry := &domain.LedgerEntry{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        100,
		Timestamp:     time.Now().String(),
	}

	mockCol.On("InsertOne", mock.Anything, entry).Return(nil, nil)

	err := repo.SaveEntry(context.TODO(), entry)
	assert.NoError(t, err.(error))
	mockCol.AssertExpectations(t)
}

func TestGetEntriesByAccountID_Success(t *testing.T) {
	mockCol := new(MockCollection)
	_ = mongo1.NewLedgerRepositoryFromInterface(mockCol)

	// Simulate expected entries
	expected := []domain.LedgerEntry{
		{
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        100,
			Timestamp:     time.Now().String(),
		},
		{
			FromAccountID: 3,
			ToAccountID:   1,
			Amount:        50,
			Timestamp:     time.Now().String(),
		},
	}

	mockCursor := &MockCursor{entries: expected}
	cursor := &mongo.Cursor{} // dummy placeholder

	mockCol.On("Find", mock.Anything, mock.Anything).Return(cursor, nil)
	mockCursor.On("Next", mock.Anything).Return(true).Once()
	mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		val := args.Get(0).(*domain.LedgerEntry)
		*val = expected[0]

	}).Return(nil).Once()
	mockCursor.On("Next", mock.Anything).Return(true).Once()
	mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		val := args.Get(0).(*domain.LedgerEntry)
		*val = expected[1]
	}).Return(nil).Once()

}
