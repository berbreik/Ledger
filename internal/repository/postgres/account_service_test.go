package postgres_test

import (
	"context"
	"database/sql"
	"ledger/internal/domain"
	"ledger/internal/repository/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return db, mock, cleanup
}

func TestGetAll(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "owner_name", "balance"}).
		AddRow("acc1", "Alice", 100.0).
		AddRow("acc2", "Bob", 200.0)

	mock.ExpectQuery(`SELECT id, owner_name, balance FROM accounts`).
		WillReturnRows(rows)

	repo := postgres.NewAccountRepository(db)
	accounts, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, "Alice", accounts[0].OwnerName)
}

func TestDelete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM accounts WHERE id = \$1`).
		WithArgs("acc1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := postgres.NewAccountRepository(db)
	err := repo.Delete(context.Background(), "acc1")

	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	account := &domain.Account{ID: "acc1", OwnerName: "Alice", Balance: 100}

	mock.ExpectExec(`INSERT INTO accounts \(id, owner_name, balance\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(account.ID, account.OwnerName, account.Balance).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := postgres.NewAccountRepository(db)
	err := repo.Create(context.Background(), account)

	assert.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	row := sqlmock.NewRows([]string{"id", "owner_name", "balance"}).
		AddRow("acc1", "Alice", 100.0)

	mock.ExpectQuery(`SELECT id, owner_name, balance FROM accounts WHERE id = \$1`).
		WithArgs("acc1").
		WillReturnRows(row)

	repo := postgres.NewAccountRepository(db)
	acc, err := repo.GetByID(context.Background(), "acc1")

	assert.NoError(t, err)
	assert.Equal(t, "Alice", acc.OwnerName)
	assert.Equal(t, 100.0, acc.Balance)
}

func TestUpdateBalance(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE accounts SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(50.0, "acc1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := postgres.NewAccountRepository(db)
	err := repo.UpdateBalance(context.Background(), "acc1", 50.0)

	assert.NoError(t, err)
}
