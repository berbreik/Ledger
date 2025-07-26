package accounts

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, acc *Account) error
	GetById(ctx context.Context, id string) (*Account, error)
	Update(ctx context.Context, id string, amount int64) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, acc *Account) error {
	acc.ID = uuid.New().String()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (id, name, balance)
		VALUES ($1, $2, $3)
	`, acc.ID, acc.Name, acc.Balance)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetById(ctx context.Context, id string) (*Account, error) {
	var acc Account
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, balance, created_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(&acc.ID, &acc.Name, &acc.Balance, &acc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) Update(ctx context.Context, id string, amount int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, amount, id)
	return err
}
