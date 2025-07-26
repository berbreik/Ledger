package accounts

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Insert(ctx context.Context, acc *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	UpdateBalance(ctx context.Context, id string, delta int64) error
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Insert(ctx context.Context, acc *Account) error {
	acc.ID = uuid.New().String()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (id, name, balance)
		VALUES ($1, $2, $3)
	`, acc.ID, acc.Name, acc.Balance)
	return err
}

func (r *repo) FindByID(ctx context.Context, id string) (*Account, error) {
	var acc Account
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, balance, created_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(&acc.ID, &acc.Name, &acc.Balance, &acc.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &acc, err
}

func (r *repo) UpdateBalance(ctx context.Context, id string, delta int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, delta, id)
	return err
}
