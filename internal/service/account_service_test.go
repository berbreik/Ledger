package service_test

import (
	"context"
	"errors"
	"ledger/internal/domain"
	"ledger/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) Create(ctx context.Context, acc *domain.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

func (m *MockAccountRepo) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) GetAll(ctx context.Context) ([]*domain.Account, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) UpdateBalance(ctx context.Context, id string, newBalance float64) error {
	args := m.Called(ctx, id, newBalance)
	return args.Error(0)
}

func (m *MockAccountRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	err := svc.CreateAccount(context.Background(), "John", 100.0)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAccount(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	expected := &domain.Account{ID: "acc1", OwnerName: "John", Balance: 100.0}
	mockRepo.On("GetByID", mock.Anything, "acc1").Return(expected, nil)

	account, err := svc.GetAccount(context.Background(), "acc1")

	assert.NoError(t, err)
	assert.Equal(t, expected, account)
	mockRepo.AssertExpectations(t)
}

func TestGetAllAccounts(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	expected := []*domain.Account{
		{ID: "acc1", OwnerName: "John", Balance: 100},
		{ID: "acc2", OwnerName: "Jane", Balance: 200},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

	accounts, err := svc.GetAllAccounts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, accounts)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccountBalance(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("UpdateBalance", mock.Anything, "acc1", 150.0).Return(nil)

	err := svc.UpdateAccountBalance(context.Background(), "acc1", 150.0)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "acc1").Return(nil)

	err := svc.DeleteAccount(context.Background(), "acc1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_Failure(t *testing.T) {
	mockRepo := new(MockAccountRepo)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("insert failed"))

	err := svc.CreateAccount(context.Background(), "John", 100.0)

	assert.Error(t, err)
	assert.EqualError(t, err, "insert failed")
	mockRepo.AssertExpectations(t)
}
